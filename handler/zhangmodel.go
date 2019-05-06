package handler

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func movingAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	if prePosition.Longitude == position.Longitude &&
		prePosition.Latitude == position.Latitude &&
		prePosition.SOG > 2 {
		fmt.Println("moving", prePosition, position) // TODO: remove
		return false
	}
	return true
}

func speedAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	if position.SOG > 16 {
		//fmt.Println("speed", prePosition.SOG, 16, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func driftAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	diff := helper.TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	preV := (16 * 1.852) / 3.6
	D := helper.PositionSpacing(&constant.Position{
		Latitude:  prePosition.Latitude,
		Longitude: prePosition.Longitude,
	}, &constant.Position{
		Latitude:  position.Latitude,
		Longitude: position.Longitude,
	})
	acturalV := D / float64(diff)
	if acturalV > preV {
		fmt.Println("drift", acturalV, preV, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func accelerationAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	diff := helper.TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	shipInfo := cache.GetShipInfo(position.MMSI)
	a := helper.MaxAcceleration(shipInfo.Length, 16.0)
	if position.SOG > prePosition.SOG+float64(diff)*a {
		fmt.Println("acceleration", position.SOG, prePosition.SOG+float64(diff)*a, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func rateAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	rate := math.Max(position.COG-prePosition.COG, prePosition.COG-position.COG)
	absRate := math.Min(rate, 360.0-rate)
	diff := helper.TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	shipInfo := cache.GetShipInfo(position.MMSI)
	Vnm := (prePosition.SOG + position.SOG) / 2
	maxRate := helper.MaxRate(shipInfo.Length, Vnm)
	if absRate > float64(diff)*maxRate {
		fmt.Println("rate", absRate, float64(diff)*maxRate, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func PositionAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	return movingAvailable(position, prePosition) && speedAvailable(position, prePosition) &&
		driftAvailable(position, prePosition) && accelerationAvailable(position, prePosition) &&
		rateAvailable(position, prePosition)
}

func ZhangDataAvailable(fileName string) {
	file, err := os.Open("data/" + fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	rawData := make([]constant.PositionMeta, 0)
	bfRd := bufio.NewReader(file)
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		sog, _ := strconv.ParseFloat(strs[2], 64)
		cog, _ := strconv.ParseFloat(strs[3], 64)
		rawData = append(rawData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
			SOG:       sog,
			COG:       cog,
			Year:      2019,
			Month:     5,
			Day:       1,
			Hour:      0,
			Minute:    0,
			Second:    0,
			MMSI:      99999999,
		})
	}
	preData := rawData[0]
	tmp := 0
	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		preData.Second = 0
		data.Second = 10 + tmp
		//fmt.Println(preData)
		//fmt.Println(data)
		if PositionAvailable(data, preData) == false {
			//fmt.Println(i)
			tmp += 10
			continue
		} else {
			tmp = 0
		}
		preData = data
	}
	fmt.Println("All Data Pass")
}
