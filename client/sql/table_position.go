package sql

import (
	"database/sql"
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"log"
)

func GetPositionWithShipID(shipID int) ([]constant.PositionMeta, error) {
	rows, err := mysql.DB.Query("SELECT * from position where MMSI = ? order by id", shipID)
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return nil, errors.New("查询出错了: SELECT * from position where MMSI = " + string(shipID))
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

func GetNewPositionWithDuration(beginTime *constant.Data, endTime *constant.Data) (*sql.Rows, error) {
	rows, err := mysql.DB.Query(
		"select * from newposition where (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?))))) order by id",
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

func GetNewPositionWithShipIDWithDuration(shipID int, beginTime *constant.Data, endTime *constant.Data) ([]constant.PositionMeta, error) {
	rows, err := mysql.DB.Query(
		"select * from positiontest where MMSI = ? and (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?)))))",
		shipID,
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
	if err != nil {
		return nil, errors.New("查询出错了: SELECT * from position where MMSI = " + string(shipID))
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

func GetAllPosition() (*sql.Rows, error) {
	rows, err := mysql.DB.Query("select * from position order by id")
	return rows, err
}

func AddNewShipPosition(position constant.PositionMeta) {
	_, err := mysql.NewPositionFMT.Exec(
		position.MessageType,
		position.RepeatIndicator,
		position.MMSI,
		position.NavigationStatus,
		position.ROT,
		position.SOG,
		position.PositionAccuracy,
		position.Longitude,
		position.Latitude,
		position.COG,
		position.HDG,
		position.TimeStamp,
		position.ReservedForRegional,
		position.RAIMFlag,
		position.Year,
		position.Month,
		position.Day,
		position.Hour,
		position.Minute,
		position.Second)
	if err != nil {
		log.Println(err)
	}
}
