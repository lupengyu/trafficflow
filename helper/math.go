package helper

import (
	"github.com/lupengyu/trafficflow/constant"
	"math"
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

/*
	if a >= b
		return true
	else
		return false
*/
func DayBigger(a *constant.Data, b *constant.Data) bool {
	return TimeDeviation(a, b) >= 0
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

func TimeDeviation(a *constant.Data, b *constant.Data) int64 {
	aTime := time.Date(a.Year, time.Month(a.Month), a.Day, a.Hour, a.Minute, a.Second, 0, time.UTC)
	bTime := time.Date(b.Year, time.Month(b.Month), b.Day, b.Hour, b.Minute, b.Second, 0, time.UTC)
	delta := aTime.Unix() - bTime.Unix()
	return delta
}

func PositionSpacing(a *constant.Position, b *constant.Position) float64 {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 := a.Latitude * rad
	lng1 := a.Longitude * rad
	lat2 := b.Latitude * rad
	lng2 := b.Longitude * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius * 1000
}

// TODO: (P3)优化插值算法
func TrackInterpolation(tracks []*constant.Track) *constant.Track {
	if len(tracks) == 0 {
		return nil
	} else if len(tracks) == 1 {
		return &constant.Track{
			PrePosition: tracks[0].PrePosition,
			COG:         tracks[0].COG,
			SOG:         tracks[0].SOG,
		}
	} else {
		// track1(low) - track2(height)
		for i := 0; i < len(tracks); i++ {
			track := tracks[i]
			if track.Deviation == 0 {
				return &constant.Track{
					PrePosition: tracks[0].PrePosition,
					COG:         tracks[0].COG,
					SOG:         tracks[0].SOG,
				}
			}
		}
		track1 := tracks[0]
		track2 := tracks[0]
		startI := 1
		for ; startI < len(tracks); startI++ {
			track2 = tracks[startI]
			if track2.Deviation != track1.Deviation {
				break
			}
		}
		if math.Abs(float64(track2.Deviation)) < math.Abs(float64(track1.Deviation)) {
			tmp := track1
			track1 = track2
			track2 = tmp
		}
		for ; startI < len(tracks); startI++ {
			track := tracks[startI]
			if math.Abs(float64(track.Deviation)) < math.Abs(float64(track1.Deviation)) {
				track2 = track1
				track1 = track
			} else if math.Abs(float64(track.Deviation)) < math.Abs(float64(track2.Deviation)) && math.Abs(float64(track.Deviation)) != math.Abs(float64(track1.Deviation)) {
				track2 = track
			}
		}
		diff := track1.Deviation - track2.Deviation
		if diff == 0 {
			return &constant.Track{
				PrePosition: track1.PrePosition,
				COG:         track1.COG,
				SOG:         track1.SOG,
			}
		}
		longitudeK := (track1.PrePosition.Longitude - track2.PrePosition.Longitude) / float64(diff)
		latitudeK := (track1.PrePosition.Latitude - track2.PrePosition.Latitude) / float64(diff)
		cogK := (track1.COG - track2.COG) / float64(diff)
		sogK := (track1.SOG - track2.SOG) / float64(diff)
		return &constant.Track{
			PrePosition: &constant.Position{
				Longitude: track1.PrePosition.Longitude - float64(track1.Deviation)*longitudeK,
				Latitude:  track1.PrePosition.Latitude - float64(track1.Deviation)*latitudeK,
			},
			COG: track1.COG - float64(track1.Deviation)*cogK,
			SOG: track1.SOG - float64(track1.Deviation)*sogK,
		}
	}
}

func ArcSin(value float64) float64 {
	return 180 * math.Asin(value) / math.Pi
}

func ArcCos(value float64) float64 {
	return 180 * math.Acos(value) / math.Pi
}

func sin(angle float64) float64 {
	return math.Sin(angle * math.Pi / 180)
}

func cos(angle float64) float64 {
	return math.Cos(angle * math.Pi / 180)
}

func InEllipse(a float64, b float64, S float64, T float64, x float64, y float64, angle float64) bool {
	sin := sin(angle)
	cos := cos(angle)
	cul := math.Pow(x*sin+y*cos-S, 2)/math.Pow(a, 2) + math.Pow(x*cos-y*sin-T, 2)/math.Pow(b, 2)
	return cul <= 1
}

/*
	求第二点相对第一点方向
*/
func PositionAzimuth(master *constant.Position, slave *constant.Position) float64 {
	cosc := cos(90-slave.Latitude)*cos(90-master.Latitude) + sin(90-slave.Latitude)*sin(90-master.Latitude)*cos(slave.Longitude-master.Longitude)
	sinc := math.Sqrt(1 - math.Pow(cosc, 2))
	A := ArcSin((sin(90-slave.Latitude) * sin(slave.Longitude-master.Longitude)) / sinc)
	if slave.Longitude >= master.Longitude && slave.Latitude >= master.Latitude {
		return A
	} else if slave.Longitude <= master.Longitude && slave.Latitude >= master.Latitude {
		return 360 + A
	} else {
		return 180 - A
	}
}

/*
	求最近会遇距离与会遇时间(平面近似，有误差)
*/
func CulMeetingIntersection(master *constant.Track, slave *constant.Track) *constant.MeetingIntersection {
	response := &constant.MeetingIntersection{}
	V0 := master.SOG
	C0 := master.COG
	Vt := slave.SOG
	Ct := slave.COG
	D := PositionSpacing(master.PrePosition, slave.PrePosition)
	q := PositionAzimuth(master.PrePosition, slave.PrePosition)
	Vr := math.Sqrt(math.Pow(V0, 2) + math.Pow(Vt, 2) - 2*V0*Vt*cos(C0-Ct))
	k := Vt / V0
	dH := C0 - Ct
	Cr := ArcCos((1 - k*cos(dH)) / (math.Sqrt(1 - 2*k*cos(dH) + math.Pow(k, 2))))
	response.DCPA = D * sin(Cr-q)
	response.TCPA = D * cos(Cr-q) / Vr
	return response
}

/*
	求第二点经纬度
	first 第一点位置
	L 距离
	Azimuth 方位角
*/
func CulSecondPointPosition(first *constant.Position, L float64, Azimuth float64) *constant.Position {
	var radius float64 = 6378137
	c := L / radius * 180 / math.Pi
	a := ArcCos(cos(90-first.Latitude)*cos(c) + sin(90-first.Latitude)*sin(c)*cos(Azimuth))
	C := ArcSin((sin(c) * sin(Azimuth)) / sin(a))
	return &constant.Position{
		Latitude:  90 - a,
		Longitude: first.Longitude + C,
	}
}
