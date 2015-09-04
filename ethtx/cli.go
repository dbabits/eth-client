package main

import (
	"fmt"

	"github.com/eris-ltd/eth-client/client"
	"github.com/eris-ltd/eth-client/ethtx/core"

	"github.com/eris-ltd/common/go/common"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

// TODO
func init() {
	client.HOST = HOST
}

func cliSend(cmd *cobra.Command, args []string) {
	tx, err := core.Send(HOST, AddressFlag, ToFlag, AmtFlag, GasFlag, GasPriceFlag, NonceFlag)
	common.IfExit(err)
	fmt.Println(tx)
	r, err := core.SignAndBroadcast(HOST, SIGN, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	fmt.Println(tx)
	fmt.Println("TxID", r)
}

func cliCreate(cmd *cobra.Command, args []string) {
	tx, err := core.Create(HOST, AddressFlag, AmtFlag, GasFlag, GasPriceFlag, DataFlag, NonceFlag)
	common.IfExit(err)
	fmt.Println(tx)
	r, err := core.SignAndBroadcast(HOST, SIGN, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	fmt.Println(tx)
	fmt.Println("TxID", r)
}

func cliCall(cmd *cobra.Command, args []string) {
	tx, err := core.Call(HOST, AddressFlag, ToFlag, AmtFlag, GasFlag, GasPriceFlag, DataFlag, NonceFlag)
	common.IfExit(err)
	fmt.Println(tx)
	r, err := core.SignAndBroadcast(HOST, SIGN, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	fmt.Println(tx)
	fmt.Println("TxID", r)
}

func cliName(cmd *cobra.Command, args []string) {
	fmt.Println("not implemented yet")
}
