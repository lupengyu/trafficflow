package cache

import "log"

func InitCache() {
	err := InitShipInfo()
	if err != nil {
		log.Println("create cache fail:", err)
		return
	}
	log.Println("create cache success")
}
