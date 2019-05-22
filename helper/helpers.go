package helper

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func SliceDividePrintln(slice []int, divisor float64) {
	for _, v := range slice {
		fmt.Printf("%8.2f", float64(v)/divisor)
	}
	fmt.Print("\n")
}

func CulTrafficResponsePrint(response *constant.CulTrafficResponse, lotDivide int, latDivide int, day float64) {
	fmt.Println("=========Traffic=========")
	fmt.Println("Traffic:             ", response.TrafficData.Traffic)
	fmt.Println("BigShipTraffic:      ", response.TrafficData.BigShipTraffic)
	fmt.Println("SmallShipTraffic:    ", response.TrafficData.SmallShipTraffic)
	fmt.Println("Type0ShipTraffic:    ", response.TrafficData.Type0ShipTraffic)
	fmt.Println("Type6xShipTraffic:   ", response.TrafficData.Type6xShipTraffic)
	fmt.Println("Type7xShipTraffic:   ", response.TrafficData.Type7xShipTraffic)
	fmt.Println("Type8xShipTraffic:   ", response.TrafficData.Type8xShipTraffic)
	fmt.Println("OtherTypeShipTraffic:", response.TrafficData.OtherTypeShipTraffic)
	fmt.Println("=========DayTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].Traffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayBigShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].BigShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DaySmallShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].SmallShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType0ShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].Type0ShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType6xShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].Type6xShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType7xShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].Type7xShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayType8xShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].Type8xShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========DayOtherTypeShipTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8.2f", float64(response.AreaTraffics[i][j].OtherTypeShipTraffic)/day)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========HourTraffic=========")
	fmt.Println("Hour, Traffic, BigShipTraffic, SmallShipTraffic Type0ShipTraffic Type6xShipTraffic Type7xShipTraffic Type8xShipTraffic OtherTypeShipTraffic")
	for i := 0; i < 24; i += 1 {
		fmt.Printf("%2.d: %8.2f%8.2f%8.2f%8.2f%8.2f%8.2f%8.2f%8.2f\n",
			i,
			float64(response.TrafficData.HourTrafficSum[i])/day,
			float64(response.TrafficData.HourBigShipTrafficSum[i])/day,
			float64(response.TrafficData.HourSmallShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType0ShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType6xShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType7xShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType8xShipTrafficSum[i])/day,
			float64(response.TrafficData.HourOtherTypeShipTraffic[i]),
		)
	}

	fmt.Println("Traffic")
	for _, v := range response.TrafficData.HourTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("BigShipTraffic")
	for _, v := range response.TrafficData.HourBigShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("SmallShipTraffic")
	for _, v := range response.TrafficData.HourSmallShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("Type0ShipTraffic")
	for _, v := range response.TrafficData.HourType0ShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("Type6xShipTraffic")
	for _, v := range response.TrafficData.HourType6xShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("Type7xShipTraffic")
	for _, v := range response.TrafficData.HourType7xShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("Type8xShipTraffic")
	for _, v := range response.TrafficData.HourType8xShipTrafficSum {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("OtherTypeShipTraffic")
	for _, v := range response.TrafficData.HourOtherTypeShipTraffic {
		fmt.Print(float64(v)/day, ",")
	}
	fmt.Println("")

	fmt.Println("=========AreaTraffic=========")
	fmt.Print("[")
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			traffic := response.AreaTraffics[i][j].Traffic
			value := 5
			if traffic > 100*int(day) {
				value = 0
			} else if traffic > 75*int(day) {
				value = 1
			} else if traffic > 50*int(day) {
				value = 2
			} else if traffic > 25*int(day) {
				value = 3
			} else if traffic > 0 {
				value = 4
			}
			fmt.Print("[", i, ",", j, ",", value, "],")
		}
	}
	fmt.Print("]\r\n")
	//for j := latDivide - 1; j >= 0; j -= 1 {
	//	for i := 0; i < lotDivide; i += 1 {
	//		fmt.Println(i, ",", j, ":")
	//		fmt.Printf("Traffic:                   %8.2f\n", float64(response.AreaTraffics[i][j].Traffic)/day)
	//		fmt.Printf("Repetition Rate            %8.2f\n", 100*float64(SliceSum(response.AreaTraffics[i][j].HourTraffic))/float64(response.AreaTraffics[i][j].Traffic))
	//		fmt.Print("Hour Traffic:                 ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourTraffic, day)
	//		fmt.Print("Hour Big Ship Traffic:        ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourBigShipTraffic, day)
	//		fmt.Print("Hour Small Ship Traffic:      ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourSmallShipTraffic, day)
	//		fmt.Print("Hour Type 0 Ship Traffic:     ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourType0ShipTraffic, day)
	//		fmt.Print("Hour Type 6x Ship Traffic:    ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourType6xShipTraffic, day)
	//		fmt.Print("Hour Type 7x Ship Traffic:    ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourType7xShipTraffic, day)
	//		fmt.Print("Hour Type 8x Ship Traffic:    ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourType8xShipTraffic, day)
	//		fmt.Print("Hour Other Type Ship Traffic: ")
	//		SliceDividePrintln(response.AreaTraffics[i][j].HourOtherTypeShipTraffic, day)
	//		fmt.Println("==================")
	//	}
	//}
}

