package handler

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/helper"
)

func CulTraffice(request *constant.CulTrafficeRequest) (response *constant.CulTrafficeResponse, err error) {
	rows, err := mysql.DB.Query(
		"select * from position where year >= ? and year <= ? " + "and month >= ? and month <= ? and day >= ? and day <= ?",
			request.StartTime.Year, request.EndTime.Year,
			request.StartTime.Month, request.EndTime.Month,
			request.StartTime.Day, request.EndTime.Day)
	if err != nil {
		return nil, err
	}
	var areaTraffic [][]constant.AreaTraffic
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaTraffic, request.LatDivide)
		areaTraffic = append(areaTraffic, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			areaTraffic[i][j].ShipMap = make(map[int]int)
		}
	}
	for rows.Next() {
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second)
		if err != nil {
			return nil, err
		}
		longitudeArea := helper.LongitudeArea(pos.Longitude, request.LotDivide)
		latitudeArea := helper.LatitudeArea(pos.Latitude, request.LatDivide)
		if longitudeArea == -1 || latitudeArea == -1 {
			continue
		}
		// 总记录
		for i := 0; i < request.LotDivide; i += 1 {
			for j := 0; j < request.LatDivide; j += 1 {
				if i == longitudeArea && j == latitudeArea {
					// 添加本区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 0 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 1
						areaTraffic[i][j].Traffic += 1
					}
				} else {
					// 清除其他区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 1 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 1
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
	fmt.Println("====================")
	fmt.Println("Latitude", constant.LatitudeMin, "-", constant.LatitudeMax)
	fmt.Println("Longitude", constant.LongitudeMin, "-", constant.LongitudeMax)
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Traffic, "\t")
		}
		fmt.Print("\n")
	}
	return nil, nil
}