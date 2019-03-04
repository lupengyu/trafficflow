package handler

import (
	"fmt"
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
	nowIndex       int
	shipMeetingMap map[int]map[int]string
	meetingIndex   int
}

/*
	计算会遇
	TODO:
		1.加入文件缓存减缓内存占用
*/
func CulMeeting(request *constant.CulMeetingRequest) (response *constant.CulMeetingResponse, err error) {
	// 协程池方案
	defer ants.Release()

	// 数据初始化
	nowTime := request.StartTime
	var wg sync.WaitGroup
	total := helper.TimeDeviation(request.EndTime, nowTime)
	now := helper.TimeDeviation(nowTime, request.StartTime)
	syncValue := syncSafe{
		nowIndex:     0,
		meetingIndex: 0,
	}
	shipList, _ := sql.GetShip()
	syncValue.shipMeetingMap = make(map[int]map[int]string)
	for _, v := range shipList {
		syncValue.shipMeetingMap[v.MMSI] = make(map[int]string)
		for _, v2 := range shipList {
			syncValue.shipMeetingMap[v.MMSI][v2.MMSI] = ""
		}
	}
	// init file cache
	if meetingFileCache {

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
		// 输出当前同步状态
		spend := helper.TimeDeviation(nowTime, request.StartTime)
		if spend > now {
			now = spend
			percent := float64(100.0*spend) / float64(total)
			log.Println("Progress:", percent, "%", value.index)
		}
		syncValue.nowIndex += 1
		syncValue.meetingIndex += 1

		/*
			flag value
				4 船舶距离小于0.5海里
				5 船舶距离大于等于0.5海里小于1海里
				6 船舶距离大于1海里
				9 无两船舶数据
		*/
		for k1, v1 := range response.ShipSpacing {
			for k2, v2 := range v1 {
				flag := ""
				if v2 < constant.HalfNauticalMile {
					// 判断是否小于0.5 nmi
					flag = "4"
				} else if v2 > constant.NauticalMile {
					// 判断是否大于1 nmi
					flag = "6"
				} else {
					// 其他情况
					flag = "5"
				}
				syncValue.shipMeetingMap[k1][k2] += flag
				syncValue.shipMeetingMap[k2][k1] += flag
			}
		}
		meetingIndex := syncValue.meetingIndex
		for k1, v1 := range syncValue.shipMeetingMap {
			for k2, v2 := range v1 {
				if len(v2) != meetingIndex {
					syncValue.shipMeetingMap[k1][k2] += "9"
					syncValue.shipMeetingMap[k2][k1] += "9"
				}
			}
		}
		// flush cache
		if meetingFileCache {

		}
		syncValue.Unlock()
		// 无需同步语句
		response.MinSpacing = 0.0
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

	// Test
	fmt.Println("shipMeetingMap[413788252][413484720]:", syncValue.shipMeetingMap[413788252][413484720])

	// 输出结果
	return &constant.CulMeetingResponse{}, nil
}
