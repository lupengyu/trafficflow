package http

import (
	"encoding/json"
	"github.com/lupengyu/trafficflow/constant"
)

func MakeGetInfoWithShipIDRequest(shipID string, option *RequestOption) (int, *constant.GetInfoWithShipIDResponse, error) {
	resp, err := MakeRequest(
		GetInfoWithShipID,
		&RequestData{
			UrlPostfix: shipID,
		},
		option,
	)

	body := &constant.GetInfoWithShipIDResponse{}
	err = json.NewDecoder(resp.Body).Decode(body)
	return resp.StatusCode, body, err
}