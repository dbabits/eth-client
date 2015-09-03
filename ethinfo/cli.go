package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc/shared"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

func RequestResponse(s *shared.Request) (b []byte, err error) {
	if b, err = json.Marshal(s); err != nil {
		return nil, fmt.Errorf("Client side error: %v", err)
	}
	buf := bytes.NewBuffer(b)
	resp, err := http.Post(HOST, "text/json", buf)
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

func hexToInt(s string) int64 {
	d, _ := strconv.ParseInt(s, 0, 64)
	return d
}

func stringToBool(s string) bool {
	d, _ := strconv.ParseBool(s)
	return d
}

// eth: blockNumber, protocolVersion, coinbase, mining, gasPrice
// net: peerCount, listening, version
func cliStatus(cmd *cobra.Command, args []string) {
	var status Status

	r, err := requestResponse("eth", "blockNumber")
	common.IfExit(err)
	status.ChainStatus.BlockNumber = hexToInt(r.(string))

	r, err = requestResponse("eth", "protocolVersion")
	common.IfExit(err)
	status.ChainStatus.ProtocolVersion = r.(string)

	r, err = requestResponse("net", "peerCount")
	common.IfExit(err)
	status.NetStatus.Version = r.(string)

	r, err = requestResponse("net", "listening")
	common.IfExit(err)
	status.NetStatus.Listening = r.(bool)

	r, err = requestResponse("net", "version")
	common.IfExit(err)
	status.NetStatus.Version = r.(string)

	r, err = requestResponse("eth", "coinbase")
	common.IfExit(err)
	status.MiningStatus.Coinbase = r.(string)

	r, err = requestResponse("eth", "mining")
	common.IfExit(err)
	status.MiningStatus.Mining = r.(bool)

	r, err = requestResponse("eth", "gasPrice")
	common.IfExit(err)
	status.MiningStatus.GasPrice = r.(string)

	b, err := json.MarshalIndent(status, "", "\t")
	common.IfExit(err)
	fmt.Println(string(b))

}

type ChainStatus struct {
	ProtocolVersion string `json:"protocol_version"`
	BlockNumber     int64  `json:"block_number"`
}

type NetStatus struct {
	Version   string `json:"version"`
	PeerCount int64  `json:"peer_count"`
	Listening bool   `json:"listening"`
}

type MiningStatus struct {
	Mining   bool   `json:"mining"`
	Coinbase string `json:"coinbase"`
	GasPrice string `json:"gas_price"` // hex
}

type Status struct {
	ChainStatus  `json:"chain"`
	NetStatus    `json:"net"`
	MiningStatus `json:"mining"`
}

var JSONRPC = "2.0"

func requestResponse(api, method string, args ...interface{}) (interface{}, error) {
	var msg json.RawMessage
	msg = json.RawMessage([]byte("{}")) // TODO
	request := &shared.Request{
		Jsonrpc: JSONRPC,
		Method:  fmt.Sprintf("%s_%s", api, method),
		Params:  msg,
		Id:      "",
	}
	body, err := RequestResponse(request)
	if err != nil {
		return nil, err
	}
	r, err := unmarshalCheckError(body)
	if err != nil {
		return nil, err
	}
	return r, err
}
