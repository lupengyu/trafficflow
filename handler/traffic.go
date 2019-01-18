package handler

import (
	"fmt"
	"github.com/lupengyu/trafficflow/client/sql"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"strconv"
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
	// 查询时间段内的数据
	rows, err := mysql.DB.Query(
		"select * from position where (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?))))) order by id",
		request.StartTime.Year, request.StartTime.Year,
		request.StartTime.Month, request.StartTime.Month,
		request.StartTime.Day, request.StartTime.Day,
		request.StartTime.Hour, request.StartTime.Hour,
		request.StartTime.Minute, request.StartTime.Minute,
		request.StartTime.Second,
		request.EndTime.Year, request.EndTime.Year,
		request.EndTime.Month, request.EndTime.Month,
		request.EndTime.Day, request.EndTime.Day,
		request.EndTime.Hour, request.EndTime.Hour,
		request.EndTime.Minute, request.EndTime.Minute,
		request.EndTime.Second,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

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

	infoMetas := make(map[int]*constant.InfoMeta)
	//infoMetas, err := sql.InitShipInfo()
	//if err != nil {
	//	log.Println(err)
	//	return nil, err
	//}
	//log.Println("InitShipInfo success")

	// 数据循环遍历计算
	index := 0
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
		nowTime := &constant.Data{
			Year:  pos.Year,
			Month: pos.Month,
			Day:   pos.Day,
			Hour:  pos.Hour,
		}

		// 获取船舶船体数据
		shipInfo := infoMetas[pos.MMSI]
		if shipInfo == nil {
			item, err := sql.GetInfoFirstWithShipID(strconv.Itoa(pos.MMSI))
			if err != nil {
				shipInfo = &constant.InfoMeta{}
			} else {
				infoMetas[pos.MMSI] = &item
				shipInfo = &item
			}
		}

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
					} else if shipInfo.ShipType / 10 == 6 {
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
					} else if shipInfo.ShipType / 10 == 7 {
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
					} else if shipInfo.ShipType / 10 == 8 {
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
					} else if shipInfo.ShipType / 10 == 6 {
						// type 6x
						if areaTraffic[i][j].Type6xShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type6xShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType6xShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					} else if shipInfo.ShipType / 10 == 7 {
						// type 7x
						if areaTraffic[i][j].Type7xShipMap[pos.MMSI] == 1 {
							areaTraffic[i][j].Type7xShipMap[pos.MMSI] = 0
						}
						if areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] == 1 {
							areaTraffic[i][j].HourType7xShipMap[nowTime.Hour][pos.MMSI] = 0
						}
					} else if shipInfo.ShipType / 10 == 8 {
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
	log.Println("Rows:", index)
	fmt.Println("=========DayTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Traffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DaySmallShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].SmallShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayBigShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].BigShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType0ShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Type0ShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType6xShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Type6xShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType7xShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Type7xShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType8xShipTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Print(areaTraffic[i][j].Type8xShipTraffic, "\t")
		}
		fmt.Print("\n")
	}
	fmt.Println("=========HourTraffic=========")
	fmt.Println("Hour, Traffic, BigShipTraffic, SmallShipTraffic Type0ShipTraffic Type6xShipTraffic Type7xShipTraffic Type8xShipTraffic")
	for i := 0; i < 24; i += 1 {
		fmt.Println(i, ":", trafficData.HourTrafficSum[i], trafficData.HourBigShipTrafficSum[i], trafficData.HourSmallShipTrafficSum[i], trafficData.HourType0ShipTrafficSum[i], trafficData.HourType6xShipTrafficSum[i], trafficData.HourType7xShipTrafficSum[i], trafficData.HourType8xShipTrafficSum[i])
	}
	fmt.Println("=========AreaTraffic=========")
	for j := request.LatDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < request.LotDivide; i += 1 {
			fmt.Println(i, ",", j, ":")
			fmt.Println("Traffic:                  ", areaTraffic[i][j].Traffic)
			fmt.Println("Hour Traffic:             ", areaTraffic[i][j].HourTraffic)
			fmt.Println("Hour Big Ship Traffic:    ", areaTraffic[i][j].HourBigShipTraffic)
			fmt.Println("Hour Small Ship Traffic:  ", areaTraffic[i][j].HourSmallShipTraffic)
			fmt.Println("Hour Type 0 Ship Traffic: ", areaTraffic[i][j].HourType0ShipTraffic)
			fmt.Println("Hour Type 6x Ship Traffic:", areaTraffic[i][j].HourType6xShipTraffic)
			fmt.Println("Hour Type 7x Ship Traffic:", areaTraffic[i][j].HourType7xShipTraffic)
			fmt.Println("Hour Type 8x Ship Traffic:", areaTraffic[i][j].HourType8xShipTraffic)
			fmt.Println("==================")
		}
	}
	return nil, nil
}
