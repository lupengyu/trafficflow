package sql

import (
	"database/sql"
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"log"
)

func GetPositionWithShipID(shipID string) ([]constant.PositionMeta, error) {
	rows, err := mysql.DB.Query("SELECT * from position where MMSI = ?", shipID)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, errors.New("查询出错了: SELECT * from position where MMSI = " + shipID)
	}
	positions := make([]constant.PositionMeta, 0)
	for rows.Next() {
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second)
		if err != nil {
			return nil, errors.New("rows fail")
		}
		positions = append(positions, pos)
	}
	return positions, nil
}

func GetPositionWithDuration(beginTime *constant.Data, endTime *constant.Data) (*sql.Rows, error) {
	rows, err := mysql.DB.Query(
		"select * from position where (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?))))) order by id",
		beginTime.Year, beginTime.Year,
		beginTime.Month, beginTime.Month,
		beginTime.Day, beginTime.Day,
		beginTime.Hour, beginTime.Hour,
		beginTime.Minute, beginTime.Minute,
		beginTime.Second,
		endTime.Year, endTime.Year,
		endTime.Month, endTime.Month,
		endTime.Day, endTime.Day,
		endTime.Hour, endTime.Hour,
		endTime.Minute, endTime.Minute,
		endTime.Second,
	)
	return rows, err
}
