package handler

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/cnkei/gospline"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateRawData(fileName string) {
	doorLine, err := os.Create("data/rawdata.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)

	file, err := os.Open("data/" + fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	bfRd := bufio.NewReader(file)
	protect := false
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			protect = true
			str := string(line)
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			continue
		}
		if protect == true {
			protect = false
			str := string(line)
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			continue
		}
		str := "XXXXXXXXX"
		if rand.Intn(100) > 20 {
			str = string(line)
		}
		n, err := doorLineWriter.WriteString(str + "\r\n")
		if n < 0 && err != nil {
			log.Println(err)
		}
		doorLineWriter.Flush()
	}
}

func CleanRawData(fileName string) {
	doorLine, err := os.Create("data/small/cleandata.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)

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
			if strs[0] == "XXXXXXXXX" {
				rawData = append(rawData, constant.PositionMeta{
					MMSI: -1,
				})
			}
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
	str := strconv.FormatFloat(preData.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(preData.Latitude, 'f', -1, 64) +
		"," + strconv.FormatFloat(preData.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(preData.COG, 'f', -1, 64)
	n, err := doorLineWriter.WriteString(str + "\r\n")
	if n < 0 && err != nil {
		log.Println(err)
	}
	doorLineWriter.Flush()
	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		if data.MMSI == -1 {
			str := "XXXXXXXXX"
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp += 10
			continue
		}
		preData.Second = 0
		data.Second = 10 + tmp
		if helper.PositionAvailable(data, preData) == false {
			str := "XXXXXXXXX"
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp += 10
			continue
		} else {
			str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
				"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp = 0
		}
		preData = data
	}
	fmt.Println("All Data Pass")
}

func ZhangCleanRawData(fileName string) {
	doorLine, err := os.Create("data/small/zhangcleandata.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)

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
			if strs[0] == "XXXXXXXXX" {
				rawData = append(rawData, constant.PositionMeta{
					MMSI: -1,
				})
			}
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
	str := strconv.FormatFloat(preData.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(preData.Latitude, 'f', -1, 64) +
		"," + strconv.FormatFloat(preData.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(preData.COG, 'f', -1, 64)
	n, err := doorLineWriter.WriteString(str + "\r\n")
	if n < 0 && err != nil {
		log.Println(err)
	}
	doorLineWriter.Flush()
	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		if data.MMSI == -1 {
			str := "XXXXXXXXX"
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp += 10
			continue
		}
		preData.Second = 0
		data.Second = 10 + tmp
		if PositionAvailable(data, preData) == false {
			str := "XXXXXXXXX"
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp += 10
			continue
		} else {
			str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64) +
				"," + strconv.FormatFloat(data.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(data.COG, 'f', -1, 64)
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n < 0 && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
			tmp = 0
		}
		preData = data
	}
	fmt.Println("All Data Pass")
}

func RepairCleanData(fileName string) {
	doorLine, err := os.Create("data/small/cleandata_repair.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)

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
	baseTime := time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	index := 0
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		duration := time.Duration(10*index) * time.Second
		index += 1
		resultTime := baseTime.Add(duration)
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			if strs[0] == "XXXXXXXXX" {
				rawData = append(rawData, constant.PositionMeta{
					MMSI:   -1,
					Year:   resultTime.Year(),
					Month:  int(resultTime.Month()),
					Day:    resultTime.Day(),
					Hour:   resultTime.Hour(),
					Minute: resultTime.Minute(),
					Second: resultTime.Second(),
				})
			}
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
			Year:      resultTime.Year(),
			Month:     int(resultTime.Month()),
			Day:       resultTime.Day(),
			Hour:      resultTime.Hour(),
			Minute:    resultTime.Minute(),
			Second:    resultTime.Second(),
		})
	}

	cleanLength := len(rawData)
	beforeList := list.New()
	beforeList.PushBack(rawData[0])

	for i := 1; i < cleanLength; i++ {
		nowPosition := rawData[i]
		nowTime := &constant.Data{
			Year:   nowPosition.Year,
			Month:  nowPosition.Month,
			Day:    nowPosition.Day,
			Hour:   nowPosition.Hour,
			Minute: nowPosition.Minute,
			Second: nowPosition.Second,
		}
		if nowPosition.MMSI != -1 {
			if beforeList.Len() == 3 {
				beforeList.Remove(beforeList.Front())
			}
			beforeList.PushBack(nowPosition)
		} else {
			afterList := list.New()
			for start := i + 1; start < cleanLength; start++ {
				if afterList.Len() == 3 {
					break
				}
				if rawData[start].MMSI == -1 {
					continue
				}
				afterList.PushBack(rawData[start])
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
				}, nowTime)

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
				}, nowTime)

				x = append(x, float64(diffX))
				longitudeY = append(longitudeY, position.Longitude)
				latitudeY = append(latitudeY, position.Latitude)
				cogY = append(cogY, helper.RateRange(position.COG, preCOG))
				sogY = append(sogY, position.SOG)
			}

			longitudeFunc := gospline.NewCubicSpline(x, longitudeY)
			latitudeFunc := gospline.NewCubicSpline(x, latitudeY)
			cogFunc := gospline.NewCubicSpline(x, cogY)
			sogFunc := gospline.NewCubicSpline(x, sogY)

			cog := cogFunc.At(0) + preCOG
			sog := sogFunc.At(0)
			sog = math.Max(0, sog)
			if cog < 0 {
				cog = cog - float64(int(cog/360)-1)*360.0
			} else if cog > 360 {
				cog = cog - float64(int(cog/360))*360.0
			}

			rawData[i].Latitude = latitudeFunc.At(0)
			rawData[i].Longitude = longitudeFunc.At(0)
			rawData[i].COG = cog
			rawData[i].SOG = sog

			//itself
			//if beforeList.Len() == 3 {
			//	beforeList.Remove(beforeList.Front())
			//}
			//beforeList.PushBack(rawData[i])
		}
	}

	for _, v := range rawData {
		fmt.Println(v)
		str := strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.Latitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(v.SOG, 'f', -1, 64) + "," + strconv.FormatFloat(v.COG, 'f', -1, 64)
		n, err := doorLineWriter.WriteString(str + "\r\n")
		if n < 0 && err != nil {
			log.Println(err)
		}
		doorLineWriter.Flush()
	}
}

