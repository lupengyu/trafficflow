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
func culTraffice() {
	lotDivide := 10
	latDivide := 10
	//var day float64 = 1
	//response, err := handler.CulTraffic(
	//	&constant.CulTrafficRequest{
	//		StartTime: &constant.Data{
	//			Year:   2018,
	//			Month:  12,
	//			Day:    25,
	//			Hour:   0,
	//			Minute: 0,
	//          Second: 0,
	//		},
	//		EndTime: &constant.Data{
	//			Year:   2018,
	//			Month:  12,
	//			Day:    25,
	//			Hour:   23,
	//			Minute: 59,
	//          Second: 59,
	//		},
	//		LotDivide: lotDivide,
	//		LatDivide: latDivide,
	//	},
	//)
	var day float64 = 12
	response, err := handler.CulTraffic(
		&constant.CulTrafficRequest{
			StartTime: &constant.Data{
				Year:   2018,
				Month:  12,
				Day:    22,
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
			LotDivide: lotDivide,
			LatDivide: latDivide,
		},
	)
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
			Time:      &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    1,
				Hour:   0,
				Minute: 0,
				Second: 0,
			},
			DeltaT:    &constant.Data{
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

func main() {
	mysql.InitMysql()
	culDensity()
}
