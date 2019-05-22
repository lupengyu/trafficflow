package handler

import (
	"bufio"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"github.com/panjf2000/ants"
	"log"
	"os"
	"strconv"
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
*/
func CulMeeting(request *constant.CulMeetingRequest) (response *constant.CulMeetingResponse, err error) {
	// 协程池方案
	defer ants.Release()

	// 数据初始化
	danger1, err := os.Create("data/danger1.txt")
	if err != nil {
		log.Println(err)
		return
	}
	danger2, err := os.Create("data/danger2.txt")
	if err != nil {
		log.Println(err)
		return
	}
	danger3, err := os.Create("data/danger3.txt")
	if err != nil {
		log.Println(err)
		return
	}
	danger4, err := os.Create("data/danger4.txt")
	if err != nil {
		log.Println(err)
		return
	}
	danger5, err := os.Create("data/danger5.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		danger1.Close()
		danger2.Close()
		danger3.Close()
		danger4.Close()
		danger5.Close()
	}()
	danger1.Sync()
	danger2.Sync()
	danger3.Sync()
	danger4.Sync()
	danger5.Sync()
	danger1Writer := bufio.NewWriter(danger1)
	danger2Writer := bufio.NewWriter(danger2)
	danger3Writer := bufio.NewWriter(danger3)
	danger4Writer := bufio.NewWriter(danger4)
	danger5Writer := bufio.NewWriter(danger5)
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
			shipInfo := cache.GetShipInfo(k1)
			L := 0.0
			if shipInfo.A != 511 && shipInfo.B != 511 &&
				shipInfo.A != 0 && shipInfo.B != 0 {
				L = float64(shipInfo.A + shipInfo.B)
			} else {
				continue
			}
			if ship1.SOG < constant.StaticShip {
				continue
			}
			// 判断插值后点是不是在合理区间内
			longitudeArea := helper.LongitudeArea(ship1.PrePosition.Longitude, 10)
			latitudeArea := helper.LatitudeArea(ship1.PrePosition.Latitude, 10)
			if longitudeArea == -1 || latitudeArea == -1 {
				continue
			}
			COG := ship1.COG
			if COG > 360 || COG < 0 {
				continue
			}
			I := 0
			II := 0
			III := 0
			IV := 0
			V := 0
			for k2, v2 := range v1 {
				if k1 != k2 {
					ship2 := response.TrackMap[k2]
					if ship2.COG > 360 || ship2.COG < 0 {
						continue
					}
					azimuth := helper.PositionRelativeAzimuth(ship1.PrePosition, ship1.COG, ship2.PrePosition)
					uDCPA := 0.0
					uTCPA := 0.0
					uB := helper.MeetingDangerUB(azimuth)
					uD := 0.0
					uV := helper.MeetingDangerUV(ship1.SOG)
					if ship1.COG != ship2.COG {
						meetingIntersection := helper.CulMeetingIntersection(ship1, ship2)
						if meetingIntersection.TCPA > 0 {
							intersectionShip1Position := helper.CulSecondPointPosition(ship1.PrePosition, meetingIntersection.TCPA*ship1.SOG, ship1.COG)
							intersectionShip2Position := helper.CulSecondPointPosition(ship2.PrePosition, meetingIntersection.TCPA*ship2.SOG, ship2.COG)
							meetingIntersection.Azimuth = helper.PositionRelativeAzimuth(intersectionShip1Position, ship1.COG, intersectionShip2Position)
							meetingIntersection.DCPA = helper.PositionSpacing(intersectionShip1Position, intersectionShip2Position)
							a := 5 * L
							b := 2.5 * L
							S := 0.75 * L
							T := 1.1 * L
							//fmt.Println(meetingIntersection.DCPA, meetingIntersection.TCPA)
							uDCPA = helper.MeetingDangerUDCPA(a, b, S, T, meetingIntersection.Azimuth, meetingIntersection.DCPA)
							uTCPA = helper.MeetingDangerUTCPA(a, b, S, T, meetingIntersection.Azimuth, meetingIntersection.DCPA,
								meetingIntersection.TCPA, meetingIntersection.VR)
							uD = helper.MeetingDangerUD(a, b, S, T, azimuth, v2)
						}
					}
					danger := uD*0.4 + uDCPA*uTCPA*0.4 + uB*0.1 + uV*0.1
					if danger <= 0.2 {
						I += 1
					} else if danger <= 0.4 {
						II += 1
					} else if danger <= 0.6 {
						III += 1
					} else if danger <= 0.8 {
						IV += 1
					} else if danger <= 1.0 {
						V += 1
					}
				}
			}
			// 文件输出
			if I != 0 {
				str := strconv.FormatFloat(ship1.PrePosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(ship1.PrePosition.Latitude, 'f', -1, 64) + "," + strconv.Itoa(I) + "\r\n"
				n, err := danger1Writer.WriteString(str)
				if n != len(str) && err != nil {
					log.Println(err)
				}
				danger1Writer.Flush()
			}
			if II != 0 {
				str := strconv.FormatFloat(ship1.PrePosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(ship1.PrePosition.Latitude, 'f', -1, 64) + "," + strconv.Itoa(II) + "\r\n"
				n, err := danger2Writer.WriteString(str)
				if n != len(str) && err != nil {
					log.Println(err)
				}
				danger2Writer.Flush()
			}
			if III != 0 {
				str := strconv.FormatFloat(ship1.PrePosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(ship1.PrePosition.Latitude, 'f', -1, 64) + "," + strconv.Itoa(III) + "\r\n"
				n, err := danger3Writer.WriteString(str)
				if n != len(str) && err != nil {
					log.Println(err)
				}
				danger3Writer.Flush()
			}
			if IV != 0 {
				str := strconv.FormatFloat(ship1.PrePosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(ship1.PrePosition.Latitude, 'f', -1, 64) + "," + strconv.Itoa(IV) + "\r\n"
				n, err := danger4Writer.WriteString(str)
				if n != len(str) && err != nil {
					log.Println(err)
				}
				danger4Writer.Flush()
			}
			if V != 0 {
				str := strconv.FormatFloat(ship1.PrePosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(ship1.PrePosition.Latitude, 'f', -1, 64) + "," + strconv.Itoa(V) + "\r\n"
				n, err := danger5Writer.WriteString(str)
				if n != len(str) && err != nil {
					log.Println(err)
				}
				danger5Writer.Flush()
			}
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
