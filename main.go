package main

import (
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/handler"
	"github.com/lupengyu/trafficflow/helper"
	"log"
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
			//	Longitude: 118.04939,
			//	Latitude:  24.444706,
			//},
			//EndPosition: &constant.Position{
			//	Longitude: 118.074398,
			//	Latitude:  24.41378,
			//},
			//StartPosition: &constant.Position{
			//	Longitude: 118.124272,
			//	Latitude:  24.244077,
			//},
			//EndPosition: &constant.Position{
			//	Longitude: 118.319528,
			//	Latitude:  24.393509,
			//},
			StartPosition: &constant.Position{
				Longitude: 118.344824,
				Latitude:  24.393509,
			},
			EndPosition: &constant.Position{
				Longitude: 118.593188,
				Latitude:  24.221604,
			},
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
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulDoorLineResponsePrint(response)
}

func culSpacing() {
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
	helper.CulSpacingResponsePrint(response)
}

func main() {
	mysql.InitMysql()
	//culTraffic()
	//culDensity()
	//culSpeed()
	//culDoorLine()
	//culSpacing()
}