func CulDensityResponsePrint(response *constant.CulDensityResponse, lotDivide int, latDivide int) {
	fmt.Println("ShipDensity:      ", response.DensityData.ShipDensity)
	fmt.Println("SmallShipDensity: ", response.DensityData.SmallShipDensity)
	fmt.Println("BigShipDensity:   ", response.DensityData.BigShipDensity)
	fmt.Println("Type0ShipDensity: ", response.DensityData.Type0Density)
	fmt.Println("Type6xShipDensity:", response.DensityData.Type6xDensity)
	fmt.Println("Type7xShipDensity:", response.DensityData.Type7xDensity)
	fmt.Println("Type8xShipDensity:", response.DensityData.Type8xDensity)
	fmt.Println("=========AreaDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].Density)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaSmallShipDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].SmallShipDensity)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaBigShipDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].BigShipDensity)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaType0Density=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].Type0Density)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaType6xDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].Type6xDensity)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaType7xDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].Type7xDensity)
		}
		fmt.Print("\n")
	}
	fmt.Println("=========AreaType8xDensity=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%8d", response.AreaDensity[i][j].Type8xDensity)
		}
		fmt.Print("\n")
	}
}

func CulSpeedResponsePrint(response *constant.CulSpeedResponse, lotDivide int, latDivide int) {
	fmt.Println("ShipCnt       :", response.SpeedData.ShipCnt)
	fmt.Println("ShipSpeed     :", response.SpeedData.ShipSpeed)
	fmt.Println("ShipSpeedRange:", response.SpeedData.ShipSpeedRange)
	fmt.Println("=========AreaShipSpeed=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Printf("%2d %6.2f|", response.AreaSpeed[i][j].ShipCnt, response.AreaSpeed[i][j].ShipSpeed)
		}
		fmt.Print("\n")
	}
}

func CulDoorLineResponsePrint(response *constant.CulDoorLineResponse) {
	fmt.Println("Cnt           ", response.Cnt)
	fmt.Println("DeWeightingCnt", response.DeWeightingCnt)
}

func CulSpacingResponsePrint(response *constant.CulSpacingResponse) {
	fmt.Println("MinSpacing  :", response.MinSpacing, " ", response.MinSpaceA, "(", response.APosition.Longitude, response.APosition.Latitude, ")", "---", response.MinSpaceB, "(", response.BPosition.Longitude, response.BPosition.Latitude, ")")
	fmt.Println("SpacingRange:", response.SpacingRange)
	for k, v := range response.SpacingMap {
		fmt.Println(k, ":", v)
	}
}

