package http

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_MakeGetInfoRequest(t *testing.T) {
	status, body, err := MakeGetInfoWithShipIDRequest("412596777", &RequestOption{Timeout: time.Duration(time.Hour)})
	assert.Empty(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, 0, body.Code)
	assert.Equal(t, "success", body.Message)
	assert.Equal(t, 157, len(body.Data))
}
