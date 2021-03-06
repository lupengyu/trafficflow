package helper

import (
	"fmt"
	"github.com/cnkei/gospline"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/cache"
	"math"
	"math/rand"
	"sort"
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
	if w*z > 0 || u*v > 0 {
		return false
	}
	return true
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
	lat1 := a.Latitude * rad  //纬度
	lng1 := a.Longitude * rad //经度
	lat2 := b.Latitude * rad
	lng2 := b.Longitude * rad
	theta := lng2 - lng1
	distTheta := math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta)
	if distTheta > 1 {
		distTheta = 1
	}
	dist := math.Acos(distTheta)
	return dist * radius * 1000
}

/*
	插值算法(三次样条插值)
*/
func TrackInterpolation(tracks []*constant.Track) *constant.Track {
	if len(tracks) == 0 {
		return nil
	} else if len(tracks) == 1 {
		return tracks[0]
	}
	sorter := &TrackSorter{tracks: tracks}
	sort.Sort(sorter)
	sorter.DeWeighting()
	if sorter.Len() == 1 {
		return sorter.tracks[0]
	}

	x := make([]float64, 0)
	longitudeY := make([]float64, 0)
	latitudeY := make([]float64, 0)
	cogY := make([]float64, 0)
	sogY := make([]float64, 0)
	preCOG := sorter.tracks[0].COG
	for _, v := range sorter.tracks {
		x = append(x, float64(v.Deviation))
		longitudeY = append(longitudeY, v.PrePosition.Longitude)
		latitudeY = append(latitudeY, v.PrePosition.Latitude)
		cogY = append(cogY, RateRange(v.COG, preCOG))
		sogY = append(sogY, v.SOG)
	}
	longitude := gospline.NewCubicSpline(x, longitudeY).At(0)
	latitude := gospline.NewCubicSpline(x, latitudeY).At(0)
	cog := gospline.NewCubicSpline(x, cogY).At(0) + preCOG
	sog := gospline.NewCubicSpline(x, sogY).At(0)
	sog = math.Max(0, sog)
	if cog < 0 {
		cog = cog - float64(int(cog/360)-1)*360.0
	} else if cog > 360 {
		cog = cog - float64(int(cog/360))*360.0
	}
	return &constant.Track{
		PrePosition: &constant.Position{
			Longitude: longitude,
			Latitude:  latitude,
		},
		COG: cog,
		SOG: sog,
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
	椭圆上一点到(-S, -T)的距离
*/
func EllipseR(a float64, b float64, S float64, T float64, angle float64) float64 {
	x0 := -S
	y0 := -T
	x := 0.0
	y := 0.0
	if angle == 0 || angle == 360 {
		x = -S
		y = math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		return y + T
	} else if angle == 180 {
		x = -S
		y = math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		return y - T
	} else if angle == 90 {
		y = -T
		x = math.Sqrt((math.Pow(b, 2)) * (1 - (math.Pow(y, 2) / math.Pow(a, 2))))
		return x + S
	} else if angle == 270 {
		y = -T
		x = math.Sqrt((math.Pow(b, 2)) * (1 - (math.Pow(y, 2) / math.Pow(a, 2))))
		return x - S
	} else if angle < 180 {
		k := cos(angle) / sin(angle)
		A := math.Pow(k, 2) + math.Pow(a, 2)/math.Pow(b, 2)
		B := 2 * k * (k*S - T)
		C := math.Pow(k*S, 2) - 2*k*T*S + math.Pow(T, 2) - math.Pow(a, 2)
		x = (-B + math.Sqrt(math.Pow(B, 2)-4*A*C)) / (2 * A)
		if angle <= 90 {
			y = math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		} else {
			y = -math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		}
	} else {
		k := cos(angle) / sin(angle)
		A := math.Pow(k, 2) + math.Pow(a, 2)/math.Pow(b, 2)
		B := 2 * k * (k*S - T)
		C := math.Pow(k*S, 2) - 2*k*T*S + math.Pow(T, 2) - math.Pow(a, 2)
		x = (-B + math.Sqrt(math.Pow(B, 2)-4*A*C)) / (2 * A)
		if angle > 270 {
			y = math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		} else {
			y = -math.Sqrt((math.Pow(a, 2)) * (1 - (math.Pow(x, 2) / math.Pow(b, 2))))
		}
	}
	length := math.Sqrt(math.Pow(x-x0, 2) + math.Pow(y-y0, 2))
	return length
}

/*
	求第二点相对第一点方向
*/
func PositionAzimuth(master *constant.Position, slave *constant.Position) float64 {
	cosc := cos(90-slave.Latitude)*cos(90-master.Latitude) + sin(90-slave.Latitude)*sin(90-master.Latitude)*cos(slave.Longitude-master.Longitude)
	sinc := math.Sqrt(1 - math.Pow(cosc, 2))
	A := ArcSin((sin(90-slave.Latitude) * sin(slave.Longitude-master.Longitude)) / sinc)
	if (sin(90-slave.Latitude)*sin(slave.Longitude-master.Longitude))/sinc > 1 {
		A = ArcSin(1)
	}
	if slave.Longitude >= master.Longitude && slave.Latitude >= master.Latitude {
		return A
	} else if slave.Longitude <= master.Longitude && slave.Latitude >= master.Latitude {
		return 360 + A
	} else {
		return 180 - A
	}
}

/*
	求第二点相对第一点速度的方向
*/
func PositionRelativeAzimuth(master *constant.Position, masterAzimuth float64, slave *constant.Position) float64 {
	azimuth := PositionAzimuth(master, slave)
	relativeAzimuth := azimuth - masterAzimuth
	if relativeAzimuth < 0 {
		return relativeAzimuth + 360
	}
	return relativeAzimuth
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
	response.VR = Vr
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

func PositionRelativeVector(v0 float64, RelativeAzimuth float64, vt float64) float64 {
	vx := v0 - vt*cos(RelativeAzimuth)
	vy := vt * sin(RelativeAzimuth)
	return math.Sqrt(math.Pow(vx, 2) + math.Pow(vy, 2))
}

func BoundaryR(angle float64) float64 {
	return constant.NauticalMile * (1.7*cos(angle-19) + math.Sqrt(4.4+2.89*math.Pow(cos(angle-19), 2)))
}

func MeetingDangerUDCPA(a float64, b float64, S float64, T float64, angle float64, DCPA float64) float64 {
	min := EllipseR(a, b, S, T, angle)
	//fmt.Println(min)
	max := constant.NauticalMile
	//fmt.Println(max)
	if DCPA < min {
		return 1.0
	} else if DCPA > max {
		return 0.0
	}
	return 0.5 - math.Sin(math.Pi*(DCPA-(min+max)/2)/(max-min))/2.0
}

func MeetingDangerUTCPA(a float64, b float64, S float64, T float64, angle float64, DCPA float64, TCPA float64, Vr float64) float64 {
	//min := EllipseR(a, b, S, T, angle)
	//max := constant.NauticalMile
	if TCPA < 0 {
		return 0.0
	}
	max := BoundaryR(angle)
	min := 12 * a / 5
	t1 := (DCPA - min) / Vr
	if min >= DCPA {
		t1 = math.Sqrt(math.Pow(min, 2)-math.Pow(DCPA, 2)) / Vr
	}
	t2 := math.Sqrt(math.Pow(max, 2)-math.Pow(DCPA, 2)) / Vr
	//fmt.Println(t1, t2, TCPA)
	if DCPA > max {
		return 0.0
	}
	if TCPA <= t1 {
		return 1.0
	} else if TCPA > t2 {
		return 0.0
	}
	return math.Pow((t2-TCPA)/(t2-t1), 2)
}

func MeetingDangerUB(angle float64) float64 {
	return (cos(angle-19)+math.Sqrt(440.0/289.0+math.Pow(cos(angle-19), 2)))/2.0 - 5.0/17.0
}

func MeetingDangerUD(a float64, b float64, S float64, T float64, angle float64, D float64) float64 {
	max := BoundaryR(angle)
	min := EllipseR(a, b, S, T, angle)
	if D <= min {
		return 1.0
	} else if D > max {
		return 0.0
	}
	return math.Pow((max-D)/(max-min), 2)
}

func MeetingDangerUV(v float64) float64 {
	min := constant.StaticShip
	max := 13.7624
	if v < min {
		return 0
	} else if v > max {
		return 1
	}
	return math.Pow((v-min)/(max-min), 2)
}

func MeetingDangerUT(v0 float64, vt float64, T float64) float64 {
	T0 := 90.0
	if v0 > vt {
		T0 = 180.0
	} else if v0 < vt {
		T0 = 40.0
	}
	if T == 360 {
		T = 0
	}
	if T >= 0 && T < 180 {
		return 1.0 / (1.0 + math.Pow(T/T0, 2))
	} else if T >= 180 && T < 360 {
		return 1.0 / math.Pow((360.0-T)/T0, 2)
	} else {
		return 0.0
	}
}

func MaxRate(length int, v float64) float64 {
	return 360.0 * v * 1.852 / (2.0 * 3.6 * math.Pi * float64(length))
}

func MaxAcceleration(length int, maxSpeed float64) float64 {
	maxSpeed = maxSpeed * 1.852 / 3.6
	return math.Pow(maxSpeed, 2) / float64(20*length)
}

func MinAcceleration(length int, maxSpeed float64) float64 {
	maxSpeed = maxSpeed * 1.852 / 3.6
	return -1 * math.Pow(maxSpeed, 2) / float64(16*length)
}

func movingAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	if prePosition.Longitude == position.Longitude &&
		prePosition.Latitude == position.Latitude &&
		prePosition.SOG > 2 {
		fmt.Println("moving", prePosition, position) // TODO: remove
		return false
	}
	return true
}

func driftAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	diff := TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	shipInfo := cache.GetShipInfo(position.MMSI)
	nowV := prePosition.SOG * 1852.0 / 3600.0
	endV := position.SOG * 1852.0 / 3600.0
	maxA := MaxAcceleration(shipInfo.Length, 16.0)
	minA := MinAcceleration(shipInfo.Length, 16.0)
	maxL := nowV*float64(diff) + 0.5*maxA*float64(diff)*float64(diff)
	maxV := nowV + maxA*float64(diff)
	time := (endV - nowV - minA*float64(diff)) / (maxA - minA)
	l := (maxV - endV) * (float64(diff) - time) * 0.5
	L := maxL - l
	D := PositionSpacing(&constant.Position{
		Latitude:  prePosition.Latitude,
		Longitude: prePosition.Longitude,
	}, &constant.Position{
		Latitude:  position.Latitude,
		Longitude: position.Longitude,
	})
	if D > L {
		fmt.Println("drift", L, D, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func accelerationAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	diff := TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	shipInfo := cache.GetShipInfo(position.MMSI)
	acturalA := (position.SOG - prePosition.SOG) * 1852.0 / 3600.0 / float64(diff) //单位 m/s^2
	aMax := MaxAcceleration(shipInfo.Length, 16.0)                                 //单位 m/s^2
	aMin := MinAcceleration(shipInfo.Length, 16.0)                                 //单位 m/s^2
	if acturalA > aMax || acturalA < aMin {
		fmt.Println("acceleration", acturalA, aMax, aMin, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func rateAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	rate := math.Max(position.COG-prePosition.COG, prePosition.COG-position.COG)
	absRate := math.Min(rate, 360.0-rate) // 当前变化角度
	diff := TimeDeviation(&constant.Data{
		Year:   position.Year,
		Month:  position.Month,
		Day:    position.Day,
		Hour:   position.Hour,
		Minute: position.Minute,
		Second: position.Second,
	}, &constant.Data{
		Year:   prePosition.Year,
		Month:  prePosition.Month,
		Day:    prePosition.Day,
		Hour:   prePosition.Hour,
		Minute: prePosition.Minute,
		Second: prePosition.Second,
	})
	shipInfo := cache.GetShipInfo(position.MMSI)
	nowV := prePosition.SOG * 1852.0 / 3600.0
	endV := position.SOG * 1852.0 / 3600.0
	maxA := MaxAcceleration(shipInfo.Length, 16.0)
	minA := MinAcceleration(shipInfo.Length, 16.0)
	maxL := nowV*float64(diff) + 0.5*maxA*float64(diff)*float64(diff)
	maxV := nowV + maxA*float64(diff)
	time := (endV - nowV - minA*float64(diff)) / (maxA - minA)
	l := (maxV - endV) * (float64(diff) - time) * 0.5
	L := maxL - l
	maxW := 360.0 * L / (2 * math.Pi * 110.0)
	if absRate > maxW {
		fmt.Println("rate", absRate, maxW, prePosition, position) // TODO: remove
		return false
	}
	return true
}

func PositionAvailable(position constant.PositionMeta, prePosition constant.PositionMeta) bool {
	return movingAvailable(position, prePosition) &&
		driftAvailable(position, prePosition) && accelerationAvailable(position, prePosition) &&
		rateAvailable(position, prePosition)
}

func RateRange(now float64, pre float64) float64 {
	rateRange := now - pre
	if rateRange < -180 {
		return rateRange + 360.0
	} else if rateRange > 180 {
		return rateRange - 360.0
	}
	return rateRange
}

type AvailableDataType struct {
	SOG       float64
	COG       float64
	Longitude float64
	Latitude  float64
	Length    int
	VMin      int
	VMax      int
	RateTurn  float64
}

func AvailableDataTest(data *AvailableDataType, Time int) *AvailableDataType {
	a := MaxAcceleration(data.Length, 16.0)
	preV := (data.SOG * 1.852) / 3.6
	maxVPrecent := preV + float64(Time)*a*float64(rand.Intn(data.VMax-data.VMin)+data.VMin)*0.1
	vnm := maxVPrecent * 3.6 / 1.852
	VPrecent := (maxVPrecent + preV) / 2
	VPrecentnm := VPrecent * 3.6 / 1.852
	maxRatePrecent := MaxRate(data.Length, VPrecentnm) * data.RateTurn
	endPosition := CulSecondPointPosition(&constant.Position{Latitude: data.Latitude, Longitude: data.Longitude}, VPrecent*float64(Time), data.COG)
	return &AvailableDataType{
		SOG:       vnm,
		COG:       data.COG - maxRatePrecent*float64(Time),
		Longitude: endPosition.Longitude,
		Latitude:  endPosition.Latitude,
		Length:    data.Length,
	}
}

func AvailableData(prePosition constant.PositionMeta, length int, time int, precent float64) {
	a := MaxAcceleration(length, 16.0)
	preV := (prePosition.SOG * 1.852) / 3.6
	maxV := preV + float64(time)*a
	maxVPrecent := preV + float64(time)*a*precent
	v := maxV * 3.6 / 1.852
	vnm := maxVPrecent * 3.6 / 1.852
	V := (maxV + preV) / 2
	VPrecent := (maxVPrecent + preV) / 2
	Vnm := V * 3.6 / 1.852
	VPrecentnm := VPrecent * 3.6 / 1.852
	maxRate := MaxRate(length, Vnm)
	maxRatePrecent := MaxRate(length, VPrecentnm)
	fmt.Println("Pre v km:", preV)
	fmt.Println("Max a   :", a)
	fmt.Println("Max v km:", maxV)
	fmt.Println("Max v nm:", v)
	fmt.Println("Max v nm precent=============:", vnm)
	fmt.Println("===================")
	fmt.Println("Max rate:", maxRate)
	fmt.Println("Max rate precent=============:", maxRatePrecent*float64(time))
	fmt.Println("===================")
	fmt.Println("Max L   :", V*float64(time))
	fmt.Println("Max L   precent=============:", VPrecent*float64(time))
	endPosition := CulSecondPointPosition(&constant.Position{Latitude: prePosition.Latitude, Longitude: prePosition.Longitude}, VPrecent, prePosition.COG)
	fmt.Println("Position:                   :", endPosition)
}
