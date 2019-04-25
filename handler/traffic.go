package handler

import (
	"fmt"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

/*
	计算交通量
		计算维度：
			区域流量（10x10）
			小时流量
			船舶类型流量：
				0：	渔船
				6x:	客轮
				7x:	货轮
				8x:	油轮
			船舶大小：
				小型船:	length < 100m
				大型船:	length >= 100m
*/
func CulTraffic(request *constant.CulTrafficRequest) (response *constant.CulTrafficResponse, err error) {
	fmt.Println("start traffic")
	// 查询时间段内的数据
	rows, err := sql.GetPositionWithDuration(request.StartTime, request.EndTime)
	if err != nil {
		fmt.Println("sql.GetPositionWithDuration error")
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// 初始化初始数据
	preTime := &constant.Data{
		Year:  0,
		Month: 0,
		Day:   0,
	}
	trafficData := &constant.TrafficData{
		HourTrafficSum:           make([]int, 24),
		HourBigShipTrafficSum:    make([]int, 24),
		HourSmallShipTrafficSum:  make([]int, 24),
		HourType0ShipTrafficSum:  make([]int, 24),
		HourType6xShipTrafficSum: make([]int, 24),
		HourType7xShipTrafficSum: make([]int, 24),
		HourType8xShipTrafficSum: make([]int, 24),
	}
	var areaTraffic [][]constant.AreaTraffic
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaTraffic, request.LatDivide)
		areaTraffic = append(areaTraffic, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			// ship
			areaTraffic[i][j].ShipMap = make(map[int]int)
			areaTraffic[i][j].HourShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourTraffic = make([]int, 24)
			// small ship
			areaTraffic[i][j].SmallShipMap = make(map[int]int)
			areaTraffic[i][j].HourSmallShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourSmallShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourSmallShipTraffic = make([]int, 24)
			// big ship
			areaTraffic[i][j].BigShipMap = make(map[int]int)
			areaTraffic[i][j].HourBigShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourBigShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourBigShipTraffic = make([]int, 24)
			// type 0 ship
			areaTraffic[i][j].Type0ShipMap = make(map[int]int)
			areaTraffic[i][j].HourType0ShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType0ShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType0ShipTraffic = make([]int, 24)
			// type 6x ship
			areaTraffic[i][j].Type6xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType6xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType6xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType6xShipTraffic = make([]int, 24)
			// type 7x ship
			areaTraffic[i][j].Type7xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType7xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType7xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType7xShipTraffic = make([]int, 24)
			// type 8x ship
			areaTraffic[i][j].Type8xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType8xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType8xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType8xShipTraffic = make([]int, 24)
		}
	}

	// 数据循环遍历计算
	index := 0
	miss := 0
	for rows.Next() {
		// 计数显示
		if index%10000 == 0 {
			log.Println(index)
		}
		index += 1
		// 数据绑定
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second,
		)
		if err != nil {
			return nil, err
		}

		/*
			计算当前船只经纬度所处地图分块
			对于不在地图经纬度区域内的数据剔除
		*/
		if pos.SOG == 0 {
			continue
		}
		longitudeArea := helper.LongitudeArea(pos.Longitude, request.LotDivide)
		latitudeArea := helper.LatitudeArea(pos.Latitude, request.LatDivide)
		if longitudeArea == -1 || latitudeArea == -1 {
			continue
		}
		if pos.Hour > 23 {
			continue
		}
		nowTime := &constant.Data{
			Year:  pos.Year,
			Month: pos.Month,
			Day:   pos.Day,
			Hour:  pos.Hour,
		}

		// 获取船舶船体数据
		shipInfo := cache.GetShipInfo(pos.MMSI)

		/*
			判断日期有没有刷新
				有：刷新记录
		*/
		if !helper.DayEqual(preTime, nowTime) {
			preTime = nowTime
			for i := 0; i < request.LotDivide; i += 1 {
				for j := 0; j < request.LatDivide; j += 1 {
					areaTraffic[i][j].ShipMap = make(map[int]int)
					areaTraffic[i][j].SmallShipMap = make(map[int]int)
					areaTraffic[i][j].BigShipMap = make(map[int]int)
					areaTraffic[i][j].Type0ShipMap = make(map[int]int)
					areaTraffic[i][j].Type6xShipMap = make(map[int]int)
					areaTraffic[i][j].Type7xShipMap = make(map[int]int)
					areaTraffic[i][j].Type8xShipMap = make(map[int]int)
					for k := 0; k < 24; k += 1 {
						areaTraffic[i][j].HourShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourSmallShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourBigShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourType0ShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourType6xShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourType7xShipMap[k] = make(map[int]int)
						areaTraffic[i][j].HourType8xShipMap[k] = make(map[int]int)
					}
				}
			}
		}

		// 数据记录
		for i := 0; i < request.LotDivide; i += 1 {
			for j := 0; j < request.LatDivide; j += 1 {
				if i == longitudeArea && j == latitudeArea {
					// 添加本区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 0 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 1
						areaTraffic[i][j].Traffic += 1
					}
					if areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] == 0 {
						areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] = 1
						areaTraffic[i][j].HourTraffic[nowTime.Hour] += 1
						trafficData.HourTrafficSum[nowTime.Hour] += 1
					}
					// 船舶大小记录
					if shipInfo.Length != 0 && shipInfo.Length != 1022 {
						if shipInfo.Length >= constant.BigShipLength {
							if areaTraffic[i][j].BigShipMap[pos.MMSI] == 0 {
								areaTraffic[i][j].BigShipMap[pos.MMSI] = 1
								areaTraffic[i][j].BigShipTraffic += 1
							}
							if areaTraffic[i][j].HourBigShipMap[nowTime.Hour][pos.MMSI] == 0 {
								areaTraffic[i][j].HourBigShipMap[nowTime.Hour][pos.MMSI] = 1
								areaTraffic[i][j].HourBigShipTraffic[nowTime.Hour] += 1
								trafficData.HourBigShipTrafficSum[nowTime.Hour] += 1
							}
						} else {
							if areaTraffic[i][j].SmallShipMap[pos.MMSI] == 0 {
								areaTraffic[i][j].SmallShipMap[pos.MMSI] = 1
								areaTraffic[i][j].SmallShipTraffic += 1
							}
							if areaTraffic[i][j].HourSmallShipMap[nowTime.Hour][pos.MMSI] == 0 {
								areaTraffic[i][j].HourSmallShipMap[nowTime.Hour][pos.MMSI] = 1
								areaTraffic[i][j].HourSmallShipTraffic[nowTime.Hour] += 1
								trafficData.HourSmallShipTrafficSum[nowTime.Hour] += 1
							}
						}
					}
					// 船舶类型记录
					if shipInfo.ShipType == 0 {
						// type 0
						if areaTraffic[i][j].Type0ShipMap[pos.MMSI] == 0 {
							areaTraffic[i][j].Type0ShipMap[pos.MMSI] = 1
							areaTraffic[i][j].Type0ShipTraffic += 1
						}
						if areaTraffic[i][j].HourType0ShipMap[nowTime.Hour][pos.MMSI] == 0 {
							areaTraffic[i][j].HourType0ShipMap[nowTime.Hour][pos.MMSI] = 1
							areaTraffic[i][j].HourType0ShipTraffic[nowTime.Hour] += 1
							trafficData.HourType0ShipTrafficSum[nowTime.Hour] += 1
						}
					} else if shipInfo.ShipType/10 == 6 {
						// type 6x
						if areaTraffic[i][j].Type6xShipMap[pos.MMSI] == 0 {
							areaTraffic[i][j].Type6xShipMap[pos.MMSI] = 1
							areaTraffic[i][j].Type6xShipTraffic += 1
						}
						if areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] == 0 {
							areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] = 1
							areaTraffic[i][j].HourType6xShipTraffic[nowTime.Hour] += 1
							trafficData.HourType6xShipTrafficSum[nowTime.Hour] += 1
						}
					} else if shipInfo.ShipType/10 == 7 {
						// type 7x
						if areaTraffic[i][j].Type7xShipMap[pos.MMSI] == 0 {
							areaTraffic[i][j].Type7xShipMap[pos.MMSI] = 1
							areaTraffic[i][j].Type7xShipTraffic += 1
						}
						if areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] == 0 {
							areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] = 1
							areaTraffic[i][j].HourType7xShipTraffic[nowTime.Hour] += 1
							trafficData.HourType7xShipTrafficSum[nowTime.Hour] += 1
						}
					} else if shipInfo.ShipType/10 == 8 {
						// type 8x
						if areaTraffic[i][j].Type8xShipMap[pos.MMSI] == 0 {
							areaTraffic[i][j].Type8xShipMap[pos.MMSI] = 1
							areaTraffic[i][j].Type8xShipTraffic += 1
						}
						if areaTraffic[i][j].HourType8xShipMap[nowTime.Hour][pos.MMSI] == 0 {
							areaTraffic[i][j].HourType8xShipMap[nowTime.Hour][pos.MMSI] = 1
							areaTraffic[i][j].HourType8xShipTraffic[nowTime.Hour] += 1
							trafficData.HourType8xShipTrafficSum[nowTime.Hour] += 1
						}
					}
				} else {
					// 清除其他区域的记录
					if areaTraffic[i][j].ShipMap[pos.MMSI] == 1 {
						areaTraffic[i][j].ShipMap[pos.MMSI] = 0
					}
					if areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] == 1 {
						areaTraffic[i][j].HourShipMap[nowTime.Hour][pos.MMSI] = 0
					}
					// 船舶大小记录
					if shipInfo.Length != 0 && shipInfo.Length != 1022 {
						if shipInfo.Length >= constant.BigShipLength {
							if areaTraffic[i][j].BigShipMap[pos.MMSI] == 1 {
								areaTraffic[i][j].BigShipMap[pos.MMSI] = 0
							}
							if areaTraffic[i][j].HourBigShipMap[nowTime.Hour][pos.MMSI] == 1 {
								areaTraffic[i][j].HourBigShipMap[nowTime.Hour][pos.MMSI] = 0
							}
						} else {
							if areaTraffic[i][j].SmallShipMap[pos.MMSI] == 1 {
								areaTraffic[i][j].SmallShipMap[pos.MMSI] = 0
							}
							if areaTraffic[i][j].HourSmallShipMap[nowTime.Hour][pos.MMSI] == 1 {
								areaTraffic[i][j].HourSmallShipMap[nowTime.Hour][pos.MMSI] = 0
							}
						}
					}
					// 船舶类型记录
					if shipInfo.ShipType == 0 {
						// type 0
						if areaTraffic[i][j].Type0ShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type0ShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType0ShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType0ShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					} else if shipInfo.ShipType/10 == 6 {
						// type 6x
						if areaTraffic[i][j].Type6xShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type6xShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					} else if shipInfo.ShipType/10 == 7 {
						// type 7x
						if areaTraffic[i][j].Type7xShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type7xShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					} else if shipInfo.ShipType/10 == 8 {
						// type 8x
						if areaTraffic[i][j].Type8xShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type8xShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType8xShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType8xShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					}
				}
			}
		}
	}

	// 输出结果
	log.Println("Rows:", index, "Miss:", miss)
	return &constant.CulTrafficResponse{
		AreaTraffics: areaTraffic,
		TrafficData:  trafficData,
	}, nil
}

