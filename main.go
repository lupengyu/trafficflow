package main

import (
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/handler"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

func main() {
	// Traffic 船舶交通量统计
	mysql.InitMysql()
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
	//		},
	//		EndTime: &constant.Data{
	//			Year:   2018,
	//			Month:  12,
	//			Day:    25,
	//			Hour:   23,
	//			Minute: 59,
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
			},
			EndTime: &constant.Data{
				Year:   2019,
				Month:  1,
				Day:    2,
				Hour:   23,
				Minute: 59,
			},
			LotDivide: 10,
			LatDivide: 10,
		},
	)
	if err != nil {
		log.Println(err)
	}
	helper.CulTrafficResponsePrint(response, lotDivide, latDivide, day)
}
