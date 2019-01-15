package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RequestInfo struct {
	Url         string
	Method      string
	ContentType string
}

type RequestName string

const (
	GetShip               RequestName = "GetShip"
	GetPositionWithShipID RequestName = "GetPositionWithShipID"
	GetInfoWithShipID     RequestName = "GetInfoWithShipID"
)

var (
	MRequestInfo = map[RequestName]RequestInfo{
		GetShip:               {"/api/get/ship/", "GET", "application/octet-stream"},
		GetPositionWithShipID: {"/api/get/position/", "GET", "application/octet-stream"},
		GetInfoWithShipID:     {"/api/get/info/", "GET", "application/octet-stream"},
	}
)

type RequestData struct {
	UrlPostfix  string
	QueryParams map[string]string
	HeaderData  map[string]string
	BodyData    map[string]interface{}
	Body        *bytes.Buffer
}

type RequestOption struct {
	Cookie     string
	Timeout    time.Duration
	ErrorMuted bool
}

func MakeRequest(requestName RequestName, requestData *RequestData, requestOption *RequestOption) (*http.Response, error) {
	reqInfo, ok := MRequestInfo[requestName]
	if !ok {
		return nil, errors.New("undefined request")
	}

	params := ""
	for k, v := range requestData.QueryParams {
		if len(params) == 0 {
			params += fmt.Sprintf("?%s=%s", k, v)
		} else {
			params += fmt.Sprintf("&%s=%s", k, v)
		}
	}
	url := domain + reqInfo.Url + requestData.UrlPostfix + params
	if requestData.HeaderData == nil {
		requestData.HeaderData = make(map[string]string)
	}
	switch reqInfo.Method {
	case "GET":
		return getRequest(url, requestData, requestOption.Timeout)
	case "POST":
		return postRequest(url, requestData, requestOption.Timeout)
	default:
		return nil, errors.New("undefined request type")
	}
}

func getRequest(url string, requestData *RequestData, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return makeRequestWithHeader(req, requestData.HeaderData, timeout)
}

func postRequest(url string, requestData *RequestData, timeout time.Duration) (*http.Response, error) {
	if len(requestData.BodyData) > 0 {
		bodyBytes, err := json.Marshal(requestData.BodyData)
		if err != nil {
			return nil, err
		}
		requestData.Body = bytes.NewBuffer(bodyBytes)
		requestData.HeaderData["Content-Type"] = "application/json"
	}

	req, err := http.NewRequest("POST", url, requestData.Body)
	if err != nil {
		return nil, err
	}

	return makeRequestWithHeader(req, requestData.HeaderData, timeout)
}

func makeRequestWithHeader(request *http.Request, headerData map[string]string, timeout time.Duration) (*http.Response, error) {
	for k, v := range headerData {
		request.Header.Set(k, v)
	}

	client := http.Client{Timeout: timeout}
	resp, err := client.Do(request)
	return resp, err
}
