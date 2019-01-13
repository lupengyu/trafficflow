package main

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/handler"
)

func main() {
	mysql.InitMysql()
	response, err := handler.CulTraffice(
		&constant.CulTrafficeRequest{
			StartTime: 	&constant.Data{
				Year: 	2018,
				Month: 	12,
				Day: 	25,
				Hour: 	0,
				Minute: 0,
			},
			EndTime: 	&constant.Data{
				Year: 	2018,
				Month: 	12,
				Day: 	25,
				Hour: 	23,
				Minute: 59,
			},
			LotDivide: 	10,
			LatDivide: 	10,
		},
	)
	fmt.Println(response, err)
}