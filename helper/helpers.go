package helper

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"time"
)

func LongitudeArea(longitude float64, lotDivide int) int {
	divideRange := (constant.LongitudeMax - constant.LongitudeMin) / float64(lotDivide)
	startLongitude := constant.LongitudeMin
	endLongitude := startLongitude + divideRange
	for area := 0; area < lotDivide; area += 1 {
		if longitude >= startLongitude && longitude <= endLongitude {
			return area
		}
		startLongitude = endLongitude
		endLongitude = startLongitude + divideRange
	}
	return -1
}

func LatitudeArea(latitude float64, latDivide int) int {
	divideRange := (constant.LatitudeMax - constant.LatitudeMin) / float64(latDivide)
	startLatitude := constant.LatitudeMin
	endLatitude := startLatitude + divideRange
	for area := 0; area < latDivide; area += 1 {
		if latitude >= startLatitude && latitude <= endLatitude {
			return area
		}
		startLatitude = endLatitude
		endLatitude = startLatitude + divideRange
	}
	return -1
}

func SpeedRange(sog float64) int {
	if sog <= 5 {
		return 0
	} else if sog > 5 && sog <= 10 {
		return 1
	} else if sog > 10 && sog <= 15 {
		return 2
	} else if sog > 15 && sog <= 20 {
		return 3
	}
	return 4
}

func DataEqual(data1 *constant.Data, data2 *constant.Data) bool {
	if data1 == nil || data2 == nil {
		return false
	}
	if data1.Year == data2.Year &&
		data1.Month == data2.Month &&
		data1.Day == data2.Day &&
		data1.Hour == data2.Hour &&
		data1.Minute == data2.Minute &&
		data1.Second == data2.Second {
		return true
	}
	return false
}

func DayEqual(data1 *constant.Data, data2 *constant.Data) bool {
	if data1 == nil || data2 == nil {
		return false
	}
	if data1.Year == data2.Year &&
		data1.Month == data2.Month &&
		data1.Day == data2.Day {
		return true
	}
	return false
}

func DayDecrease(data *constant.Data, delta *constant.Data) *constant.Data {
	baseTime := time.Date(data.Year, time.Month(data.Month), data.Day, data.Hour, data.Minute, data.Second, 0, time.UTC)
	deltaTime := time.Duration(delta.Hour)*time.Hour - time.Duration(delta.Minute)*time.Minute - time.Duration(delta.Second)*time.Second
	resultTime := baseTime.Add(deltaTime)
	return &constant.Data{
		Year:   resultTime.Year(),
		Month:  int(resultTime.Month()),
		Day:    resultTime.Day(),
		Hour:   resultTime.Hour(),
		Minute: resultTime.Minute(),
		Second: resultTime.Second(),
	}
}

func DayIncrease(data *constant.Data, delta *constant.Data) *constant.Data {
	baseTime := time.Date(data.Year, time.Month(data.Month), data.Day, data.Hour, data.Minute, data.Second, 0, time.UTC)
	deltaTime := time.Duration(delta.Hour)*time.Hour + time.Duration(delta.Minute)*time.Minute + time.Duration(delta.Second)*time.Second
	resultTime := baseTime.Add(deltaTime)
	return &constant.Data{
		Year:   resultTime.Year(),
		Month:  int(resultTime.Month()),
		Day:    resultTime.Day(),
		Hour:   resultTime.Hour(),
		Minute: resultTime.Minute(),
		Second: resultTime.Second(),
	}
}

func SliceSum(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum
}

func IsLineInterSect(a *constant.Position, b *constant.Position, c *constant.Position, d *constant.Position) bool {
	u := (c.Longitude-a.Longitude)*(b.Latitude-a.Latitude) - (b.Longitude-a.Longitude)*(c.Latitude-a.Latitude)
	v := (d.Longitude-a.Longitude)*(b.Latitude-a.Latitude) - (b.Longitude-a.Longitude)*(d.Latitude-a.Latitude)
	w := (a.Longitude-c.Longitude)*(d.Latitude-c.Latitude) - (d.Longitude-c.Longitude)*(a.Latitude-c.Latitude)
	z := (b.Longitude-c.Longitude)*(d.Latitude-c.Latitude) - (d.Longitude-c.Longitude)*(b.Latitude-c.Latitude)
	return u*v <= 0.00000001 && w*z <= 0.00000001
}

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
