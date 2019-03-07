package handler

import (
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"github.com/lupengyu/trafficflow/helper"
	"log"
)

/*
	计算船舶密度
*/
func CulDensity(request *constant.CulDensityRequest) (response *constant.CulDensityResponse, err error) {
	beginTime := helper.DayDecrease(request.Time, request.DeltaT)
	endTime := helper.DayIncrease(request.Time, request.DeltaT)
	// 查询时间段内的数据
	rows, err := sql.GetPositionWithDuration(beginTime, endTime)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// 数据初始化
	shipMap := make(map[int]int)
	shipDensity := 0
	smallShipMap := make(map[int]int)
	smallShipDensity := 0
	bigShipMap := make(map[int]int)
	bigShipDensity := 0
	type0ShipMap := make(map[int]int)
	type0ShipDensity := 0
	type6xShipMap := make(map[int]int)
	type6xShipDensity := 0
	type7xShipMap := make(map[int]int)
	type7xShipDensity := 0
	type8xShipMap := make(map[int]int)
	type8xShipDensity := 0
	var areaDensity [][]constant.AreaDensity
	for i := 0; i < request.LotDivide; i += 1 {
		tmp := make([]constant.AreaDensity, request.LatDivide)
		areaDensity = append(areaDensity, tmp)
		for j := 0; j < request.LatDivide; j += 1 {
			// ship
			areaDensity[i][j].ShipMap = make(map[int]int)
			areaDensity[i][j].SmallShipShipMap = make(map[int]int)
			areaDensity[i][j].BigShipShipMap = make(map[int]int)
			areaDensity[i][j].Type0ShipMap = make(map[int]int)
			areaDensity[i][j].Type6xShipMap = make(map[int]int)
			areaDensity[i][j].Type7xShipMap = make(map[int]int)
			areaDensity[i][j].Type8xShipMap = make(map[int]int)
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
		longitudeArea := helper.LongitudeArea(pos.Longitude, request.LotDivide)
		latitudeArea := helper.LatitudeArea(pos.Latitude, request.LatDivide)
		if longitudeArea == -1 || latitudeArea == -1 {
			continue
		}
		if pos.Hour > 23 {
			continue
		}

		shipInfo := cache.GetShipInfo(pos.MMSI)

		// shipMap记录
		if shipMap[pos.MMSI] == 0 {
			shipMap[pos.MMSI] = 1
			shipDensity += 1
		}
		// 船舶大小记录
		if shipInfo.Length != 0 && shipInfo.Length != 1022 {
			if shipInfo.Length >= constant.BigShipLength {
				if bigShipMap[pos.MMSI] == 0 {
					bigShipMap[pos.MMSI] = 1
					bigShipDensity += 1
				}
			} else {
				if smallShipMap[pos.MMSI] == 0 {
					smallShipMap[pos.MMSI] = 1
					smallShipDensity += 1
				}
			}
		}
		// 船舶类型记录
		if shipInfo.ShipType == 0 {
			// type 0
			if type0ShipMap[pos.MMSI] == 0 {
				type0ShipMap[pos.MMSI] = 1
				type0ShipDensity += 1
			}
		} else if shipInfo.ShipType/10 == 6 {
			// type 6x
			if type6xShipMap[pos.MMSI] == 0 {
				type6xShipMap[pos.MMSI] = 1
				type6xShipDensity += 1
			}
		} else if shipInfo.ShipType/10 == 7 {
			// type 7x
			if type7xShipMap[pos.MMSI] == 0 {
				type7xShipMap[pos.MMSI] = 1
				type7xShipDensity += 1
			}
		} else if shipInfo.ShipType/10 == 8 {
			// type 8x
			if type8xShipMap[pos.MMSI] == 0 {
				type8xShipMap[pos.MMSI] = 1
				type8xShipDensity += 1
			}
		}

		// 数据记录
		for i := 0; i < request.LotDivide; i += 1 {
			for j := 0; j < request.LatDivide; j += 1 {
				if i == longitudeArea && j == latitudeArea {
					// 添加本区域的记录
					if areaDensity[i][j].ShipMap[pos.MMSI] == 0 {
						areaDensity[i][j].ShipMap[pos.MMSI] = 1
						areaDensity[i][j].Density += 1
					}
					// 船舶大小记录
					if shipInfo.Length != 0 && shipInfo.Length != 1022 {
						if shipInfo.Length >= constant.BigShipLength {
							if areaDensity[i][j].BigShipShipMap[pos.MMSI] == 0 {
								areaDensity[i][j].BigShipShipMap[pos.MMSI] = 1
								areaDensity[i][j].BigShipDensity += 1
							}
						} else {
							if areaDensity[i][j].SmallShipShipMap[pos.MMSI] == 0 {
								areaDensity[i][j].SmallShipShipMap[pos.MMSI] = 1
								areaDensity[i][j].SmallShipDensity += 1
							}
						}
					}
					// 船舶类型记录
					if shipInfo.ShipType == 0 {
						// type 0
						if areaDensity[i][j].Type0ShipMap[pos.MMSI] == 0 {
							areaDensity[i][j].Type0ShipMap[pos.MMSI] = 1
							areaDensity[i][j].Type0Density += 1
						}
					} else if shipInfo.ShipType/10 == 6 {
						// type 6x
						if areaDensity[i][j].Type6xShipMap[pos.MMSI] == 0 {
							areaDensity[i][j].Type6xShipMap[pos.MMSI] = 1
							areaDensity[i][j].Type6xDensity += 1
						}
					} else if shipInfo.ShipType/10 == 7 {
						// type 7x
						if areaDensity[i][j].Type7xShipMap[pos.MMSI] == 0 {
							areaDensity[i][j].Type7xShipMap[pos.MMSI] = 1
							areaDensity[i][j].Type7xDensity += 1
						}
					} else if shipInfo.ShipType/10 == 8 {
						// type 8x
						if areaDensity[i][j].Type8xShipMap[pos.MMSI] == 0 {
							areaDensity[i][j].Type8xShipMap[pos.MMSI] = 1
							areaDensity[i][j].Type8xDensity += 1
						}
					}
				} else {
					// 清除其他区域的记录
					if areaDensity[i][j].ShipMap[pos.MMSI] == 1 {
						areaDensity[i][j].ShipMap[pos.MMSI] = 0
						areaDensity[i][j].Density -= 1
					}
					// 船舶大小记录
					if shipInfo.Length != 0 && shipInfo.Length != 1022 {
						if shipInfo.Length >= constant.BigShipLength {
							if areaDensity[i][j].BigShipShipMap[pos.MMSI] == 1 {
								areaDensity[i][j].BigShipShipMap[pos.MMSI] = 0
								areaDensity[i][j].BigShipDensity -= 1
							}
						} else {
							if areaDensity[i][j].SmallShipShipMap[pos.MMSI] == 1 {
								areaDensity[i][j].SmallShipShipMap[pos.MMSI] = 0
								areaDensity[i][j].SmallShipDensity -= 1
							}
						}
					}
					// 船舶类型记录
					if shipInfo.ShipType == 0 {
						// type 0
						if areaDensity[i][j].Type0ShipMap[pos.MMSI] == 1 {
							areaDensity[i][j].Type0ShipMap[pos.MMSI] = 0
							areaDensity[i][j].Type0Density -= 1
						}
					} else if shipInfo.ShipType/10 == 6 {
						// type 6x
						if areaDensity[i][j].Type6xShipMap[pos.MMSI] == 1 {
							areaDensity[i][j].Type6xShipMap[pos.MMSI] = 0
							areaDensity[i][j].Type6xDensity -= 1
						}
					} else if shipInfo.ShipType/10 == 7 {
						// type 7x
						if areaDensity[i][j].Type7xShipMap[pos.MMSI] == 1 {
							areaDensity[i][j].Type7xShipMap[pos.MMSI] = 0
							areaDensity[i][j].Type7xDensity -= 1
						}
					} else if shipInfo.ShipType/10 == 8 {
						// type 8x
						if areaDensity[i][j].Type8xShipMap[pos.MMSI] == 1 {
							areaDensity[i][j].Type8xShipMap[pos.MMSI] = 0
							areaDensity[i][j].Type8xDensity -= 1
						}
					}
				}
			}
		}
	}

	// 输出结果
	log.Println("Rows:", index, "Miss:", miss)
	return &constant.CulDensityResponse{
		DensityData: &constant.DensityData{
			ShipDensity:      shipDensity,
			SmallShipDensity: smallShipDensity,
			BigShipDensity:   bigShipDensity,
			Type0Density:     type0ShipDensity,
			Type6xDensity:    type6xShipDensity,
			Type7xDensity:    type7xShipDensity,
			Type8xDensity:    type8xShipDensity,
		},
		AreaDensity: areaDensity,
	}, nil
}
