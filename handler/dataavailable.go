package handler

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"os"
	"strconv"
	"strings"
)

func DataAvailable(fileName string) {
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
		if helper.PositionAvailable(data, preData) == false {
			tmp += 10
			continue
		} else {
			tmp = 0
		}
		preData = data
	}
	fmt.Println("All Data Pass")
}
//118.07376715856918,24.482426457322347,10.156921212121214,199.4075012994731