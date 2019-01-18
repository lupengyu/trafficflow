package sql

import (
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"log"
	"strconv"
)

func GetInfoWithShipID(shipID string) ([]constant.InfoMeta, error) {
	rows, err := mysql.DB.Query("SELECT * from info where MMSI = ?", shipID)
	if err != nil {
		return nil, errors.New("查询出错了: SELECT * from info where MMSI = " + shipID)
	}
	infos := make([]constant.InfoMeta, 0)
	for rows.Next() {
		var inf constant.InfoMeta
		err := rows.Scan(
			&inf.ID, &inf.NavigationStatus, &inf.MMSI, &inf.AIS, &inf.IMO, &inf.CallSign, &inf.Name,
			&inf.ShipType, &inf.A, &inf.B, &inf.C, &inf.D, &inf.Length, &inf.Width,
			&inf.PositionType, &inf.ETAMonth, &inf.ETADay, &inf.ETAHour, &inf.ETAMinute, &inf.Draft, &inf.Destination,
			&inf.Year, &inf.Month, &inf.Day, &inf.Hour, &inf.Minute, &inf.Second)
		if err != nil {
			return nil, errors.New("rows fail")
		}
		infos = append(infos, inf)
	}
	return infos, nil
}

func GetInfoFirstWithShipID(shipID string) (constant.InfoMeta, error) {
	rows, queryErr := mysql.DB.Query("SELECT * from info where MMSI = ? LIMIT 1", shipID)
	defer rows.Close()
	if queryErr != nil {
		log.Println(queryErr)
		return constant.InfoMeta{}, errors.New("查询出错了: SELECT * from info where MMSI = " + shipID + " LIMIT 1")
	}
	inf := constant.InfoMeta{}
	if rows.Next() {
		err := rows.Scan(
			&inf.ID, &inf.NavigationStatus, &inf.MMSI, &inf.AIS, &inf.IMO, &inf.CallSign, &inf.Name,
			&inf.ShipType, &inf.A, &inf.B, &inf.C, &inf.D, &inf.Length, &inf.Width,
			&inf.PositionType, &inf.ETAMonth, &inf.ETADay, &inf.ETAHour, &inf.ETAMinute, &inf.Draft, &inf.Destination,
			&inf.Year, &inf.Month, &inf.Day, &inf.Hour, &inf.Minute, &inf.Second)
		if err != nil {
			return constant.InfoMeta{}, errors.New("rows fail")
		}
	}
	return inf, nil
}

func InitShipInfo() (map[int]constant.InfoMeta, error) {
	shipMetas, err := GetShip()
	if err != nil {
		return nil, err
	}
	infoMetas := make(map[int]constant.InfoMeta)
	for _, metas := range shipMetas {
		item, err := GetInfoFirstWithShipID(strconv.Itoa(metas.MMSI))
		if err != nil {
			return nil, err
		}
		infoMetas[metas.MMSI] = item
	}
	return infoMetas, nil
}
