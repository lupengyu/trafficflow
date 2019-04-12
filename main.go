package main

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/handler"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"time"
)

/*
	Traffic 船舶交通量统计
*/
func culTraffic() {
	lotDivide := 10
	latDivide := 10
	var day float64 = 1
	response, err := handler.CulTraffic(
		&constant.CulTrafficRequest{
			StartTime: &constant.Data{
				Year:   2018,
				Month:  12,
				Day:    25,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			EndTime: &constant.Data{
				Year:   2018,
				Month:  12,
				Day:    25,
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
			LotDivide: lotDivide,
			LatDivide: latDivide,
		},
	)
	//var day float64 = 12
	//response, err := handler.CulTraffic(
	//	&constant.CulTrafficRequest{
	//		StartTime: &constant.Data{
	//			Year:   2018,
	//			Month:  12,
	//			Day:    22,
	//			Hour:   0,
	//			Minute: 0,
	//			Second: 0,
	//		},
	//		EndTime: &constant.Data{
	//			Year:   2019,
	//			Month:  1,
	//			Day:    2,
	//			Hour:   23,
	//			Minute: 59,
	//			Second: 59,
	//		},
	//		LotDivide: lotDivide,
	//		LatDivide: latDivide,
	//	},
	//)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulTrafficResponsePrint(response, lotDivide, latDivide, day)
}

/*
	Density 船舶密度统计
*/
func culDensity() {
	lotDivide := 10
	latDivide := 10
	response, err := handler.CulDensity(
		&constant.CulDensityRequest{
			Time: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			DeltaT: &constant.Data{
				Year:   0,
				Month:  0,
				Day:    0,
				Hour:   0,
				Minute: 1,
				Second: 0,
			},
			LotDivide: lotDivide,
			LatDivide: latDivide,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulDensityResponsePrint(response, lotDivide, latDivide)
}

/*
	Speed 船舶航速统计
*/
func culSpeed() {
	lotDivide := 10
	latDivide := 10
	response, err := handler.CulSpeed(
		&constant.CulSpeedRequest{
			Time: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			DeltaT: &constant.Data{
				Year:   0,
				Month:  0,
				Day:    0,
				Hour:   0,
				Minute: 1,
				Second: 0,
			},
			LotDivide: lotDivide,
			LatDivide: latDivide,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulSpeedResponsePrint(response, lotDivide, latDivide)
}

func culDoorLine() {
	response, err := handler.CulDoorLine(
		&constant.CulDoorLineRequest{
			//StartPosition: &constant.Position{
			//	Longitude: 118.04939 - 0.0105,
			//	Latitude:  24.444706 - 0.0035,
			//},
			//EndPosition: &constant.Position{
			//	Longitude: 118.074398 - 0.0105,
			//	Latitude:  24.41378 - 0.0035,
			//},
			//StartPosition: &constant.Position{
			//	Longitude: 118.049497 - 0.0105,
			//	Latitude:  24.451812 - 0.0035,
			//},
			//EndPosition: &constant.Position{
			//	Longitude: 118.064822 - 0.0105,
			//	Latitude:  24.448818 - 0.0035,
			//},
			StartPosition: &constant.Position{
				Longitude: 118.049084 - 0.0105,
				Latitude:  24.444986 - 0.0035,
			},
			EndPosition: &constant.Position{
				Longitude: 118.046353 - 0.0105,
				Latitude:  24.419127 - 0.0035,
			},
			StartTime: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    2,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			EndTime: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    2,
				Hour:   23,
				Minute: 59,
				Second: 59,
			},
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulDoorLineResponsePrint(response)
}

func culSpacing() {
	startT := time.Now()
	response, err := handler.CulSpacing(
		&constant.CulSpacingRequest{
			Time: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			DeltaT: &constant.Data{
				Year:   0,
				Month:  0,
				Day:    0,
				Hour:   0,
				Minute: 1,
				Second: 0,
			},
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
	helper.CulSpacingResponsePrint(response)
}

func culMeeting() {
	startT := time.Now()
	response, err := handler.CulMeeting(
		&constant.CulMeetingRequest{
			StartTime: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			EndTime: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   23,
				Minute: 59,
				Second: 0,
			},
			TimeRange: &constant.Data{
				Year:   0,
				Month:  0,
				Day:    0,
				Hour:   0,
				Minute: 1,
				Second: 0,
			},
			DeltaT: &constant.Data{
				Year:   0,
				Month:  0,
				Day:    0,
				Hour:   0,
				Minute: 1,
				Second: 0,
			},
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
	helper.CulMeetingResponsePrint(response)
}

func earlyWarning() {
	startT := time.Now()
	response, err := handler.EarlyWarning(&constant.EarlyWarningRequest{
		StartTime: &constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 0,
			Second: 0,
		},
		EndTime: &constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 59,
			Second: 0,
		},
		TimeRange: &constant.Data{
			Year:   0,
			Month:  0,
			Day:    0,
			Hour:   0,
			Minute: 1,
			Second: 0,
		},
		DeltaT: &constant.Data{
			Year:   0,
			Month:  0,
			Day:    0,
			Hour:   0,
			Minute: 1,
			Second: 0,
		},
		MMSI: 413694190,
	})
	if err != nil {
		log.Println(err)
		return
	}
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
	helper.EarlyWarningResponsePrint(response)
}

func getTrajectory() {
	startT := time.Now()
	handler.GetTrajectory(&constant.GetTrajectoryRequest{
		MMSI: 412596777,
	})
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
}

func dataSegmentation() {
	startT := time.Now()
	handler.DataSegmentation(&constant.DataSegmentationRequest{
		MMSI: 412596777,
	})
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
}

func dataClean() {
	startT := time.Now()
	handler.DataClean()
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
}

func main() {
	mysql.InitMysql()
	cache.InitCache()
	//culTraffic()
	//culDensity()
	//culSpeed()
	//culDoorLine()
	//culSpacing()
	//culMeeting()
	//earlyWarning()
	//getTrajectory()
	//dataSegmentation()
	dataClean()
}
