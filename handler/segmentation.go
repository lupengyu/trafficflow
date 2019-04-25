package handler

import (
	"bufio"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"os"
	"strconv"
)

func DataSegmentation(request *constant.DataSegmentationRequest) {
	positions, err := sql.GetPositionWithShipID(request.MMSI)
	if err != nil {
		log.Println("查询失败")
		return
	}
	length := len(positions)
	if length == 0 {
		return
	}
	log.Println("查询完成，总长度:", length)
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
			// 删除后点时间为之前点或之前点之前的数据
			if rg <= 0 {
				ignore = append(ignore, start)
			} else {
				prePosition = positions[start]
				if rg > 5*60 {
					// 间隔时间大于5分钟，区分之
					ends = append(ends, start)
					break
				}
			}
		}
		index = start
	}
	log.Println("分段完成，一共", len(ends), "段")
	log.Println(ends, ignore)
	// 输出原始数据
	start := 0
	segmentation, err := os.Create("data/segmentation.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		segmentation.Close()
	}()
	segmentation.Sync()
	segmentationWriter := bufio.NewWriter(segmentation)
	for _, v := range ends {
		for i := start; i < v; i++ {
			str := strconv.FormatFloat(positions[i].Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(positions[i].Latitude, 'f', -1, 64)
			if i != v-1 {
				str += "-"
			}
			n, err := segmentationWriter.WriteString(str)
			if n != len(str) && err != nil {
				log.Println(err)
			}
			segmentationWriter.Flush()
		}
		segmentationWriter.WriteString("\r\n")
		segmentationWriter.Flush()
		start = v
	}
	if start < length {
		for i := start; i < length; i++ {
			str := strconv.FormatFloat(positions[i].Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(positions[i].Latitude, 'f', -1, 64)
			if i != length-1 {
				str += "-"
			}
			n, err := segmentationWriter.WriteString(str)
			if n != len(str) && err != nil {
				log.Println(err)
			}
			segmentationWriter.Flush()
		}
	}
	log.Println("分段输出完成")
}
