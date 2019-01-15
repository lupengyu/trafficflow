package sql

import (
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
)

func GetShip() ([]constant.ShipMeta, error) {
	rows, err := mysql.DB.Query("SELECT * from ship")
	if err != nil {
		return nil, errors.New("查询出错了: SELECT * from info where MMSI = ")
	}
	ships := make([]constant.ShipMeta, 0)
	for rows.Next() {
		var info constant.ShipMeta
		err := rows.Scan(&info.MMSI)
		if err != nil {
			return nil, errors.New("rows fail")
		}
		ships = append(ships, info)
	}
	return ships, nil
}
