package sql

import (
	"errors"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"log"
)

func GetShip() ([]constant.ShipMeta, error) {
	rows, err := mysql.DB.Query("SELECT * from ship")
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
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
