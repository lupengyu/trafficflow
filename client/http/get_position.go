package http

import (
	"encoding/json"
	"github.com/lupengyu/trafficflow/constant"
)

func MakeGetPositionWithShipIDRequest(shipID string, option *RequestOption) (int, *constant.GetPositionWithShipIDResponse, error) {
	resp, err := MakeRequest(
		GetPositionWithShipID,
		&RequestData{
			UrlPostfix: shipID,
		},
		option,
	)

	body := &constant.GetPositionWithShipIDResponse{}
	err = json.NewDecoder(resp.Body).Decode(body)
	return resp.StatusCode, body, err
}
