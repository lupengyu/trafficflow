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
