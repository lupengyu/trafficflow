package http

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_MakeGetShipRequest(t *testing.T) {
	status, getShipResponse, err := MakeGetShipRequest(&RequestOption{Timeout: time.Duration(time.Hour)})
	assert.Empty(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, 0, getShipResponse.Code)
	assert.Equal(t, "success", getShipResponse.Message)
	assert.Equal(t, 3844, len(getShipResponse.Data))
}