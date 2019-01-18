package sql

import (
	"github.com/lupengyu/trafficflow/dal/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetPositionWithShipID(t *testing.T) {
	defer func() {
		mysql.FreeConn()
	}()
	mysql.InitMysql()
	body, err := GetPositionWithShipID("412596777")
	assert.Empty(t, err)
	assert.Equal(t, 570, len(body))
}
