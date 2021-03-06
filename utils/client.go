package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/ethereum/go-ethereum/rpc/shared"
)

const JSONRPC = "2.0"

type Client struct {
	Host string
}

func NewClient(h string) *Client {
	return &Client{h}
}

func (c *Client) RequestResponse(api, method string, args ...interface{}) (interface{}, error) {
	var msg json.RawMessage
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	msg = json.RawMessage(b)
	request := &shared.Request{
		Jsonrpc: JSONRPC,
		Method:  fmt.Sprintf("%s_%s", api, method),
		Params:  msg,
		Id:      "",
	}
	body, err := c.requestResponse(request)
	if err != nil {
		return nil, err
	}
	r, err := unmarshalCheckError(body)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Client) requestResponse(s *shared.Request) (b []byte, err error) {
	if b, err = json.Marshal(s); err != nil {
		return nil, fmt.Errorf("Client side error: %v", err)
	}
	buf := bytes.NewBuffer(b)
	resp, err := http.Post(c.Host, "text/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func unmarshalCheckError(body []byte) (interface{}, error) {
	var errResponse shared.ErrorResponse
	var successResponse shared.SuccessResponse
	if err := json.Unmarshal(body, &errResponse); err == nil {
		if errResponse.Error != nil {
			return nil, fmt.Errorf("error code %d: %s", errResponse.Error.Code, errResponse.Error.Message)
		}
	}

	if err := json.Unmarshal(body, &successResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling success response", err)
	}
	return successResponse.Result, nil
}

func HexToInt(s string) int64 {
	d, _ := strconv.ParseInt(s, 0, 64)
	return d
}

func StripHex(s string) string {
	if strings.HasPrefix(s, "0x") {
		if len(s) > 2 {
			return s[2:]
		}
		return ""
	}
	return s
}
