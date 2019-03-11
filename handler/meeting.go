package handler

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"github.com/panjf2000/ants"
	"log"
	"sync"
	"time"
)

type unitFuncValue struct {
	index int
	time  *constant.Data
}

type syncSafe struct {
	sync.Mutex
	nowIndex                      int
	shipMeetingList               map[int]map[int]int
	shipMeetingNum                map[int]int
	shipDamageMeetingList         map[int]map[int]int
	shipDamageMeetingNum          map[int]int
	shipForecastDamageMeetingList map[int]map[int]int
}

/*
	计算会遇
	TODO:
		(P0)会遇热力图展示
*/
func CulMeeting(request *constant.CulMeetingRequest) (response *constant.CulMeetingResponse, err error) {
	// 协程池方案
	defer ants.Release()

	// 数据初始化
	resp := &constant.CulMeetingResponse{
		SimpleMeeting:        0,
		ComplexMeeting:       0,
		SimpleDamageMeeting:  0,
		ComplexDamageMeeting: 0,
	}
	resp.AngleForecastDamageMeeting = make([]int, 12)
	resp.AngleDamageMeetingAvoid = make([]int, 12)
	for i := 0; i < 12; i++ {
		resp.AngleForecastDamageMeeting[i] = 0
		resp.AngleDamageMeetingAvoid[i] = 0
	}
	nowTime := request.StartTime
	var wg sync.WaitGroup
	total := helper.TimeDeviation(request.EndTime, nowTime)
	now := helper.TimeDeviation(nowTime, request.StartTime)
	syncValue := syncSafe{
		nowIndex: 0,
	}
	shipList, _ := sql.GetShip()
	syncValue.shipMeetingList = make(map[int]map[int]int)
	syncValue.shipMeetingNum = make(map[int]int)
	syncValue.shipDamageMeetingList = make(map[int]map[int]int)
	syncValue.shipDamageMeetingNum = make(map[int]int)
	syncValue.shipForecastDamageMeetingList = make(map[int]map[int]int)
	for _, v := range shipList {
		syncValue.shipMeetingList[v.MMSI] = make(map[int]int)
		syncValue.shipMeetingNum[v.MMSI] = 0
		syncValue.shipDamageMeetingList[v.MMSI] = make(map[int]int)
		syncValue.shipDamageMeetingNum[v.MMSI] = 0
		syncValue.shipForecastDamageMeetingList[v.MMSI] = make(map[int]int)
	}

	/*
		计算协程
	*/
	unitFunc, _ := ants.NewPoolWithFunc(16, func(v interface{}) {
		defer wg.Done()
		value := v.(*unitFuncValue)
		nowTime := value.time
		response, err := CulSpacing(
			&constant.CulSpacingRequest{
				Time:   nowTime,
				DeltaT: request.DeltaT,
			},
		)
		if err != nil {
			log.Println(err)
			return
		}
		// 协程同步
		for syncValue.nowIndex != value.index {
			time.Sleep(time.Millisecond)
		}
		syncValue.Lock()
		// 遍历判断
		for k1, v1 := range response.ShipSpacing {
			// main: k1 主船: k1
			ship1 := response.TrackMap[k1]
			// 静态船舶剔除
			if ship1.SOG < constant.StaticShip {
				continue
			}
			newMeetingShipNum := 0
			meetingShipNum := syncValue.shipMeetingNum[k1]
			newDamageMeetingShipNum := 0
			damageMeetingShipNum := syncValue.shipDamageMeetingNum[k1]
			shipInfo := cache.GetShipInfo(k1)
			COG := ship1.COG
			for k2, v2 := range v1 {
				if k1 != k2 {
					// 会遇计算
					if v2 < constant.HalfNauticalMile {
						// 静态船舶剔除
						ship2 := response.TrackMap[k2]
						if ship2.SOG < constant.StaticShip {
							continue
						}
						// 如果之前没有会遇
						if syncValue.shipMeetingList[k1][k2] == 0 {
							syncValue.shipMeetingList[k1][k2] = 1
							newMeetingShipNum += 1
							meetingShipNum += 1
						}
						// 计算危险会遇
						if shipInfo.A != 511 && shipInfo.B != 511 &&
							shipInfo.A != 0 && shipInfo.B != 0 {
							L := float64(shipInfo.A + shipInfo.B)
							// 初筛除
							if v2 <= 10*L {
								a := 5 * L
								b := 2.5 * L
								S := 0.75 * L
								T := 1.1 * L
								x := helper.PositionSpacing(&constant.Position{
									Longitude: ship2.PrePosition.Longitude,
									Latitude:  ship1.PrePosition.Latitude,
								}, &constant.Position{
									Longitude: ship1.PrePosition.Longitude,
									Latitude:  ship1.PrePosition.Latitude,
								})
								if ship2.PrePosition.Longitude < ship1.PrePosition.Longitude {
									x *= -1
								}
								y := helper.PositionSpacing(&constant.Position{
									Longitude: ship1.PrePosition.Longitude,
									Latitude:  ship2.PrePosition.Latitude,
								}, &constant.Position{
									Longitude: ship1.PrePosition.Longitude,
									Latitude:  ship1.PrePosition.Latitude,
								})
								if ship2.PrePosition.Latitude < ship1.PrePosition.Latitude {
									y *= -1
								}
								// 危险接触
								if COG <= 360 && helper.InEllipse(a, b, S, T, x, y, COG) {
									// 如果之前没有会遇
									if syncValue.shipDamageMeetingList[k1][k2] == 0 {
										syncValue.shipDamageMeetingList[k1][k2] = 1
										newDamageMeetingShipNum += 1
										damageMeetingShipNum += 1
									}
								}
							}
						}
						// 会遇点计算, 抛出已经进入船舶领域的情况
						if syncValue.shipDamageMeetingList[k1][k2] == 0 {
							if ship1.COG != ship2.COG {
								// 航向不平行, 计算船舶会遇点
								meetingIntersection := helper.CulMeetingIntersection(ship1, ship2)
								if shipInfo.A != 511 && shipInfo.B != 511 &&
									shipInfo.A != 0 && shipInfo.B != 0 {
									L := float64(shipInfo.A + shipInfo.B)
									// 初筛除
									if meetingIntersection.TCPA > 0 {
										intersectionShip1Position := helper.CulSecondPointPosition(ship1.PrePosition, meetingIntersection.TCPA*ship1.SOG, ship1.COG)
										intersectionShip2Position := helper.CulSecondPointPosition(ship2.PrePosition, meetingIntersection.TCPA*ship2.SOG, ship2.COG)
										a := 5 * L
										b := 2.5 * L
										S := 0.75 * L
										T := 1.1 * L
										x := helper.PositionSpacing(&constant.Position{
											Longitude: intersectionShip2Position.Longitude,
											Latitude:  intersectionShip1Position.Latitude,
										}, &constant.Position{
											Longitude: intersectionShip1Position.Longitude,
											Latitude:  intersectionShip1Position.Latitude,
										})
										if intersectionShip2Position.Longitude < intersectionShip1Position.Longitude {
											x *= -1
										}
										y := helper.PositionSpacing(&constant.Position{
											Longitude: intersectionShip1Position.Longitude,
											Latitude:  intersectionShip2Position.Latitude,
										}, &constant.Position{
											Longitude: intersectionShip1Position.Longitude,
											Latitude:  intersectionShip1Position.Latitude,
										})
										if intersectionShip2Position.Latitude < intersectionShip1Position.Latitude {
											y *= -1
										}
										// 预测危险接触
										if COG <= 360 && helper.InEllipse(a, b, S, T, x, y, COG) {
											// 如果之前没有预测接触
											if syncValue.shipForecastDamageMeetingList[k1][k2] == 0 {
												Azimuth := helper.PositionRelativeAzimuth(intersectionShip1Position, ship1.COG, intersectionShip2Position)
												index := int(Azimuth / 30)
												if index == 12 {
													index = 0
												}
												syncValue.shipForecastDamageMeetingList[k1][k2] = 1 + index
												resp.ForecastDamageMeeting += 1
												resp.AngleForecastDamageMeeting[index] += 1
											}
										}
									}
								}
							}
						}
					} else if v2 > constant.NauticalMile {
						// 接触脱离, 如果之前有会遇
						if syncValue.shipMeetingList[k1][k2] == 1 {
							syncValue.shipMeetingList[k1][k2] = 0
							meetingShipNum -= 1
							// 如果之前预测会遇点在船舶领域内部, 规避成功
							if syncValue.shipForecastDamageMeetingList[k1][k2] != 0 && syncValue.shipDamageMeetingList[k1][k2] == 0 {
								resp.DamageMeetingAvoid += 1
								resp.AngleDamageMeetingAvoid[syncValue.shipForecastDamageMeetingList[k1][k2]-1] += 1
							}
							syncValue.shipForecastDamageMeetingList[k1][k2] = 0
							// 接触脱离, 如果之前有危险会遇
							if syncValue.shipDamageMeetingList[k1][k2] == 1 {
								syncValue.shipDamageMeetingList[k1][k2] = 0
								damageMeetingShipNum -= 1
							}
						}
					} else {
						//船舶间距在0.5海里与1海里之前，不做处理
					}
				}
			}
			// 会遇数据汇总
			if meetingShipNum == 1 && newMeetingShipNum == 1 {
				resp.SimpleMeeting += 1
			} else if meetingShipNum > 1 && newMeetingShipNum > 0 {
				resp.ComplexMeeting += 1
			}
			syncValue.shipMeetingNum[k1] = meetingShipNum
			// 危险会遇数据汇总
			if damageMeetingShipNum == 1 && newDamageMeetingShipNum == 1 {
				resp.SimpleDamageMeeting += 1
			} else if damageMeetingShipNum > 1 && newDamageMeetingShipNum > 0 {
				resp.ComplexDamageMeeting += 1
			}
			syncValue.shipDamageMeetingNum[k1] = damageMeetingShipNum
		}
		// 输出当前同步状态
		spend := helper.TimeDeviation(nowTime, request.StartTime)
		if spend > now {
			now = spend
			percent := float64(100.0*spend) / float64(total)
			log.Println("Progress:", percent, "%", value.index)
		}
		syncValue.nowIndex += 1
		syncValue.Unlock()
	})

	// 多协程调用
	index := 0
	for helper.DayBigger(request.EndTime, nowTime) {
		wg.Add(1)
		err := unitFunc.Invoke(&unitFuncValue{
			index: index,
			time:  nowTime,
		})
		if err != nil {
			log.Println(err)
		}
		nowTime = helper.DayIncrease(nowTime, request.TimeRange)
		index += 1
	}
	wg.Wait()

	// 输出结果
	return resp, nil
}