func CulDeviation(createFile string, cleanFile string, repairFile string) {
	file, err := os.Open("data/" + createFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	bfRd := bufio.NewReader(file)
	createData := make([]constant.PositionMeta, 0)
	cleanData := make([]constant.PositionMeta, 0)
	repairData := make([]constant.PositionMeta, 0)

	tmp := 0
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			tmp += 1
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		sog, _ := strconv.ParseFloat(strs[2], 64)
		cog, _ := strconv.ParseFloat(strs[3], 64)
		createData = append(createData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
			SOG:       sog,
			COG:       cog,
		})
	}
	fmt.Println(tmp)

	file, err = os.Open("data/" + cleanFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	bfRd = bufio.NewReader(file)
	sum := 0
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			if strs[0] == "XXXXXXXXX" {
				sum += 1
				cleanData = append(cleanData, constant.PositionMeta{
					MMSI: -1,
				})
			}
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		sog, _ := strconv.ParseFloat(strs[2], 64)
		cog, _ := strconv.ParseFloat(strs[3], 64)
		cleanData = append(cleanData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
			SOG:       sog,
			COG:       cog,
		})
	}
	fmt.Println(sum)

	cnt := 0
	file, err = os.Open("data/" + repairFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	bfRd = bufio.NewReader(file)
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			cnt += 1
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
		repairData = append(repairData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
			SOG:       sog,
			COG:       cog,
		})
	}
	fmt.Println(cnt)

	need := 0
	lonDeviation := 0.0
	lanDeviation := 0.0
	cogDeviation := 0.0
	sogDeviation := 0.0
	for index, v := range cleanData {
		if v.MMSI == -1 {
			need += 1
		}
		preData := createData[index]
		nowData := repairData[index]
		lonDeviation += math.Abs(preData.Longitude-nowData.Longitude) * 100 / preData.Longitude
		lanDeviation += math.Abs(preData.Latitude-nowData.Latitude) * 100 / preData.Latitude
		cogDeviation += math.Abs(preData.COG-nowData.COG) * 100 / preData.COG
		sogDeviation += math.Abs(preData.SOG-nowData.SOG) * 100 / preData.SOG
	}
	fmt.Println(len(createData), len(cleanData), len(repairData))
	fmt.Println("need        :", need)
	fmt.Println("lonDeviation:", lonDeviation/float64(need))
	fmt.Println("lanDeviation:", lanDeviation/float64(need))
	fmt.Println("cogDeviation:", cogDeviation/float64(need))
	fmt.Println("sogDeviation:", sogDeviation/float64(need))
}
