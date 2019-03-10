package helper

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"time"
)

func SliceDividePrintln(slice []int, divisor float64) {
	for _, v := range slice {
		fmt.Printf("%8.2f", float64(v)/divisor)
	}
	fmt.Print("\n")
}

func CulTrafficResponsePrint(response *constant.CulTrafficResponse, lotDivide int, latDivide int, day float64) {
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
	fmt.Println("=========HourTraffic=========")
	fmt.Println("Hour, Traffic, BigShipTraffic, SmallShipTraffic Type0ShipTraffic Type6xShipTraffic Type7xShipTraffic Type8xShipTraffic")
	for i := 0; i < 24; i += 1 {
		fmt.Printf("%2.d: %8.2f%8.2f%8.2f%8.2f%8.2f%8.2f%8.2f\n",
			i,
			float64(response.TrafficData.HourTrafficSum[i])/day,
			float64(response.TrafficData.HourBigShipTrafficSum[i])/day,
			float64(response.TrafficData.HourSmallShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType0ShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType6xShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType7xShipTrafficSum[i])/day,
			float64(response.TrafficData.HourType8xShipTrafficSum[i])/day,
		)
	}
	fmt.Println("=========AreaTraffic=========")
	for j := latDivide - 1; j >= 0; j -= 1 {
		for i := 0; i < lotDivide; i += 1 {
			fmt.Println(i, ",", j, ":")
			fmt.Printf("Traffic:                   %8.2f\n", float64(response.AreaTraffics[i][j].Traffic)/day)
			fmt.Printf("Repetition Rate            %8.2f\n", 100*float64(SliceSum(response.AreaTraffics[i][j].HourTraffic))/float64(response.AreaTraffics[i][j].Traffic))
			fmt.Print("Hour Traffic:              ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourTraffic, day)
			fmt.Print("Hour Big Ship Traffic:     ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourBigShipTraffic, day)
			fmt.Print("Hour Small Ship Traffic:   ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourSmallShipTraffic, day)
			fmt.Print("Hour Type 0 Ship Traffic:  ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourType0ShipTraffic, day)
			fmt.Print("Hour Type 6x Ship Traffic: ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourType6xShipTraffic, day)
			fmt.Print("Hour Type 7x Ship Traffic: ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourType7xShipTraffic, day)
			fmt.Print("Hour Type 8x Ship Traffic: ")
			SliceDividePrintln(response.AreaTraffics[i][j].HourType8xShipTraffic, day)
			fmt.Println("==================")
		}
	}
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
}

func DataFmt(data *constant.Data) string {
	Time := time.Date(data.Year, time.Month(data.Month), data.Day, data.Hour, data.Minute, data.Second, 0, time.UTC)
	return Time.Format("2006-01-02 03:04:05 PM")
}

func AlertPrint(alert *constant.Alert) {
	fmt.Printf("%8d: Distance: %6.2fm, Azimuth: %3.2f°, COG: %3.2f°, SOG: %3.2fnm/h",
		alert.MMSI, alert.Distance, alert.Azimuth, alert.ShipTrack.COG, alert.ShipTrack.SOG)
}

func EarlyWarningResponsePrint(response *constant.EarlyWarningResponse) {
	for _, v := range response.Warning {
		fmt.Println("Time:", DataFmt(v.Time))
		for _, v2 := range v.Alerts {
			AlertPrint(v2)
			if v2.IsEmergency {
				fmt.Print(" Emergency !!!: ship domain invaded")
			} else if v2.MeetingIntersection != nil {
				fmt.Printf(" Emergency !!!: ship domain will be invaded, TCPA: %4.2fs, DCPA: %4.2fm, Azimuth: %3.2f°",
					v2.MeetingIntersection.TCPA, v2.MeetingIntersection.DCPA, v2.MeetingIntersection.Azimuth)
			}
			fmt.Printf("\n")
		}
	}
}
