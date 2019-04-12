package handler

import (
	"container/list"
	"fmt"
	"github.com/cnkei/gospline"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"math"
	"time"
)

/*
	清洗程序
*/
func DataClean() {
	positions, err := sql.GetPositionWithShipID(412596777)
	if err != nil {
		log.Println("查询失败")
		return
	}
	length := len(positions)
	if length == 0 {
		return
	}
	ignore := make([]int, 0)
	ends := make([]int, 0)
	signal := true
	prePosition := constant.PositionMeta{}
	index := 0

	for index < length {
		start := index
		for ; start < length; start++ {
			if signal {
				prePosition = positions[start]
				signal = false
				continue
			}
			rg := helper.TimeDeviation(&constant.Data{
				Year:   positions[start].Year,
				Month:  positions[start].Month,
				Day:    positions[start].Day,
				Hour:   positions[start].Hour,
				Minute: positions[start].Minute,
				Second: positions[start].Second,
			}, &constant.Data{
				Year:   prePosition.Year,
				Month:  prePosition.Month,
				Day:    prePosition.Day,
				Hour:   prePosition.Hour,
				Minute: prePosition.Minute,
				Second: prePosition.Second,
			})
			if rg <= 0 {
				// 忽略时间错位数据
				ignore = append(ignore, start)
			} else {
				prePosition = positions[start]
				if rg > 5*60 {
					// 间隔时间大于5分钟，区分之
					ends = append(ends, start)
					signal = true
					break
				}
			}
		}
		index = start
	}
	// 分段清洗
	start := 0
	cnt := 1
	//log.Println(ignore)
	//log.Println(ends)
	ignoreIndex := 0
	for _, v := range ends {
		dataList := make([]constant.PositionMeta, 0)
		for i := start; i < v; i++ {
			if ignoreIndex < len(ignore) && i == ignore[ignoreIndex] {
				ignoreIndex += 1
				continue
			}
			dataList = append(dataList, positions[i])
		}
		response, err := CleaningAndRepairPositionMeta(&constant.CleaningAndRepairPositionMetaRequest{
			DataList: dataList,
		})
		if err != nil {
			log.Println(err)
			return
		}
		SaveCleanData(&constant.SaveCleanDataRequest{
			DataList: response.DataList,
		})
		return // TODO: remove return
		start = v
		cnt += 1
	}
	if start < length {
		dataList := make([]constant.PositionMeta, 0)
		for i := start; i < length; i++ {
			if ignoreIndex < len(ignore) && i == ignore[ignoreIndex] {
				ignoreIndex += 1
				continue
			}
			dataList = append(dataList, positions[i])
		}
		response, err := CleaningAndRepairPositionMeta(&constant.CleaningAndRepairPositionMetaRequest{
			DataList: dataList,
		})
		if err != nil {
			log.Println(err)
			return
		}
		SaveCleanData(&constant.SaveCleanDataRequest{
			DataList: response.DataList,
		})
	}
}

