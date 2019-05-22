package helper

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/cnkei/gospline"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
)

func Test_LongitudeArea(t *testing.T) {
	assert.Equal(t, -1, LongitudeArea(0, 10))
	assert.Equal(t, -1, LongitudeArea(188, 10))
	assert.Equal(t, 9, LongitudeArea(118.999, 10))
	assert.Equal(t, 1, LongitudeArea(118, 10))
	assert.Equal(t, 5, LongitudeArea(118.5, 10))
}

func Test_LatitudeArea(t *testing.T) {
	assert.Equal(t, -1, LatitudeArea(0, 10))
	assert.Equal(t, -1, LatitudeArea(188, 10))
	assert.Equal(t, 0, LatitudeArea(24.1, 10))
	assert.Equal(t, 4, LatitudeArea(24.5, 10))
	assert.Equal(t, 7, LatitudeArea(24.8, 10))
}

func Test_DayDecrease(t *testing.T) {
	baseTime := &constant.Data{
		Year:   2019,
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	deltaTime := &constant.Data{
		Year:   0,
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 1,
		Second: 0,
	}
	t.Log(DayDecrease(baseTime, deltaTime))
}

func Test_DayIncrease(t *testing.T) {
	baseTime := &constant.Data{
		Year:   2019,
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	deltaTime := &constant.Data{
		Year:   0,
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 1,
		Second: 0,
	}
	t.Log(DayIncrease(baseTime, deltaTime))
}

func Test_IsLineInterSect(t *testing.T) {
	assert.Equal(t, true, IsLineInterSect(
		&constant.Position{Latitude: 0, Longitude: 0},
		&constant.Position{Latitude: 0, Longitude: 1},
		&constant.Position{Latitude: 0, Longitude: 1},
		&constant.Position{Latitude: 1, Longitude: 1},
	))
	assert.Equal(t, false, IsLineInterSect(
		&constant.Position{Latitude: 0, Longitude: 0},
		&constant.Position{Latitude: 0, Longitude: 1},
		&constant.Position{Latitude: 0, Longitude: 2},
		&constant.Position{Latitude: 1, Longitude: 1},
	))
	assert.Equal(t, true, IsLineInterSect(
		&constant.Position{Latitude: 0, Longitude: 0},
		&constant.Position{Latitude: 0, Longitude: 1},
		&constant.Position{Latitude: -1, Longitude: 0},
		&constant.Position{Latitude: 1, Longitude: 0},
	))
	assert.Equal(t, false, IsLineInterSect(
		&constant.Position{Latitude: 24.444706, Longitude: 118.04939},
		&constant.Position{Latitude: 24.41378, Longitude: 118.074398},
		&constant.Position{Latitude: 24.106986666666664, Longitude: 117.99400666666666},
		&constant.Position{Latitude: 24.503766666666667, Longitude: 118.63990333333332},
	))
}

func Test_TimeDeviation(t *testing.T) {
	assert.EqualValues(t, -10, TimeDeviation(
		&constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 0,
			Second: 0,
		}, &constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 0,
			Second: 10,
		},
	))
	assert.EqualValues(t, 40, TimeDeviation(
		&constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   1,
			Minute: 0,
			Second: 10,
		}, &constant.Data{
			Year:   2019,
			Month:  1,
			Day:    1,
			Hour:   0,
			Minute: 59,
			Second: 30,
		},
	))
}

func Test_PositionSpacing(t *testing.T) {
	t.Log(PositionSpacing(&constant.Position{
		Longitude: 118.148778,
		Latitude:  24.28328,
	}, &constant.Position{
		Longitude: 118.330451,
		Latitude:  24.388639,
	}))
	a := PositionSpacing(&constant.Position{
		Longitude: 0,
		Latitude:  0,
	}, &constant.Position{
		Longitude: 1,
		Latitude:  0,
	})
	b := PositionSpacing(&constant.Position{
		Longitude: 0,
		Latitude:  0,
	}, &constant.Position{
		Longitude: 0,
		Latitude:  1,
	})
	c := PositionSpacing(&constant.Position{
		Longitude: 0,
		Latitude:  0,
	}, &constant.Position{
		Longitude: 1,
		Latitude:  1,
	})
	t.Log(a, b, c)
	t.Log(math.Sqrt(a*a + b*b))
}

