package sql

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
)

func Test_GetPositionWithShipID(t *testing.T) {
	defer func() {
		mysql.FreeConn()
	}()
	mysql.InitMysql()
	body, err := GetPositionWithShipID(412596777)
	assert.Empty(t, err)
	assert.Equal(t, 570, len(body))
}

func Test_DataOutput(t *testing.T) {
	defer func() {
		mysql.FreeConn()
	}()
	mysql.InitMysql()
	StartTime := &constant.Data{
		Year:   2018,
		Month:  12,
		Day:    22,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	EndTime := &constant.Data{
		Year:   2018,
		Month:  12,
		Day:    22,
		Hour:   23,
		Minute: 59,
		Second: 59,
	}
	file, err := os.Create("data/20181222.txt")
	if err != nil {
		log.Println(err)
		return
	}
	file.Sync()
	writer := bufio.NewWriter(file)
	shipIDs, _ := GetShip()
	length := len(shipIDs)
	fmt.Println(shipIDs)
	for index, shipID := range shipIDs {
		// clean and repair
		positions, err := GetNewPositionWithShipIDWithDuration(shipID.MMSI, StartTime, EndTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, pos := range positions {
			str := strconv.Itoa(pos.MMSI) + "," +
				strconv.FormatFloat(pos.SOG, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.COG, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.Longitude, 'f', -1, 64) + "," +
				strconv.FormatFloat(pos.Latitude, 'f', -1, 64) + "," +
				strconv.Itoa(pos.Year) + "," +
				strconv.Itoa(pos.Month) + "," +
				strconv.Itoa(pos.Day) + "," +
				strconv.Itoa(pos.Hour) + "," +
				strconv.Itoa(pos.Minute) + "," +
				strconv.Itoa(pos.Second) + "\r\n"
			n, err := writer.WriteString(str)
			if n != len(str) && err != nil {
				log.Println(err)
			}
			writer.Flush()
		}
		percent := float64(100.0*index) / float64(length)
		log.Println("Progress:", percent, "%")
	}
}