func CulNewTraffic(request *constant.CulTrafficRequest) (response *constant.CulTrafficResponse, err error) {
	trafficData := &constant.TrafficData{
		HourTrafficSum:           make([]int, 24),
		HourBigShipTrafficSum:    make([]int, 24),
		HourSmallShipTrafficSum:  make([]int, 24),
		HourType0ShipTrafficSum:  make([]int, 24),
		HourType6xShipTrafficSum: make([]int, 24),
		HourType7xShipTrafficSum: make([]int, 24),
		HourType8xShipTrafficSum: make([]int, 24),
		HourOtherTypeShipTraffic: make([]int, 24),
	}
	var areaTraffic [][]constant.AreaTraffic
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaTraffic, request.LatDivide)
		areaTraffic = append(areaTraffic, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			// ship
			areaTraffic[i][j].ShipMap = make(map[int]int)
			areaTraffic[i][j].HourShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourTraffic = make([]int, 24)
			// small ship
			areaTraffic[i][j].SmallShipMap = make(map[int]int)
			areaTraffic[i][j].HourSmallShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourSmallShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourSmallShipTraffic = make([]int, 24)
			// big ship
			areaTraffic[i][j].BigShipMap = make(map[int]int)
			areaTraffic[i][j].HourBigShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourBigShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourBigShipTraffic = make([]int, 24)
			// type 0 ship
			areaTraffic[i][j].Type0ShipMap = make(map[int]int)
			areaTraffic[i][j].HourType0ShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType0ShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType0ShipTraffic = make([]int, 24)
			// type 6x ship
			areaTraffic[i][j].Type6xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType6xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType6xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType6xShipTraffic = make([]int, 24)
			// type 7x ship
			areaTraffic[i][j].Type7xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType7xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType7xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType7xShipTraffic = make([]int, 24)
			// type 8x ship
			areaTraffic[i][j].Type8xShipMap = make(map[int]int)
			areaTraffic[i][j].HourType8xShipMap = make([]map[int]int, 24)
			for k := 0; k < 24; k += 1 {
				areaTraffic[i][j].HourType8xShipMap[k] = make(map[int]int)
			}
			areaTraffic[i][j].HourType8xShipTraffic = make([]int, 24)
			areaTraffic[i][j].HourOtherTypeShipTraffic = make([]int, 24)
		}
	}

	shipIDs, _ := sql.GetShip()
	length := len(shipIDs)
	for index, shipID := range shipIDs {
		// clean and repair
		// 单船交通量统计
		shipPositions, err := sql.GetNewPositionWithShipIDWithDuration(shipID.MMSI, request.StartTime, request.EndTime)
		if err != nil {
			fmt.Println("sql.GetNewPositionWithShipIDWithDuration error")
			return nil, err
		}
		// 获取船舶船体数据
		shipMap := 0
		shipDayMap := make([]int, request.Day)
		shipHourMap := make([][]int, request.Day)
		for i := 0; i < request.Day; i += 1 {
			shipDayMap[i] = 0
			shipHourMap[i] = make([]int, 24)
			for j := 0; j < 24; j++ {
				shipHourMap[i][j] = 0
			}
		}
		shipInfo := cache.GetShipInfo(shipID.MMSI)
		for _, pos := range shipPositions {
			if pos.SOG < constant.StaticShip {
				continue
			}
			longitudeArea := helper.LongitudeArea(pos.Longitude, request.LotDivide)
			latitudeArea := helper.LatitudeArea(pos.Latitude, request.LatDivide)
			if longitudeArea == -1 || latitudeArea == -1 {
				continue
			}
			if pos.Hour > 23 {
				continue
			}
			nowTime := &constant.Data{
				Year:   pos.Year,
				Month:  pos.Month,
				Day:    pos.Day,
				Hour:   pos.Hour,
				Minute: pos.Minute,
				Second: pos.Second,
			}
			diff := helper.TimeDeviation(nowTime, request.StartTime)
			day := diff / 86400
			hour := (diff % 86400) / 3600
			nowPosition := longitudeArea*request.LatDivide + latitudeArea + 1
			if shipMap != nowPosition {
				shipMap = nowPosition
				// 添加总记录
				trafficData.Traffic += 1
				if shipInfo.Length != 0 && shipInfo.Length != 1022 {
					if shipInfo.Length >= constant.BigShipLength {
						trafficData.BigShipTraffic += 1
					} else {
						trafficData.SmallShipTraffic += 1
					}
				}
				if shipInfo.ShipType == 0 {
					trafficData.Type0ShipTraffic += 1
				} else if shipInfo.ShipType/10 == 6 {
					trafficData.Type6xShipTraffic += 1
				} else if shipInfo.ShipType/10 == 7 {
					trafficData.Type7xShipTraffic += 1
				} else if shipInfo.ShipType/10 == 8 {
					trafficData.Type8xShipTraffic += 1
				} else {
					trafficData.OtherTypeShipTraffic += 1
				}
			}
			if shipDayMap[day] != nowPosition {
				shipDayMap[day] = nowPosition
				// 添加天记录
				areaTraffic[longitudeArea][latitudeArea].Traffic += 1
				if shipInfo.Length != 0 && shipInfo.Length != 1022 {
					if shipInfo.Length >= constant.BigShipLength {
						areaTraffic[longitudeArea][latitudeArea].BigShipTraffic += 1
						areaTraffic[longitudeArea][latitudeArea].HourBigShipTraffic[hour] += 1
					} else {
						areaTraffic[longitudeArea][latitudeArea].SmallShipTraffic += 1
						areaTraffic[longitudeArea][latitudeArea].HourSmallShipTraffic[hour] += 1
					}
				}
				if shipInfo.ShipType == 0 {
					areaTraffic[longitudeArea][latitudeArea].Type0ShipTraffic += 1
					areaTraffic[longitudeArea][latitudeArea].HourType0ShipTraffic[hour] += 1
				} else if shipInfo.ShipType/10 == 6 {
					areaTraffic[longitudeArea][latitudeArea].Type6xShipTraffic += 1
					areaTraffic[longitudeArea][latitudeArea].HourType6xShipTraffic[hour] += 1
				} else if shipInfo.ShipType/10 == 7 {
					areaTraffic[longitudeArea][latitudeArea].Type7xShipTraffic += 1
					areaTraffic[longitudeArea][latitudeArea].HourType7xShipTraffic[hour] += 1
				} else if shipInfo.ShipType/10 == 8 {
					areaTraffic[longitudeArea][latitudeArea].Type8xShipTraffic += 1
					areaTraffic[longitudeArea][latitudeArea].HourType8xShipTraffic[hour] += 1
				} else {
					areaTraffic[longitudeArea][latitudeArea].OtherTypeShipTraffic += 1
					areaTraffic[longitudeArea][latitudeArea].HourOtherTypeShipTraffic[hour] += 1
				}
			}
			if shipHourMap[day][hour] != nowPosition {
				shipHourMap[day][hour] = nowPosition
				// 添加小时记录
				trafficData.HourTrafficSum[hour] += 1
				areaTraffic[longitudeArea][latitudeArea].HourTraffic[hour] += 1
				if shipInfo.Length != 0 && shipInfo.Length != 1022 {
					if shipInfo.Length >= constant.BigShipLength {
						trafficData.HourBigShipTrafficSum[hour] += 1
					} else {
						trafficData.HourSmallShipTrafficSum[hour] += 1
					}
				}
				if shipInfo.ShipType == 0 {
					trafficData.HourType0ShipTrafficSum[hour] += 1
				} else if shipInfo.ShipType/10 == 6 {
					trafficData.HourType6xShipTrafficSum[hour] += 1
				} else if shipInfo.ShipType/10 == 7 {
					trafficData.HourType7xShipTrafficSum[hour] += 1
				} else if shipInfo.ShipType/10 == 8 {
					trafficData.HourType8xShipTrafficSum[hour] += 1
				} else {
					trafficData.HourOtherTypeShipTraffic[hour] += 1
				}
			}
		}
		percent := float64(100.0*index) / float64(length)
		log.Println("Progress:", percent, "%")
	}

	return &constant.CulTrafficResponse{
		AreaTraffics: areaTraffic,
		TrafficData:  trafficData,
	}, nil
}
