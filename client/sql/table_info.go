package sql

import (
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
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
