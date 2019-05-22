package main

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/handler"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"os"
	"strconv"
	"time"
)

/*
	Traffic 船舶交通量统计
*/
func culTraffic() {
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
	//			Second: 0,
	//		},
	//		EndTime: &constant.Data{
	//			Year:   2018,
	//			Month:  12,
	//			Day:    25,
	//			Hour:   23,
	//			Minute: 59,
	//			Second: 59,
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
			StartPosition: &constant.Position{
				Longitude: 118.049497 - 0.0105,
				Latitude:  24.451812 - 0.0035,
			},
			EndPosition: &constant.Position{
				Longitude: 118.064822 - 0.0105,
				Latitude:  24.448818 - 0.0035,
			},
			//StartPosition: &constant.Position{
			//	Longitude: 118.049084 - 0.0105,
			//	Latitude:  24.444986 - 0.0035,
			//},
			//EndPosition: &constant.Position{
			//	Longitude: 118.046353 - 0.0105,
			//	Latitude:  24.419127 - 0.0035,
			//},
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

func culNewDoorLine() {
	response, err := handler.CulNewDoorLine(
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
		MMSI: 111333222,
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

func culNewTraffic() {
	startT := time.Now()
	lotDivide := 100
	latDivide := 100
	var day float64 = 12
	response, err := handler.CulNewTraffic(
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
			Day:       int(day),
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	helper.CulTrafficResponsePrint(response, lotDivide, latDivide, day)
	elapsed := time.Since(startT)
	fmt.Println("App elapsed: ", elapsed)
}

func test() {
	StartTime := &constant.Data{
		Year:   2018,
		Month:  12,
		Day:    22,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	EndTime := &constant.Data{
		Year:   2018,
		Month:  12,
		Day:    22,
		Hour:   23,
		Minute: 59,
		Second: 59,
	}
	file, err := os.Create("data/20181222.txt")
	if err != nil {
		log.Println(err)
		return
	}
	file.Sync()
	writer := bufio.NewWriter(file)
	shipIDs, _ := sql.GetShip()
	length := len(shipIDs)
	fmt.Println(shipIDs)
	for index, shipID := range shipIDs {
		// clean and repair
		positions, err := sql.GetNewPositionWithShipIDWithDuration(shipID.MMSI, StartTime, EndTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, pos := range positions {
			str := strconv.Itoa(pos.MMSI) + "," +
				strconv.FormatFloat(pos.SOG, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.COG, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.Longitude, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.Latitude, 'f', -1, 64) + "," +
				strconv.Itoa(pos.Year) + "," +
				strconv.Itoa(pos.Month) + "," +
				strconv.Itoa(pos.Day) + "," +
				strconv.Itoa(pos.Hour) + "," +
				strconv.Itoa(pos.Minute) + "," +
				strconv.Itoa(pos.Second) + "\r\n"
			n, err := writer.WriteString(str)
			if n != len(str) && err != nil {
				log.Println(err)
			}
			writer.Flush()
		}
		percent := float64(100.0*index) / float64(length)
		log.Println("Progress:", percent, "%")
	}
}

func cleanDataBase() {
	rows, err := sql.GetAllPosition()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	preData := &constant.Data{
		Year:   0,
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	index := 0
	for rows.Next() {
		// 计数显示
		if index%10000 == 0 {
			log.Println(index)
		}
		index += 1
		// 数据绑定
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second,
		)
		if err != nil {
			return
		}

		nowData := &constant.Data{
			Year:   pos.Year,
			Month:  pos.Month,
			Day:    pos.Day,
			Hour:   pos.Hour,
			Minute: pos.Minute,
			Second: pos.Second,
		}
		diff := helper.TimeDeviation(nowData, preData)
		if diff > 5*60 {
			fmt.Println("line:", pos.ID, diff/60)
		} else if diff < -300 {
			fmt.Println("line:", pos.ID, diff/60)
		}
		preData = nowData
	}
}

func createData() {
	//rand.Seed(time.Now().Unix())
	file, err := os.Create("data/create.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	file.Sync()
	writer := bufio.NewWriter(file)
	data := &helper.AvailableDataType{
		Longitude: 118.081333,
		Latitude:  24.486428,
		SOG:       0.1,
		COG:       350,
		Length:    110,
	}
	// 加速转弯阶段
	for i := 0; i < 30; i++ {
		data.VMin = 1
		data.VMax = 3
		data.RateTurn = 0.5

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================1")
	for i := 0; i < 30; i++ {
		data.VMin = 3
		data.VMax = 4
		data.RateTurn = 0.1

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		if data.SOG > 10 {
			break
		}
		fmt.Println(data)
	}
	fmt.Println("==========================2")
	// 以下速度控制
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = 0.05
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================3")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = 0.02
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================4")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = 0.02
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================5")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = -0.01
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================5")
	for i := 0; i < 20; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = -0.05
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================6")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = -0.01
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================7")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = 0.02
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================8")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = 0.02
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 12 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================9")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = 0.015
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 16 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================10")
	for i := 0; i < 30; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = 0.01
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 14 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================11")
	for i := 0; i < 50; i++ {
		data.VMin = -1
		data.VMax = 1
		data.RateTurn = 0.0
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 14 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================12")
	for i := 0; i < 50; i++ {
		data.VMin = -1
		data.VMax = 2
		data.RateTurn = -0.013
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 14 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================13")
	for i := 0; i < 200; i++ {
		data.VMin = -1
		data.VMax = 3
		data.RateTurn = -0.0
		if data.SOG < 8 {
			data.VMin = 1
			data.VMax = 3
		} else if data.SOG > 20 {
			data.VMin = -2
			data.VMax = 0
		}

		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		data = helper.AvailableDataTest(data, 10)
		fmt.Println(data)
	}
	fmt.Println("==========================14")
}

func dataAvailable() {
	handler.DataAvailable("create.txt")
	handler.ZhangDataAvailable("create.txt")
}

func createRawData() {
	handler.CreateRawData("create.txt")
}

func cleanRawData() {
	handler.CleanRawData("small/rawdata.txt")
}

func repairCleanData() {
	handler.RepairCleanData("small/cleandata.txt")
}

func zhangCleanRawData() {
	handler.ZhangCleanRawData("small/rawdata.txt")
}

func getShipRawData() {
	file, err := os.Create("data/412596777.txt")
	if err != nil {
		log.Println(err)
		return
	}
	file.Sync()
	writer := bufio.NewWriter(file)
	positions, err := sql.GetPositionWithShipID(412596777)
	if err != nil {
		log.Println("查询失败")
		return
	}
	for _, v := range positions {
		nowTime := &constant.Data{
			Year:   v.Year,
			Month:  v.Month,
			Day:    v.Day,
			Hour:   v.Hour,
			Minute: v.Minute,
			Second: v.Second,
		}
		preTime := &constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 0,
			Second: 0,
		}
		diff := helper.TimeDeviation(nowTime, preTime)
		str := strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(v.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(v.COG, 'f', -1, 64) + "," + strconv.FormatInt(diff, 10)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()
	}
}

func culDeviation() {
	//handler.CulDeviation("create(true).txt", "cleandata.txt", "cleandata_repair.txt")
	handler.CulDeviation("small/create.txt", "small/cleandata.txt", "small/cleandata_repair.txt")
}

func soleimani() {
	handler.Soleimani("small/create.txt", "small/testfile.txt")
}

func main() {
	mysql.InitMysql()
	cache.InitCache()
	//cleanDataBase()
	//test()

	//culTraffic()
	//culNewTraffic()

	//culDensity()

	//culSpeed()

	//culDoorLine()
	//culNewDoorLine()

	//culSpacing()

	culMeeting()

	//earlyWarning()

	//getTrajectory()

	//dataSegmentation()

	//dataClean()

	//createData()
	//createRawData()
	//cleanRawData() // 清洗原始数据
	//repairCleanData()
	zhangCleanRawData()
	//dataAvailable()
	//culDeviation()

	//getShipRawData()

	//soleimani()
}
