package helper

import (
	"github.com/lupengyu/trafficflow/constant"
	"github.com/stretchr/testify/assert"
	"math"
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
	t.Log(resp)
	newFirst := CulSecondPointPosition(first, 10*resp.TCPA, 0)
	newSecond := CulSecondPointPosition(second, 10*resp.TCPA, 350)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 10*resp.TCPA+100, 0)
	newSecond = CulSecondPointPosition(second, 10*resp.TCPA+100, 350)
	t.Log(PositionSpacing(newFirst, newSecond))
	newFirst = CulSecondPointPosition(first, 10*resp.TCPA-100, 0)
	newSecond = CulSecondPointPosition(second, 10*resp.TCPA-100, 350)
	t.Log(PositionSpacing(newFirst, newSecond))
}
