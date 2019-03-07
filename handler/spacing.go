package handler

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

/*
	计算船舶间距
*/
func CulSpacing(request *constant.CulSpacingRequest) (response *constant.CulSpacingResponse, err error) {
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

	// 数据初始化
	trackMap := make(map[int]*constant.Track)
	shipTrackList := make(map[int][]*constant.Track)

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

		// 判断船舶位置
		nowPosition := &constant.Position{
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
		}
		nowTime := &constant.Data{
			Year:   pos.Year,
			Month:  pos.Month,
			Day:    pos.Day,
			Hour:   pos.Hour,
			Minute: pos.Minute,
			Second: pos.Second,
		}
		shipTrackList[pos.MMSI] = append(shipTrackList[pos.MMSI], &constant.Track{
			PrePosition: nowPosition,
			Time:        nowTime,
			Deviation:   helper.TimeDeviation(nowTime, request.Time),
			COG:         pos.COG,
		})
	}

	shipSpacing := make(map[int]map[int]float64)

	spacing := make(map[int]float64)
	for k, v := range shipTrackList {
		trackMap[k] = helper.TrackDifference(v)
		spacing[k] = 9999999999
	}

	var minSpacing float64 = 9999999999
	minSpaceA := 0
	minSpaceB := 0
	aPosition := &constant.Position{}
	bPosition := &constant.Position{}

	size := len(trackMap)
	tracks := make([]*constant.Track, size)
	shipIDs := make([]int, size)
	index := 0
	for k, v := range trackMap {
		shipSpacing[k] = make(map[int]float64)
		tracks[index] = v
		shipIDs[index] = k
		index++
	}

	for i := 0; i < size; i++ {
		k1 := shipIDs[i]
		v1 := tracks[i]
		for j := i + 1; j < size; j++ {
			k2 := shipIDs[j]
			v2 := tracks[j]
			nowSpacing := helper.PositionSpacing(v1.PrePosition, v2.PrePosition)
			shipSpacing[k1][k2] = nowSpacing
			shipSpacing[k2][k1] = nowSpacing
			if nowSpacing < spacing[k1] {
				spacing[k1] = nowSpacing
			}
			if nowSpacing < spacing[k2] {
				spacing[k2] = nowSpacing
			}
			if nowSpacing < minSpacing {
				minSpacing = nowSpacing
				minSpaceA = k1
				minSpaceB = k2
				aPosition = v1.PrePosition
				bPosition = v2.PrePosition
			}
		}
	}

	spacingRange := make([]int, 3)
	for _, v := range spacing {
		if v < 50 {
			spacingRange[0] += 1
		} else if v > 300 {
			spacingRange[2] += 1
		} else {
			spacingRange[1] += 1
		}
	}

	// 输出结果
	return &constant.CulSpacingResponse{
		MinSpacing:   minSpacing,
		MinSpaceA:    minSpaceA,
		MinSpaceB:    minSpaceB,
		APosition:    aPosition,
		BPosition:    bPosition,
		SpacingMap:   spacing,
		SpacingRange: spacingRange,
		ShipSpacing:  shipSpacing,
		TrackMap:     trackMap,
	}, nil
}