func CulMeetingResponsePrint(response *constant.CulMeetingResponse) {
	fmt.Println("SimpleMeeting        :", response.SimpleMeeting)
	fmt.Println("ComplexMeeting       :", response.ComplexMeeting)
	fmt.Println("SimpleDamageMeeting  :", response.SimpleDamageMeeting)
	fmt.Println("ComplexDamageMeeting :", response.ComplexDamageMeeting)
	fmt.Println("ForecastDamageMeeting:", response.ForecastDamageMeeting)
	fmt.Println("DamageMeetingAvoid   :", response.DamageMeetingAvoid)
	fmt.Println("EvasionRate          :", 100*float64(response.DamageMeetingAvoid)/float64(response.ForecastDamageMeeting))
	for i := 0; i < 12; i++ {
		start := i * 30
		end := (i + 1) * 30
		fmt.Printf("%3d° - %3d°: ForecastDamageMeeting: %8d, DamageMeetingAvoid: %8d, EvasionRate: %6.2f\n",
			start, end, response.AngleForecastDamageMeeting[i], response.AngleDamageMeetingAvoid[i],
			100*float64(response.AngleDamageMeetingAvoid[i])/float64(response.AngleForecastDamageMeeting[i]))
	}
	fmt.Println("Sum:", SliceSum(response.AngleForecastDamageMeeting), SliceSum(response.AngleDamageMeetingAvoid))
}

func DataFmt(data *constant.Data) string {
	Time := time.Date(data.Year, time.Month(data.Month), data.Day, data.Hour, data.Minute, data.Second, 0, time.UTC)
	return Time.Format("2006-01-02 03:04:05 PM")
}

func AlertPrint(alert *constant.Alert) {
	fmt.Printf("%8d: Distance: %6.2fm, Relative Azimuth: %3.2f°, COG: %3.2f°, SOG: %3.2fnm/h, UDCPA: %.4f, UTCPA: %.4f, UB: %.4f, UD: %.4f, UV: %.4f, Danger: %.4f",
		alert.MMSI, alert.Distance, alert.Azimuth, alert.ShipTrack.COG, alert.ShipTrack.SOG, alert.UDCPA, alert.UTCPA, alert.UB, alert.UD, alert.UV, alert.Danger)
}

func EarlyWarningResponsePrint(response *constant.EarlyWarningResponse) {
}

func PointListOutput(fileName string, positions []constant.PositionMeta) {
	file, err := os.Create("data/" + fileName + ".txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	file.Sync()
	writer := bufio.NewWriter(file)

	for index, v := range positions {
		str := strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.Latitude, 'f', -1, 64)
		if index != len(positions)-1 {
			str += "-"
		}
		n, err := writer.WriteString(str)
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()
	}
}

func FmtPrintList(positions []constant.PositionMeta) {
	for _, v := range positions {
		fmt.Println(v)
	}
}

func GetPositionFromFile(fileName string) ([]constant.PositionMeta, error) {
	file, err := os.Open("data/" + fileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		file.Close()
	}()
	rawData := make([]constant.PositionMeta, 0)
	bfRd := bufio.NewReader(file)
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		sog, _ := strconv.ParseFloat(strs[2], 64)
		cog, _ := strconv.ParseFloat(strs[3], 64)
		year, _ := strconv.Atoi(strs[4])
		month, _ := strconv.Atoi(strs[5])
		day, _ := strconv.Atoi(strs[6])
		hour, _ := strconv.Atoi(strs[7])
		minute, _ := strconv.Atoi(strs[8])
		second, _ := strconv.Atoi(strs[9])
		rawData = append(rawData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
			SOG:       sog,
			COG:       cog,
			Year:      year,
			Month:     month,
			Day:       day,
			Hour:      hour,
			Minute:    minute,
			Second:    second,
			MMSI:      99999999,
		})
	}
	return rawData, nil
}
