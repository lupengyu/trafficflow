package cache

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"log"
	"strconv"
)

var infoCache map[int]*constant.InfoMeta

func InitShipInfo() error {
	infoCache = make(map[int]*constant.InfoMeta)
	return nil
}

func GetShipInfo(shipID int) *constant.InfoMeta {
	shipInfo := infoCache[shipID]
	if shipInfo == nil {
		item, err := sql.GetInfoFirstWithShipID(strconv.Itoa(shipID))
		if err != nil {
			log.Println(err)
			shipInfo = &constant.InfoMeta{ShipType: -1}
		} else {
			infoCache[shipID] = &item
			shipInfo = &item
		}
	}
	return shipInfo
}
