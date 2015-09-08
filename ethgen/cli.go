package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

//------------------------------------------------------------------------------
// ethgen cli

var (
	Mixhash    = "0x0000000000000000000000000000000000000000000000000000000000000000"
	Coinbase   = "0x0000000000000000000000000000000000000000"
	Timestamp  = "0x00"
	ParentHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

var DefaultBalance = "0xffffffffffffffffffffffffffffffffffffffff"

type Account struct {
	Balance string `json:"balance"`
}

type GenesisBlock struct {
	Nonce      string             `json:"nonce"`
	Timestamp  string             `json:"timestamp"`
	ParentHash string             `json:"parentHash"`
	ExtraData  string             `json:"extraData"`
	GasLimit   string             `json:"gasLimit"`
	Difficulty string             `json:"difficulty"`
	Mixhash    string             `json:"mixhash"`
	Coinbase   string             `json:"coinbase"`
	Alloc      map[string]Account `json:"alloc"`
}

func cliGen(cmd *cobra.Command, args []string) {
	g := &GenesisBlock{
		Nonce:      NonceFlag,
		Timestamp:  Timestamp,
		ParentHash: ParentHash,
		ExtraData:  ExtraDataFlag,
		GasLimit:   GasLimitFlag,
		Difficulty: DifficultyFlag,
		Mixhash:    Mixhash,
		Coinbase:   Coinbase,
	}

	if CsvPathFlag == "" && len(args) == 0 {
		Exit(fmt.Errorf("Please pass addresses as arguments or use the --csv flag"))
	}

	alloc := make(map[string]Account)

	for _, a := range args {
		alloc[a] = Account{DefaultBalance}
	}

	if CsvPathFlag != "" {
		addrs, balances, err := parseCsv(CsvPathFlag)
		IfExit(err)
		for i, addr := range addrs {
			alloc[addr] = Account{balances[i]}
		}
	}

	g.Alloc = alloc

	b, err := json.MarshalIndent(g, "", "\t")
	IfExit(err)
	fmt.Println(string(b))
}

func parseCsv(filePath string) (addrs []string, balances []string, err error) {

	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't open file: %s: %v", filePath, err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	params, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't read file: %v", err)
	}

	addrs = make([]string, len(params))
	balances = make([]string, len(params))
	for i, each := range params {
		addrs[i] = each[0]
		balances[i] = ifExistsElse(each, 1, DefaultBalance)
	}

	return addrs, balances, nil
}

func ifExistsElse(list []string, index int, defaultValue string) string {
	if len(list) > index {
		if list[index] != "" {
			return list[index]
		}
	}
	return defaultValue
}
