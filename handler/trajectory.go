package handler

import (
	"bufio"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"log"
	"os"
	"strconv"
)

func GetTrajectory(request *constant.GetTrajectoryRequest) {
	positions, err := sql.GetPositionWithShipID(request.MMSI)
	if err != nil {
		log.Println("查询失败")
		return
	}
	log.Println("查询完成，总长度:", len(positions))
	trajectory, err := os.Create("data/trajectory.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		trajectory.Close()
	}()
	trajectory.Sync()
	trajectoryWriter := bufio.NewWriter(trajectory)

	for index, v := range positions {
		str := strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.Latitude, 'f', -1, 64)
		if index != len(positions)-1 {
			str += "-"
		}
		n, err := trajectoryWriter.WriteString(str)
		if n != len(str) && err != nil {
			log.Println(err)
		}
		trajectoryWriter.Flush()
	}
}