func Test_InEllipse(t *testing.T) {
	assert.True(t, InEllipse(5, 2.5, 0, 0, 0, 5, 0))
	assert.False(t, InEllipse(5, 2.5, 0, 0, 0, 6, 0))
	assert.True(t, InEllipse(5, 2.5, 0, 0, 2.5, 0, 0))
	assert.True(t, InEllipse(5, 2.5, 1, 0, 0, 5, 0))
	assert.False(t, InEllipse(5, 2.5, 1, 0, 2.5, 0, 0))
	assert.True(t, InEllipse(5, 2.5, 0, 0, 3.535, 3.535, 45))
	assert.False(t, InEllipse(5, 2.5, 0, 0, 3.635, 3.635, 45))
	assert.True(t, InEllipse(5, 2.5, 1, 0, 4.242, 4.242, 45))
}

func Test_PositionAzimuth(t *testing.T) {
	assert.EqualValues(t, 51.82921409234056, PositionAzimuth(&constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}, &constant.Position{
		Longitude: 116.581918,
		Latitude:  40.020885,
	}))
	assert.EqualValues(t, 307.5677865347654, PositionAzimuth(&constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}, &constant.Position{
		Longitude: 116.203048,
		Latitude:  40.031051,
	}))
	assert.EqualValues(t, 121.06169999543708, PositionAzimuth(&constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}, &constant.Position{
		Longitude: 116.654357,
		Latitude:  39.796846,
	}))
	assert.EqualValues(t, 246.07882841128884, PositionAzimuth(&constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}, &constant.Position{
		Longitude: 116.180052,
		Latitude:  39.837192,
	}))
}

func Test_CulSecondPointPosition(t *testing.T) {
	first := &constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}
	second := &constant.Position{
		Longitude: 116.581918,
		Latitude:  40.020885,
	}
	D := PositionSpacing(first, second)
	q := PositionAzimuth(first, second)
	culSecond := CulSecondPointPosition(first, D, q)
	t.Log(second, culSecond)
}

func Test_CulMeetingIntersection(t *testing.T) {
	first := &constant.Position{
		Longitude: 118.06283333333334,
		Latitude:  24.481209166666666,
	}
	second := &constant.Position{
		Longitude: 118.06080466666667,
		Latitude:  24.47756788888889,
	}
	sog1 := 9.05
	sog2 := 6.786666666666667
	cog1 := 34.400000000000006
	cog2 := 13.88
	t.Log(PositionSpacing(first, second))
	resp := CulMeetingIntersection(&constant.Track{
		PrePosition: first,
		COG:         cog1,
		SOG:         sog1,
	}, &constant.Track{
		PrePosition: second,
		COG:         cog2,
		SOG:         sog2,
	})
	t.Log(resp.TCPA, resp.DCPA)
	newFirst := CulSecondPointPosition(first, sog1*resp.TCPA, cog1)
	newSecond := CulSecondPointPosition(second, sog2*resp.TCPA, cog2)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 0.5*sog1*resp.TCPA, cog1)
	newSecond = CulSecondPointPosition(second, 0.5*sog2*resp.TCPA, cog2)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, -1*sog1*resp.TCPA, cog1)
	newSecond = CulSecondPointPosition(second, -1*sog1*resp.TCPA, cog2)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 2*sog1*resp.TCPA, cog1)
	newSecond = CulSecondPointPosition(second, 2*sog2*resp.TCPA, cog2)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 3*sog1*resp.TCPA, cog1)
	newSecond = CulSecondPointPosition(second, 3*sog2*resp.TCPA, cog2)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
}

func Test_CulMeetingIntersectionOld(t *testing.T) {
	first := &constant.Position{
		Longitude: 116.403119,
		Latitude:  39.913385,
	}
	second := &constant.Position{
		Longitude: 116.581918,
		Latitude:  40.020885,
	}
	t.Log(PositionSpacing(first, second))
	resp := CulMeetingIntersection(&constant.Track{
		PrePosition: first,
		COG:         0,
		SOG:         10,
	}, &constant.Track{
		PrePosition: second,
		COG:         350,
		SOG:         10,
	})
	t.Log(resp.TCPA, resp.DCPA)
	newFirst := CulSecondPointPosition(first, 10*resp.TCPA, 0)
	newSecond := CulSecondPointPosition(second, 10*resp.TCPA, 350)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 20*resp.TCPA, 0)
	newSecond = CulSecondPointPosition(second, 20*resp.TCPA, 350)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 30*resp.TCPA, 0)
	newSecond = CulSecondPointPosition(second, 30*resp.TCPA, 350)
	t.Log(newFirst, newSecond)
	t.Log(PositionSpacing(newFirst, newSecond))
}

