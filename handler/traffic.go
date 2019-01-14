package handler

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/helper"
)

func CulTraffice(request *constant.CulTrafficRequest) (response *constant.CulTrafficResponse, err error) {
	rows, err := mysql.DB.Query(
	"select * from position where (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?)))))",
		request.StartTime.Year, request.StartTime.Year,
		request.StartTime.Month, request.StartTime.Month,
		request.StartTime.Day, request.StartTime.Day,
		request.StartTime.Hour, request.StartTime.Hour,
		request.StartTime.Minute, request.StartTime.Minute,
		request.StartTime.Second,
		request.EndTime.Year, request.EndTime.Year,
		request.EndTime.Month, request.EndTime.Month,
		request.EndTime.Day, request.EndTime.Day,
		request.EndTime.Hour, request.EndTime.Hour,
		request.EndTime.Minute, request.EndTime.Minute,
		request.EndTime.Second,
	)
	if err != nil {
		return nil, err
	}
	preTime := &constant.Data {
		Year: 	0,
		Month: 	0,
		Day: 	0,
	}
	trafficData := &constant.TrafficeData {
		HourTrafficSum: 	make([]int, 24),
	}
	var areaTraffic [][]constant.AreaTraffic
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaTraffic, request.LatDivide)
		areaTraffic = append(areaTraffic, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			areaTraffic[i][j].ShipMap = make(map[int]int)
			areaTraffic[i][j].HourShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourTraffic = make([]int, 24)
		}
	}
	for rows.Next() {
		// 数据绑定
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second,
		)
		if err != nil {
			return nil, err
		}

		/*
			计算当前船只经纬度所处地图分块
			对于不在地图经纬度区域内的数据剔除
		 */
		longitudeArea := helper.LongitudeArea(pos.Longitude, request.LotDivide)
		latitudeArea := helper.LatitudeArea(pos.Latitude, request.LatDivide)
		if longitudeArea == -1 || latitudeArea == -1 {
			continue
		}
		if pos.Hour > 23 {
			continue
		}
		nowTime := &constant.Data {
			Year: 	pos.Year,
			Month: 	pos.Month,
			Day: 	pos.Day,
			Hour: 	pos.Hour,
		}

		/*
			判断日期有没有刷新
				有：刷新记录
		 */
		if !helper.DayEqual(preTime, nowTime) {
			preTime = nowTime
			for i := 0; i < request.LotDivide; i += 1 {
				for j := 0; j < request.LatDivide; j += 1 {
					areaTraffic[i][j].ShipMap = make(map[int]int)
					for k := 0; k < 24; k += 1 {
						areaTraffic[i][j].HourShipMap[k] = make(map[int]int)
					}
				}
			}
		}

		// 数据记录
		for i := 0; i < request.LotDivide; i += 1 {
			for j := 0; j < request.LatDivide; j += 1 {
				if i == longitudeArea && j == latitudeArea {
					// 添加本区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 0 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 1
						areaTraffic[i][j].Traffic += 1
					}
					if areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] == 0 {
						areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] = 1
						areaTraffic[i][j].HourTraffic[nowTime.Hour] += 1
						trafficData.HourTrafficSum[nowTime.Hour] += 1
					}
				} else {
					// 清除其他区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 1 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 0
					}
					if areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] == 1 {
						areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] = 0
					}
				}
			}
		}
	}
	//for i := 0; i < request.LotDivide; i += 1 {
	//	for j := 0; j < request.LatDivide; j += 1 {
	//		fmt.Println(i, j, ": ")
	//		fmt.Println(areaTraffic[i][j].ShipMap)
	//		fmt.Println("Traffic:", areaTraffic[i][j].Traffic)
	//	}
	//}
	//fmt.Println("Latitude", constant.LatitudeMin, "-", constant.LatitudeMax)
	//fmt.Println("Longitude", constant.LongitudeMin, "-", constant.LongitudeMax)
	fmt.Println("=========DayTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Traffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========HourTraffic=========")
	for i := 0; i < 24; i += 1 {
		fmt.Println(i, ":", trafficData.HourTrafficSum[i])
	}
	fmt.Println("=========AreaTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Println(i, ",", j, ":")
			fmt.Println(areaTraffic[i][j].Traffic, areaTraffic[i][j].HourTraffic)
		}
	}
	return nil, nil
}