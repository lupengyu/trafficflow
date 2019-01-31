package handler

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

func CulSpeed(request *constant.CulSpeedRequest) (response *constant.CulSpeedResponse, err error) {
	beginTime := helper.DayDecrease(request.Time, request.DeltaT)
	endTime := helper.DayIncrease(request.Time, request.DeltaT)
	// 查询时间段内的数据
	rows, err := sql.GetPositionWithDuration(beginTime, endTime)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	shipSpeedRangeMap := make([]map[int]int, 5)
	for i := 0; i < 5; i += 1 {
		shipSpeedRangeMap[i] = make(map[int]int)
	}
	shipSpeedRange := make([]int, 5)
	shipSpeedSumMap := make(map[int]float64)
	shipSpeedCnt := make(map[int]int)
	var areaSpeed [][]constant.AreaSpeed
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaSpeed, request.LatDivide)
		areaSpeed = append(areaSpeed, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			areaSpeed[i][j].ShipSpeedSumMap = make(map[int]float64)
			areaSpeed[i][j].ShipSpeedCnt = make(map[int]int)
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
		// 移除时间错误的数据
		if pos.Hour > 23 {
			continue
		}
		// 移除航速不可用的数据
		if pos.SOG >= 102.3 {
			continue
		}

		shipSpeedSumMap[pos.MMSI] += pos.SOG
		shipSpeedCnt[pos.MMSI] += 1
		areaSpeed[longitudeArea][latitudeArea].ShipSpeedSumMap[pos.MMSI] += pos.SOG
		areaSpeed[longitudeArea][latitudeArea].ShipSpeedCnt[pos.MMSI] += 1
	}

	shipCnt := len(shipSpeedSumMap)
	var shipSpeedSum float64 = 0
	for k, v := range shipSpeedSumMap {
		speedAverage := v / float64(shipSpeedCnt[k])
		shipSpeedSum += speedAverage
		speedRange := helper.SpeedRange(speedAverage)
		shipSpeedRange[speedRange] += 1
	}

	for i := 0; i < request.LotDivide; i += 1 {
		for j := 0; j < request.LatDivide; j += 1 {
			areaSpeed[i][j].ShipCnt = len(areaSpeed[i][j].ShipSpeedSumMap)
			for k, v := range areaSpeed[i][j].ShipSpeedSumMap {
				speedAverage := v / float64(areaSpeed[i][j].ShipSpeedCnt[k])
				areaSpeed[i][j].ShipSpeed += speedAverage
			}
			areaSpeed[i][j].ShipSpeed /= float64(areaSpeed[i][j].ShipCnt)
		}
	}

	return &constant.CulSpeedResponse{
		SpeedData: &constant.SpeedData{
			ShipSpeed:      shipSpeedSum / float64(shipCnt),
			ShipCnt:        shipCnt,
			ShipSpeedRange: shipSpeedRange,
		},
		AreaSpeed: areaSpeed,
	}, nil
}
