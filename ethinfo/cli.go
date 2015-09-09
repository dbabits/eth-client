package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/eris-ltd/eth-client/utils"

	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

//---------------------------------------------------------------
// ethinfo status

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
	Price    string `json:"gas_price"` // hex
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
	status.ChainStatus.BlockNumber = utils.HexToInt(r.(string))

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
	status.MiningStatus.Price = r.(string)

	b, err := json.MarshalIndent(status, "", "\t")
	common.IfExit(err)
	fmt.Println(string(b))
}

//---------------------------------------------------------------
// ethinfo account

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
	blockNum := utils.HexToInt(r.(string))

	r, err = client.RequestResponse("eth", "getBalance", addr, blockNum)
	common.IfExit(err)
	acc.Balance = r.(string)

	r, err = client.RequestResponse("eth", "getTransactionCount", addr, blockNum)
	common.IfExit(err)
	acc.Nonce = uint64(utils.HexToInt(r.(string)))

	r, err = client.RequestResponse("eth", "getCode", addr, blockNum)
	common.IfExit(err)
	acc.Code = r.(string)

	b, err := json.MarshalIndent(acc, "", "\t")
	common.IfExit(err)
	fmt.Println(string(b))
}

//---------------------------------------------------------------
// ethinfo storage

func cliStorage(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		common.Exit(fmt.Errorf("must specify an account"))
	}

	addr := args[0]
	var storageKey string
	if len(args) > 1 {
		storageKey = args[1]
	}

	r, err := client.RequestResponse("eth", "blockNumber")
	common.IfExit(err)
	blockNum := utils.HexToInt(r.(string))

	if storageKey == "" {
		// get all the storage
		r, err = client.RequestResponse("eth", "getStorage", addr, blockNum)
		common.IfExit(err)
		sortPrintMap(r.(map[string]interface{}))
	} else {
		// only grab one storage entry
		r, err = client.RequestResponse("eth", "getStorageAt", addr, storageKey, blockNum)
		common.IfExit(err)
		fmt.Println(r)
	}
}

//---------------------------------------------------------------
// ethinfo receipt

func cliReceipt(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		common.Exit(fmt.Errorf("must specify tx hash"))
	}

	txHash := args[0]

	r, err := client.RequestResponse("eth", "getTransactionReceipt", txHash)
	common.IfExit(err)
	sortPrintMap(r.(map[string]interface{}))

}

//---------------------------------------------------------------
// ethinfo broadcast

// TODO: this should be able to read off stdin too
func cliBroadcast(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		common.Exit(fmt.Errorf("must pass some transaction bytes"))
	}

	txHex := args[0]
	r, err := client.RequestResponse("eth", "sendRawTransaction", txHex)
	common.IfExit(err)
	fmt.Println(r)
}

//---------------------------------------------------------------
// ethinfo estimate

type callData struct {
	From     string
	To       string
	Value    interface{}
	Gas      interface{}
	GasPrice interface{}
	Data     string
}

func cliEstimate(cmd *cobra.Command, args []string) {
	r, err := client.RequestResponse("eth", "blockNumber")
	common.IfExit(err)
	blockNum := utils.HexToInt(r.(string))

	callArgs := callData{FromFlag, ToFlag, AmtFlag, GasFlag, PriceFlag, DataFlag}

	r, err = client.RequestResponse("eth", "estimateGas", callArgs, blockNum)
	common.IfExit(err)
	fmt.Println(r)
}

//---------------------------------------------------------------
// ethinfo call

func cliCall(cmd *cobra.Command, args []string) {
	r, err := client.RequestResponse("eth", "blockNumber")
	common.IfExit(err)
	blockNum := utils.HexToInt(r.(string))

	callArgs := callData{FromFlag, ToFlag, AmtFlag, GasFlag, PriceFlag, DataFlag}

	r, err = client.RequestResponse("eth", "call", callArgs, blockNum)
	common.IfExit(err)
	fmt.Println(r)
}

//---------------------------------------------------------------
// ethinfo blocks

func cliBlocks(cmd *cobra.Command, args []string) {
}

//---------------------------------------------------------------
// utils

type pair struct {
	key   string
	value interface{}
}

func sortPrintMap(m map[string]interface{}) {
	pairs := make([]pair, len(m))
	i := 0
	for k, v := range m {
		pairs[i] = pair{k, v}
		i += 1
	}
	sort.Sort(sortPairs(pairs))

	for _, p := range pairs {
		if p.value == nil {
			p.value = interface{}("")
		}
		fmt.Printf("%s: %s\n", p.key, p.value)
	}

}

type sortPairs []pair

func (v sortPairs) Len() int           { return len(v) }
func (v sortPairs) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v sortPairs) Less(i, j int) bool { return v[i].key < v[j].key }
