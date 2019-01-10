package sql

import (
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetInfoWithShipID(t *testing.T) {
	defer func() {
		mysql.FreeConn()
	}()
	mysql.InitMysql()
	body, err := GetInfoWithShipID("412596777")
	assert.Empty(t, err, err)
	assert.Equal(t, 157, len(body))
}