package sql

import (
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetShip(t *testing.T) {
	defer func() {
		mysql.FreeConn()
	}()
	mysql.InitMysql()
	body, err := GetShip()
	assert.Empty(t, err, err)
	assert.Equal(t, 3844, len(body))
}