func Test_TrackSorter(t *testing.T) {
	mysql.InitMysql()
	Time := &constant.Data{
		Year:   2019,
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	DeltaT := &constant.Data{
		Year:   0,
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 1,
		Second: 0,
	}
	beginTime := DayDecrease(Time, DeltaT)
	endTime := DayIncrease(Time, DeltaT)
	rows, err := mysql.DB.Query(
		"select * from position where MMSI = ? and (year > ? or year = ? and (month > ? or month = ? and (day > ? or day = ? and (hour > ? or hour = ? and (minute > ? or minute = ? and second >= ?))))) and (year < ? or year = ? and (month < ? or month = ? and (day < ? or day = ? and (hour < ? or hour = ? and (minute < ? or minute = ? and second <= ?))))) order by id",
		413694190,
		beginTime.Year, beginTime.Year,
		beginTime.Month, beginTime.Month,
		beginTime.Day, beginTime.Day,
		beginTime.Hour, beginTime.Hour,
		beginTime.Minute, beginTime.Minute,
		beginTime.Second,
		endTime.Year, endTime.Year,
		endTime.Month, endTime.Month,
		endTime.Day, endTime.Day,
		endTime.Hour, endTime.Hour,
		endTime.Minute, endTime.Minute,
		endTime.Second,
	)
	if err != nil {
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	shipTrackList := make([]*constant.Track, 0)
	for rows.Next() {
		// 数据绑定
		var pos constant.PositionMeta
		err := rows.Scan(
			&pos.ID, &pos.MessageType, &pos.RepeatIndicator, &pos.MMSI, &pos.NavigationStatus, &pos.ROT, &pos.SOG,
			&pos.PositionAccuracy, &pos.Longitude, &pos.Latitude, &pos.COG, &pos.HDG, &pos.TimeStamp, &pos.ReservedForRegional,
			&pos.RAIMFlag, &pos.Year, &pos.Month, &pos.Day, &pos.Hour, &pos.Minute, &pos.Second,
		)
		if err != nil {
			return
		}
		fmt.Println(pos)
		// 判断船舶位置
		nowPosition := &constant.Position{
			Longitude: pos.Longitude,
			Latitude:  pos.Latitude,
		}
		nowTime := &constant.Data{
			Year:   pos.Year,
			Month:  pos.Month,
			Day:    pos.Day,
			Hour:   pos.Hour,
			Minute: pos.Minute,
			Second: pos.Second,
		}
		track := &constant.Track{
			PrePosition: nowPosition,
			Time:        nowTime,
			Deviation:   TimeDeviation(nowTime, Time),
			COG:         pos.COG,
			SOG:         pos.SOG,
		}
		fmt.Println(track)
		shipTrackList = append(shipTrackList, track)
	}
	sorter := &TrackSorter{tracks: shipTrackList}
	sort.Sort(sorter)
	fmt.Println("====================")
	for _, v := range sorter.tracks {
		fmt.Println(v)
	}
	fmt.Println("====================")
	sorter.DeWeighting()
	for _, v := range sorter.tracks {
		fmt.Println(v, v.PrePosition)
	}
	track1 := sorter.tracks[0]
	track2 := sorter.tracks[1]
	diff := track1.Deviation - track2.Deviation
	longitudeK := (track1.PrePosition.Longitude - track2.PrePosition.Longitude) / float64(diff)
	latitudeK := (track1.PrePosition.Latitude - track2.PrePosition.Latitude) / float64(diff)
	cogK := (track1.COG - track2.COG) / float64(diff)
	sogK := (track1.SOG - track2.SOG) / float64(diff)
	a := constant.Track{
		PrePosition: &constant.Position{
			Longitude: track1.PrePosition.Longitude - float64(track1.Deviation)*longitudeK,
			Latitude:  track1.PrePosition.Latitude - float64(track1.Deviation)*latitudeK,
		},
		COG: track1.COG - float64(track1.Deviation)*cogK,
		SOG: track1.SOG - float64(track1.Deviation)*sogK,
	}
	fmt.Println(a, a.PrePosition)
}

func Test_PositionRelativeAzimuth(t *testing.T) {
	first := &constant.Position{
		Longitude: 118.0732,
		Latitude:  24.504341114856814,
	}
	second := &constant.Position{
		Longitude: 118.07406773809524,
		Latitude:  24.504341112448856,
	}
	cog := 360.0
	a := PositionRelativeAzimuth(first, cog, second)
	fmt.Println(a)
}

func Test_EllipseR(t *testing.T) {
	t.Log(EllipseR(5, 2.5, 0, 1, 0))
	t.Log(EllipseR(5, 2.5, 0, 1, 90))
	t.Log(EllipseR(5, 2.5, 0, 1, 180))
	t.Log(EllipseR(5, 2.5, 0, 1, 270))
	t.Log(EllipseR(5, 2.5, 0, 0, 45))
	t.Log(EllipseR(5, 2.5, 0, 0, 135))
	t.Log(EllipseR(5, 2.5, 0, 0, 225))
	t.Log(EllipseR(5, 2.5, 0, 0, 315))
}

func Test_PositionRelativeAzimuthVector(t *testing.T) {
	t.Log(PositionRelativeVector(1, 0, 2))
	t.Log(PositionRelativeVector(1, 180, 2))
	t.Log(PositionRelativeVector(3, 90, 4))
}

func Test_Danger(t *testing.T) {
	//cog0 := 0
	//sog0 := 15
	//L := 75
	//Azimuth := 29.5
	//D := 3 * constant.NauticalMile
	//DCPA := 0.4 * constant.NauticalMile
	//TCPA := 7
}

func Test_DangerUV(t *testing.T) {
	t.Log(MeetingDangerUV(0.01))
	t.Log(MeetingDangerUV(0.03))
	t.Log(MeetingDangerUV(1))
	t.Log(MeetingDangerUV(2))
	t.Log(MeetingDangerUV(3))
	t.Log(MeetingDangerUV(4))
	t.Log(MeetingDangerUV(5))
	t.Log(MeetingDangerUV(6))
	t.Log(MeetingDangerUV(7))
	t.Log(MeetingDangerUV(8))
	t.Log(MeetingDangerUV(9))
	t.Log(MeetingDangerUV(10))
	t.Log(MeetingDangerUV(9.05))
}

func Test_Boundary(t *testing.T) {
	t.Log(BoundaryR(0))
	t.Log(BoundaryR(30))
	t.Log(BoundaryR(60))
	t.Log(BoundaryR(90))
	t.Log(BoundaryR(120))
	t.Log(BoundaryR(180))
	t.Log(BoundaryR(270))
}

func Test_InterSpline(t *testing.T) {
	s := gospline.NewCubicSpline([]float64{0, 1}, []float64{0, 1})
	fmt.Println(s.At(0), s.At(1), s.At(2), s.At(3), s.At(4), s.At(5))
}

func Test_math(t *testing.T) {
	x := -400.2
	x = x - float64(int(x/360)-1)*360.0
	fmt.Println(x)
	x = 400.2
	x = x - float64(int(x/360))*360.0
	fmt.Println(x)
}

func Test_MaxRate(t *testing.T) {
	fmt.Println(MaxRate(110, 8))
	fmt.Println(MaxRate(110, 0))
	fmt.Println(MaxRate(110, 4))
}

func Test_MaxAcceleration(t *testing.T) {
	fmt.Println(MaxAcceleration(110, 16))
}

func Test_RateRange(t *testing.T) {
	fmt.Println(RateRange(20, 340))
	fmt.Println(RateRange(340, 300))
	fmt.Println(RateRange(300, 340))
	fmt.Println(RateRange(340, 20))
}

func Test_List(t *testing.T) {
	List := list.New()
	List.PushBack(1)
	List.PushBack(2)
	List.PushBack(3)
	for e := List.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	List.Remove(List.Front())
	List.PushBack(4)
	for e := List.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func Test_TrafficOutput(t *testing.T) {
	fmt.Print("[")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print("[", i, ",", j, ",", rand.Int31n(5), "],")
		}
	}
	fmt.Print("]")
}

func Test_AvailableData(t *testing.T) {
	file, err := os.Create("data/create.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	file.Sync()
	writer := bufio.NewWriter(file)
	data := &AvailableDataType{
		Longitude: 118.080862,
		Latitude:  24.487965,
		SOG:       0.1,
		COG:       350,
		Length:    110,
	}
	for i := 0; i < 10; i++ {
		data.VMin = 1
		data.VMax = 3
		data.RateTurn = 0.5
		str := strconv.FormatFloat(data.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(data.Latitude, 'f', -1, 64)
		n, err := writer.WriteString(str + "\r\n")
		if n != len(str) && err != nil {
			log.Println(err)
		}
		writer.Flush()

		fmt.Println("=============")
		fmt.Println(data)
		data = AvailableDataTest(data, 10)
		fmt.Println(data)
	}
}

/*
	118.07814474308265,24.48793123327711,5.3679111111111135,239.9249450114566
	118.07790359198313,24.487804144517256,5.607361616161619,232.57263588755245->118.07790359198313,24.487804144517256,5.607361616161619,224.57263588755245

	118.09267189875042,24.440384238118966,8.600492929292939,153.03150875933466
	118.09286850853842,24.440032467273255,8.480767676767686,154.17577834900646->118.09286850853842,24.440032467273255,8.480767676767686,177.17577834900646
*/

func Test_1(t *testing.T) {
	//Vnm := (5.3679111111111135 + 5.607361616161619) / 2
	//maxRate := MaxRate(110, Vnm)
	//fmt.Println("zhang rate:", maxRate)
	//
	//a := MaxAcceleration(110, 16.0)
	//preV := (5.3679111111111135 * 1.852) / 3.6
	//maxV := preV + float64(10)*a
	//V := (maxV + preV) / 2
	//Vnm = V * 3.6 / 1.852
	//maxRate = MaxRate(110, Vnm)
	//fmt.Println("my rate:", maxRate)

	Vnm := (8.600492929292939 + 8.480767676767686) / 2
	maxRate := MaxRate(110, Vnm)
	fmt.Println("zhang rate:", maxRate)

	a := MaxAcceleration(110, 16.0)
	preV := (8.600492929292939 * 1.852) / 3.6
	maxV := preV + float64(10)*a
	V := (maxV + preV) / 2
	Vnm = V * 3.6 / 1.852
	maxRate = MaxRate(110, Vnm)
	fmt.Println("my rate:", maxRate)
}

//118.07376715856918,24.482426457322347,10.156921212121214,199.4075012994731
//118.07298950337592,24.482551247981604,10.03719595959596,198.05470218805002
func Test_2(t *testing.T) {
	endPosition := CulSecondPointPosition(&constant.Position{Latitude: 24.482426457322347, Longitude: 118.07376715856918}, 70, 150)
	fmt.Println(endPosition)
	//&{118.07349571696271 24.481998576958674}
	//&{118.07298950337592 24.482551247981604}
}

func Test_Simulation1(t *testing.T) {
	angle := 270.0 // 相对角度
	L := 50.0      // 参照船船长
	D := 500.0     // 距离
	V := 4.0       // 速度
	DCPA := 500.0
	TCPA := 0.0
	Vr := 2.0578 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation21(t *testing.T) {
	angle := 180.0 // 相对角度
	L := 50.0      // 参照船船长
	D := 500.0     // 距离
	V := 4.0       // 速度
	DCPA := 0.0
	TCPA := 242.9779
	Vr := 2.0578 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation23(t *testing.T) {
	angle := 0.0 // 相对角度
	L := 50.0    // 参照船船长
	D := 500.0   // 距离
	V := 2.0     // 速度
	DCPA := 0.0
	TCPA := -242.9779
	Vr := 2.0578 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation31(t *testing.T) {
	angle := 45.0 // 相对角度
	L := 50.0     // 参照船船长
	D := 707.11   // 距离
	V := 16.0     // 速度
	DCPA := 500.0
	TCPA := 60.75
	Vr := 8.2311 // 相对速度 1852
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 90.0, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 90, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation32(t *testing.T) {
	angle := 90.0 // 相对角度
	L := 50.0     // 参照船船长
	D := 500.0    // 距离
	V := 16.0     // 速度
	DCPA := 500.0
	TCPA := 0.0
	Vr := 8.2311 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, angle, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation33(t *testing.T) {
	angle := 135.0 // 相对角度
	L := 50.0      // 参照船船长
	D := 707.11    // 距离
	V := 16.0      // 速度
	DCPA := 500.0
	TCPA := -242.9779
	Vr := 8.2311 // 相对速度 1852
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 90, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 90, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation41(t *testing.T) {
	l := math.Sqrt(math.Pow(1000.0+750*math.Sqrt(2), 2.0) + math.Pow(250*math.Sqrt(2), 2.0))
	T := math.Atan(250 * math.Sqrt(2) / (1000.0 + 750*math.Sqrt(2)))
	fmt.Println(l, T*180.0/math.Pi)
	angle := 9.735610317245346 // 相对角度
	L := 50.0                  // 参照船船长
	D := 2090.7702751760276    // 距离
	V := 14.78                 // 速度
	DCPA := 500.0
	TCPA := 242.98
	Vr := 7.5952 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation42(t *testing.T) {
	//response := &constant.MeetingIntersection{}
	//V0 := 8.0
	//C0 := 0.0
	//Vt := 8.0
	//Ct := 225.0
	//D := 500.0
	//q := 0.0
	//Vr := math.Sqrt(math.Pow(V0, 2) + math.Pow(Vt, 2) - 2*V0*Vt*cos(C0-Ct)) * 1852.0 / 3600.0
	//k := Vt / V0
	//dH := C0 - Ct
	//Cr := ArcCos((1 - k*cos(dH)) / (math.Sqrt(1 - 2*k*cos(dH) + math.Pow(k, 2))))
	//response.DCPA = D * sin(Cr-q)
	//response.TCPA = D * cos(Cr-q) / Vr
	//response.VR = Vr
	//fmt.Println(response.DCPA, response.TCPA, response.VR)

	angle := 0.0 // 相对角度
	L := 50.0    // 参照船船长
	D := 1207.11 // 距离
	V := 14.78   // 速度
	DCPA := 500.0
	TCPA := 121.62
	Vr := 7.5952 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_Simulation43(t *testing.T) {
	angle := 270.0 // 相对角度
	L := 50.0      // 参照船船长
	D := 638.22    // 距离
	V := 14.78     // 速度
	DCPA := 500.0
	TCPA := -50.322921994571
	Vr := 7.5952 // 相对速度
	rDCPA := MeetingDangerUDCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA)
	rTCPA := MeetingDangerUTCPA(5*L, 2.5*L, 0.75*L, 1.1*L, 315, DCPA, TCPA, Vr)
	rD := MeetingDangerUD(5*L, 2.5*L, 0.75*L, 1.1*L, angle, D)
	rQ := MeetingDangerUB(angle)
	rV := MeetingDangerUV(V)
	fmt.Println("rDCPA:", rDCPA)
	fmt.Println("rTCPA:", rTCPA)
	fmt.Println("rD   :", rD)
	fmt.Println("rQ   :", rQ)
	fmt.Println("rV   :", rV)
	fmt.Println("Score:", rD*0.4+rDCPA*rTCPA*0.4+rQ*0.1+rV*0.1)
}

func Test_test(t *testing.T) {
	fmt.Println(ArcSin(1.0 / 2.613125929752753))
}

func Test_driftAvailable(t *testing.T) {
	nowV := 8.0 * 1852.0 / 3600.0
	endV := 8.0 * 1852.0 / 3600.0
	maxA := MaxAcceleration(110, 16.0)
	minA := MinAcceleration(110, 16.0)
	diff := 10.0
	maxL := nowV*diff + 0.5*maxA*diff*diff
	maxV := nowV + maxA*diff
	time := (endV - nowV - minA*diff) / (maxA - minA)
	l := (maxV - endV) * (diff - time) * 0.5
	L := maxL - l
	fmt.Println(maxL, l, L, time, maxA, minA)
}

func Test_rateAvailable(t *testing.T) {
	nowV := 8.0 * 1852.0 / 3600.0
	endV := 8.0 * 1852.0 / 3600.0
	maxA := MaxAcceleration(110, 16.0)
	minA := MinAcceleration(110, 16.0)
	diff := 10.0
	maxL := nowV*diff + 0.5*maxA*diff*diff
	maxV := nowV + maxA*diff
	time := (endV - nowV - minA*diff) / (maxA - minA)
	l := (maxV - endV) * (diff - time) * 0.5
	L := maxL - l
	maxW := 360.0 * L / (2 * math.Pi * 110.0)
	fmt.Println(maxW)
}
