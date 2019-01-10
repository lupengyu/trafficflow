package http

import (
	"encoding/json"
	"github.com/lupengyu/trafficflow/constant"
)

func MakeGetShipRequest(option *RequestOption) (int, *constant.GetShipResponse, error) {
	resp, err := MakeRequest(
		GetShip,
		&RequestData{},
		option,
	)

	body := &constant.GetShipResponse{}
	err = json.NewDecoder(resp.Body).Decode(body)
	return resp.StatusCode, body, err
}
