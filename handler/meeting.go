package handler

import (
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
	nowIndex int
}

func CulMeeting(request *constant.CulMeetingRequest) (response *constant.CulMeetingResponse, err error) {
	// 协程池方案
	defer ants.Release()
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
	unitFunc, _ := ants.NewPoolWithFunc(100, func(v interface{}) {
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
		spend := helper.TimeDeviation(nowTime, request.StartTime)
		if spend > now {
			now = spend
			percent := float64(100.0*spend) / float64(total)
			log.Println("Progress:", percent, "%", value.index)
		}
		syncValue.nowIndex += 1
		syncValue.Unlock()
		// 无需同步语句
		response.MinSpacing = 0.0
	})

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
	return &constant.CulMeetingResponse{}, nil
}

/*
	多协程方案
*/
//nowTime := request.StartTime
//var wg sync.WaitGroup
//for helper.DayBigger(request.EndTime, nowTime) {
//	wg.Add(1)
//	go func(time *constant.Data) {
//		defer wg.Done()
//		response, err := CulSpacing(
//			&constant.CulSpacingRequest{
//				Time: time,
//				DeltaT: request.DeltaT,
//			},
//		)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		fmt.Println(time, ":", response.MinSpacing)
//	}(nowTime)
//	nowTime = helper.DayIncrease(nowTime, request.TimeRange)
//}
//wg.Wait()
