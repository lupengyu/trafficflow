package helper

import "github.com/lupengyu/trafficflow/constant"

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