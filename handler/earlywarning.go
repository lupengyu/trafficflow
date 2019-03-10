package handler

import (
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"github.com/panjf2000/ants"
	"log"
	"sync"
	"time"
)

/*
	提前预警
*/
func EarlyWarning(request *constant.EarlyWarningRequest) (response *constant.EarlyWarningResponse, err error) {
	// 协程池方案
	defer ants.Release()
	nowTime := request.StartTime
	var wg sync.WaitGroup
	total := helper.TimeDeviation(request.EndTime, nowTime)
	now := helper.TimeDeviation(nowTime, request.StartTime)
	syncValue := syncSafe{
		nowIndex: 0,
	}
	shipInfo := cache.GetShipInfo(request.MMSI)
	L := 0.0
	if shipInfo.A != 511 && shipInfo.B != 511 &&
		shipInfo.A != 0 && shipInfo.B != 0 {
		L = float64(shipInfo.A + shipInfo.B)
	} else {
		log.Println("ship static info error")
	}
	resp := &constant.EarlyWarningResponse{}
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
		ship1 := response.TrackMap[request.MMSI]
		warning := &constant.Warning{
			MasterShipTrack: ship1,
			Time:            nowTime,
		}
		for k, v := range response.ShipSpacing[request.MMSI] {
			if k != request.MMSI {
				if v < constant.HalfNauticalMile {
					ship2 := response.TrackMap[k]
					if L != 0 {
						// 船舶静态数据有效
						if v <= 2*L {
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
							if ship1.COG <= 360 && helper.InEllipse(a, b, S, T, x, y, ship1.COG) {
								warning.Alerts = append(warning.Alerts, &constant.Alert{
									MMSI:        k,
									IsEmergency: true,
									ShipTrack:   ship2,
									Distance:    v,
									Azimuth:     helper.PositionAzimuth(ship1.PrePosition, ship2.PrePosition),
								})
								continue
							}
						}
					}
					// 会遇点预测
					if ship1.COG != ship2.COG {
						if L != 0 && ship2.SOG >= constant.StaticShip {
							meetingIntersection := helper.CulMeetingIntersection(ship1, ship2)
							if meetingIntersection.TCPA > 0 {
								intersectionShip1Position := helper.CulSecondPointPosition(ship1.PrePosition, meetingIntersection.TCPA*ship1.SOG, ship1.COG)
								intersectionShip2Position := helper.CulSecondPointPosition(ship2.PrePosition, meetingIntersection.TCPA*ship2.SOG, ship2.COG)
								meetingIntersection.Azimuth = helper.PositionAzimuth(intersectionShip1Position, intersectionShip2Position)
								meetingIntersection.DCPA = helper.PositionSpacing(intersectionShip1Position, intersectionShip2Position)
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
								if ship1.COG <= 360 && helper.InEllipse(a, b, S, T, x, y, ship1.COG) {
									// 如果之前没有预测接触
									warning.Alerts = append(warning.Alerts, &constant.Alert{
										MMSI:                k,
										IsEmergency:         false,
										ShipTrack:           ship2,
										Distance:            v,
										Azimuth:             helper.PositionAzimuth(ship1.PrePosition, ship2.PrePosition),
										MeetingIntersection: meetingIntersection,
									})
									continue
								}
							}
						}
					}
					warning.Alerts = append(warning.Alerts, &constant.Alert{
						MMSI:        k,
						IsEmergency: false,
						ShipTrack:   ship2,
						Distance:    v,
						Azimuth:     helper.PositionAzimuth(ship1.PrePosition, ship2.PrePosition),
					})
				}
			}
		}
		if len(warning.Alerts) != 0 {
			resp.Warning = append(resp.Warning, warning)
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

	return resp, nil
}
