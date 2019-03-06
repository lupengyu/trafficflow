package handler

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
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
	nowIndex        int
	shipMeetingList map[int]map[int]int
	shipMeetingNum  map[int]int
}

/*
	计算会遇
	TODO:
		1.加入文件缓存减缓内存占用(×, 代码优化在每轮数据后计算，优化了时间空间复杂度，不需要添加文件缓存)
		2.判断位置是否在船舶领域中
		3.判断会遇中会遇点在船舶领域中的情况
		4.计算会遇中的危险会遇(介入会遇船只的船舶领域)
		5.计算会遇中的规避情况(即原先会出现危险会遇但是经过规避避免了危险会遇)
*/
func CulMeeting(request *constant.CulMeetingRequest) (response *constant.CulMeetingResponse, err error) {
	// 协程池方案
	defer ants.Release()

	// 数据初始化
	resp := &constant.CulMeetingResponse{
		SimpleMeeting:  0,
		ComplexMeeting: 0,
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
	for _, v := range shipList {
		syncValue.shipMeetingList[v.MMSI] = make(map[int]int)
		syncValue.shipMeetingNum[v.MMSI] = 0
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
		/*
			4 船舶距离小于0.5海里
			5 船舶距离大于等于0.5海里小于1海里
			6 船舶距离大于1海里
			9 无两船舶数据
			ship meeting ship value
				"" 	   无会遇
				"1, 3" 与1,3会遇
		*/
		for k1, v1 := range response.ShipSpacing {
			// main: k1
			newMeetingShipNum := 0
			meetingShipNum := syncValue.shipMeetingNum[k1]
			for k2, v2 := range v1 {
				if k1 != k2 {
					if v2 < constant.HalfNauticalMile {
						// 如果之前没有会遇
						if syncValue.shipMeetingList[k1][k2] == 0 {
							syncValue.shipMeetingList[k1][k2] = 1
							newMeetingShipNum += 1
							meetingShipNum += 1
						}
					} else if v2 > constant.NauticalMile {
						// 如果之前有会遇
						if syncValue.shipMeetingList[k1][k2] == 1 {
							syncValue.shipMeetingList[k1][k2] = 0
							meetingShipNum -= 1
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
