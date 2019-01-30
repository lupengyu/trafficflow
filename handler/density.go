package handler

import (
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

/*
	计算船舶密度
*/
func CulDensity(request *constant.CulDensityRequest) (response *constant.CulDensityResponse, err error) {
	beginTime := helper.DayDecrease(request.Time, request.DeltaT)
	endTime := helper.DayIncrease(request.Time, request.DeltaT)
	// 查询时间段内的数据
	rows, err := mysql.DB.Query(
		"select * from position where (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?))))) order by id",
		beginTime.Year, beginTime.Year,
		beginTime.Month, beginTime.Month,
		beginTime.Day, beginTime.Day,
		beginTime.Hour, beginTime.Hour,
		beginTime.Minute, beginTime.Minute,
		beginTime.Second,
		endTime.Year, endTime.Year,
		endTime.Month, endTime.Month,
		endTime.Day, endTime.Day,
		endTime.Hour, endTime.Hour,
		endTime.Minute, endTime.Minute,
		endTime.Second,
	)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	shipMap := make(map[int]int)
	shipDensity := 0
	var areaDensity [][]constant.AreaDensity
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaDensity, request.LatDivide)
		areaDensity = append(areaDensity, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			// ship
			areaDensity[i][j].ShipMap = make(map[int]int)
		}
	}

	// 数据循环遍历计算
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

		// shipMap记录
		if shipMap[pos.MMSI] == 0 {
			shipMap[pos.MMSI] = 1
			shipDensity += 1
		}

		// 数据记录
		for i := 0; i < request.LotDivide; i += 1 {
			for j := 0; j < request.LatDivide; j += 1 {
				if i == longitudeArea && j == latitudeArea {
					// 添加本区域的记录
					if areaDensity[i][j].ShipMap[pos.MMSI] == 0 {
						areaDensity[i][j].ShipMap[pos.MMSI] = 1
						areaDensity[i][j].Density += 1
					}
				} else {
					// 清除其他区域的记录
					if areaDensity[i][j].ShipMap[pos.MMSI] == 1 {
						areaDensity[i][j].ShipMap[pos.MMSI] = 0
						areaDensity[i][j].Density -= 1
					}
				}
			}
		}
	}

	return &constant.CulDensityResponse{
		Density:     shipDensity,
		AreaDensity: areaDensity,
	}, nil
}
