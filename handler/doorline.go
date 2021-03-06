package handler

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"os"
	"strconv"
)

/*
	计算门线通过次数
*/
func CulDoorLine(request *constant.CulDoorLineRequest) (response *constant.CulDoorLineResponse, err error) {
	// 查询时间段内的数据
	rows, err := sql.GetPositionWithDuration(request.StartTime, request.EndTime)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// 初始化数据
	cnt := 0
	deWeightingCnt := 0
	TrackMap := make(map[int]*constant.Track)

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

		//判断航迹
		nowPosition := &constant.Position{
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
		}
		if TrackMap[pos.MMSI] == nil {
			TrackMap[pos.MMSI] = &constant.Track{
				PrePosition:      nowPosition,
				DeWeightDoorLine: true,
				TrackList:        []*constant.Position{nowPosition},
			}
		} else {
			TrackMap[pos.MMSI].TrackList = append(TrackMap[pos.MMSI].TrackList, nowPosition)
			if helper.IsLineInterSect(request.StartPosition, request.EndPosition, TrackMap[pos.MMSI].PrePosition, nowPosition) == true {
				if TrackMap[pos.MMSI].DeWeightDoorLine {
					TrackMap[pos.MMSI].DeWeightDoorLine = false
					deWeightingCnt += 1
				}
				cnt += 1
			}
			TrackMap[pos.MMSI].PrePosition = nowPosition
		}
	}

	doorLine, err := os.Create("data/doorline.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)

	startLine := strconv.FormatFloat(request.StartPosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(request.StartPosition.Latitude, 'f', -1, 64) + "-" +
		strconv.FormatFloat(request.EndPosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(request.EndPosition.Latitude, 'f', -1, 64) + "\r\n"
	n, err := doorLineWriter.WriteString(startLine)
	if n != len(startLine) && err != nil {
		log.Println(err)
	}
	doorLineWriter.Flush()

	for _, v := range TrackMap {
		if v.DeWeightDoorLine == false {
			// 写入路径数据
			str := ""
			for i := 0; i < len(v.TrackList); i += 1 {
				str += strconv.FormatFloat(v.TrackList[i].Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.TrackList[i].Latitude, 'f', -1, 64)
				if i != len(v.TrackList)-1 {
					str += "-"
				}
			}
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n != len(str) && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
		}
	}

	// 输出结果
	log.Println("Rows:", index)
	return &constant.CulDoorLineResponse{
		Cnt:            cnt,
		DeWeightingCnt: deWeightingCnt,
	}, nil
}

func CulNewDoorLine(request *constant.CulDoorLineRequest) (response *constant.CulDoorLineResponse, err error) {
	doorLine, err := os.Create("data/doorline.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		doorLine.Close()
	}()
	doorLine.Sync()
	doorLineWriter := bufio.NewWriter(doorLine)
	startLine := strconv.FormatFloat(request.StartPosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(request.StartPosition.Latitude, 'f', -1, 64) + "-" +
		strconv.FormatFloat(request.EndPosition.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(request.EndPosition.Latitude, 'f', -1, 64) + "\r\n"
	n, err := doorLineWriter.WriteString(startLine)
	if n != len(startLine) && err != nil {
		log.Println(err)
	}
	doorLineWriter.Flush()

	shipIDs, _ := sql.GetShip()
	length := len(shipIDs)
	cnt := 0
	deWeightingCnt := 0
	for index, shipID := range shipIDs {
		shipPositions, err := sql.GetNewPositionWithShipIDWithDuration(shipID.MMSI, request.StartTime, request.EndTime)
		if err != nil {
			fmt.Println("sql.GetNewPositionWithShipIDWithDuration error")
			return nil, err
		}
		if len(shipPositions) == 0 {
			continue
		}
		prePosition := &constant.Position{
			Longitude: shipPositions[0].Longitude,
			Latitude:  shipPositions[0].Latitude,
		}
		deWeight := true
		for _, pos := range shipPositions {
			nowPosition := &constant.Position{
				Longitude: pos.Longitude,
				Latitude:  pos.Latitude,
			}
			if helper.IsLineInterSect(request.StartPosition, request.EndPosition, prePosition, nowPosition) == true {
				cnt += 1
				if deWeight == true {
					deWeight = false
					deWeightingCnt += 1
				}
			}
			prePosition = nowPosition
		}
		if deWeight == false {
			str := ""
			for i, pos := range shipPositions {
				str += strconv.FormatFloat(pos.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(pos.Latitude, 'f', -1, 64)
				if i != len(shipPositions)-1 {
					str += "-"
				}
			}
			n, err := doorLineWriter.WriteString(str + "\r\n")
			if n != len(str) && err != nil {
				log.Println(err)
			}
			doorLineWriter.Flush()
		}
		percent := float64(100.0*index) / float64(length)
		log.Println("Progress:", percent, "%")
	}

	// 输出结果
	return &constant.CulDoorLineResponse{
		Cnt:            cnt,
		DeWeightingCnt: deWeightingCnt,
	}, nil
}
