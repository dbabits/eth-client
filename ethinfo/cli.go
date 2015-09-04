package main

import (
	"encoding/json"
	"fmt"

	"github.com/eris-ltd/eth-client/client"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

// TODO ...
func init() {
	client.HOST = HOST
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

// eth: blockNumber, protocolVersion, coinbase, mining, gasPrice
// net: peerCount, listening, version
func cliStatus(cmd *cobra.Command, args []string) {
	var status Status

	r, err := client.RequestResponse("eth", "blockNumber")
	common.IfExit(err)
	status.ChainStatus.BlockNumber = client.HexToInt(r.(string))

	r, err = client.RequestResponse("eth", "protocolVersion")
	common.IfExit(err)
	status.ChainStatus.ProtocolVersion = r.(string)

	r, err = client.RequestResponse("net", "peerCount")
	common.IfExit(err)
	status.NetStatus.Version = r.(string)

	r, err = client.RequestResponse("net", "listening")
	common.IfExit(err)
	status.NetStatus.Listening = r.(bool)

	r, err = client.RequestResponse("net", "version")
	common.IfExit(err)
	status.NetStatus.Version = r.(string)

	r, err = client.RequestResponse("eth", "coinbase")
	common.IfExit(err)
	status.MiningStatus.Coinbase = r.(string)

	r, err = client.RequestResponse("eth", "mining")
	common.IfExit(err)
	status.MiningStatus.Mining = r.(bool)

	r, err = client.RequestResponse("eth", "gasPrice")
	common.IfExit(err)
	status.MiningStatus.GasPrice = r.(string)

	b, err := json.MarshalIndent(status, "", "\t")
	common.IfExit(err)
	fmt.Println(string(b))
}

type Account struct {
	Address     string `json:"address"`
	Nonce       uint64 `json:"nonce"`
	Balance     string `json:"balance"`
	Code        string `json:"code"`
	StorageHash string `json:"storage_hash"` // not sure this is supported but a storage map is (TODO)
}

func cliAccount(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		common.Exit(fmt.Errorf("must specify an account"))
	}

	addr := args[0]
	acc := new(Account)
	acc.Address = addr

	r, err := client.RequestResponse("eth", "blockNumber")
	common.IfExit(err)
	blockNum := client.HexToInt(r.(string))

	r, err = client.RequestResponse("eth", "getBalance", addr, blockNum)
	common.IfExit(err)
	acc.Balance = r.(string)

	r, err = client.RequestResponse("eth", "getTransactionCount", addr, blockNum)
	common.IfExit(err)
	acc.Nonce = uint64(client.HexToInt(r.(string)))

	r, err = client.RequestResponse("eth", "getCode", addr, blockNum)
	common.IfExit(err)
	acc.Code = r.(string)

	b, err := json.MarshalIndent(acc, "", "\t")
	common.IfExit(err)
	fmt.Println(string(b))
}