/*
	清洗和修复数据
*/
func CleaningAndRepairPositionMeta(request *constant.CleaningAndRepairPositionMetaRequest) (*constant.CleaningAndRepairPositionMetaResponse, error) {
	fmt.Println(request.DataList)
	if len(request.DataList) == 1 {
		return &constant.CleaningAndRepairPositionMetaResponse{
			DataList: request.DataList,
		}, nil
	}
	// 数据清洗
	rawData := request.DataList
	length := len(rawData)
	ignore := make([]int, 0)
	start := 0
	prePosition := rawData[0]
	cleanData := []constant.PositionMeta{prePosition}
	// 起点适用性判断 仅当长度大于等于3时判断
	if length >= 3 && helper.PositionAvailable(rawData[1], rawData[0]) == false {
		if helper.PositionAvailable(rawData[2], rawData[1]) == true {
			// 第一个点为异常点，抛弃
			start = 1
			ignore = append(ignore, 0)
			prePosition = rawData[1]
			cleanData = []constant.PositionMeta{prePosition}
		}
	}
	for i := start + 1; i < length; i++ {
		if helper.PositionAvailable(rawData[i], prePosition) == false {
			ignore = append(ignore, i)
			continue
		}
		prePosition = rawData[i]
		cleanData = append(cleanData, prePosition)
	}

	helper.PointListOutput("clean", cleanData)

	// 数据修复
	repairData := make([]constant.PositionMeta, 0)
	beforeList := list.New() // 前队列
	prePosition = cleanData[0]
	cleanLength := len(cleanData)
	beforeList.PushBack(prePosition)
	repairData = append(repairData, prePosition)
	for i := 1; i < cleanLength; i++ {
		nowPosition := cleanData[i]
		nowTime := &constant.Data{
			Year:   nowPosition.Year,
			Month:  nowPosition.Month,
			Day:    nowPosition.Day,
			Hour:   nowPosition.Hour,
			Minute: nowPosition.Minute,
			Second: nowPosition.Second,
		}
		preTime := &constant.Data{
			Year:   prePosition.Year,
			Month:  prePosition.Month,
			Day:    prePosition.Day,
			Hour:   prePosition.Hour,
			Minute: prePosition.Minute,
			Second: prePosition.Second,
		}
		diff := helper.TimeDeviation(nowTime, preTime)
		if diff <= 30 {
			if beforeList.Len() == 3 {
				beforeList.Remove(beforeList.Front())
			}
			beforeList.PushBack(nowPosition)
			prePosition = nowPosition
			repairData = append(repairData, prePosition)
		} else {
			afterList := list.New()
			for start := i; start < cleanLength; start ++ {
				if afterList.Len() == 3 {
					break
				}
				afterList.PushBack(cleanData[start])
			}

			x := make([]float64, 0)
			longitudeY := make([]float64, 0)
			latitudeY := make([]float64, 0)
			cogY := make([]float64, 0)
			sogY := make([]float64, 0)
			preCOG := beforeList.Front().Value.(constant.PositionMeta).COG

			for e := beforeList.Front(); e != nil; e = e.Next() {
				position := e.Value.(constant.PositionMeta)
				diffX := helper.TimeDeviation(&constant.Data{
					Year:   position.Year,
					Month:  position.Month,
					Day:    position.Day,
					Hour:   position.Hour,
					Minute: position.Minute,
					Second: position.Second,
				}, preTime)

				x = append(x, float64(diffX))
				longitudeY = append(longitudeY, position.Longitude)
				latitudeY = append(latitudeY, position.Latitude)
				cogY = append(cogY, helper.RateRange(position.COG, preCOG))
				sogY = append(sogY, position.SOG)
			}

			for e := afterList.Front(); e != nil; e = e.Next() {
				position := e.Value.(constant.PositionMeta)
				diffX := helper.TimeDeviation(&constant.Data{
					Year:   position.Year,
					Month:  position.Month,
					Day:    position.Day,
					Hour:   position.Hour,
					Minute: position.Minute,
					Second: position.Second,
				}, preTime)

				x = append(x, float64(diffX))
				longitudeY = append(longitudeY, position.Longitude)
				latitudeY = append(latitudeY, position.Latitude)
				cogY = append(cogY, helper.RateRange(position.COG, preCOG))
				sogY = append(sogY, position.SOG)
			}

			fmt.Println("x", x)
			fmt.Println("longitudeY", longitudeY)
			fmt.Println("latitudeY", latitudeY)
			fmt.Println("cogY", cogY)
			fmt.Println("sogY", sogY)

			longitudeFunc := gospline.NewCubicSpline(x, longitudeY)
			latitudeFunc := gospline.NewCubicSpline(x, latitudeY)
			cogFunc := gospline.NewCubicSpline(x, cogY)
			sogFunc := gospline.NewCubicSpline(x, sogY)

			need := (int(diff) - 1) / 30
			baseTime := time.Date(preTime.Year, time.Month(preTime.Month), preTime.Day, preTime.Hour, preTime.Minute, preTime.Second, 0, time.UTC)
			fmt.Println("=============")
			fmt.Println(prePosition)
			fmt.Println(nowPosition)
			for j := 0; j < need; j++ {
				add := (j + 1) * int(diff) / (need + 1)
				resultTime := baseTime.Add(time.Duration(add)*time.Second)
				fmt.Println(baseTime, resultTime, add)
				longitude := longitudeFunc.At(float64(add))
				latitude := latitudeFunc.At(float64(add))
				cog := cogFunc.At(float64(add)) + preCOG
				sog := sogFunc.At(float64(add))
				sog = math.Max(0, sog)
				if cog < 0 {
					cog = cog - float64(int(cog/360)-1)*360.0
				} else if cog > 360 {
					cog = cog - float64(int(cog/360))*360.0
				}
				position := constant.PositionMeta{
					MMSI: prePosition.MMSI,
					SOG:sog,
					COG:cog,
					Longitude:longitude,
					Latitude:latitude,
					Year:   resultTime.Year(),
					Month:  int(resultTime.Month()),
					Day:    resultTime.Day(),
					Hour:   resultTime.Hour(),
					Minute: resultTime.Minute(),
					Second: resultTime.Second(),
				}
				if beforeList.Len() == 3 {
					beforeList.Remove(beforeList.Front())
				}
				beforeList.PushBack(position)
				prePosition = position
				repairData = append(repairData, prePosition)
				fmt.Println(position)
			}
		}
	}

	helper.PointListOutput("repair", repairData)

	return &constant.CleaningAndRepairPositionMetaResponse{
		DataList: repairData,
	}, nil
}

/*
	新数据存储
*/
func SaveCleanData(request *constant.SaveCleanDataRequest) {

}
