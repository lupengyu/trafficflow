package helper

import (
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"sort"
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