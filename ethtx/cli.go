package main

import (
	"github.com/eris-ltd/eth-client/ethtx/core"

	"github.com/eris-ltd/common/go/common"

	"github.com/spf13/cobra"
)

func cliSend(cmd *cobra.Command, args []string) {
	tx, err := core.Send(AddressFlag, ToFlag, AmtFlag, GasFlag, GasPriceFlag, NonceFlag)
	common.IfExit(err)
	logger.Infoln(tx)
	r, err := core.SignAndBroadcast(SignAddrFlag, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	logger.Infof("Signature %X\n", tx.Signature())
	if BroadcastFlag {
		logger.Println("TxID", r)
	}
}

func cliCreate(cmd *cobra.Command, args []string) {
	tx, err := core.Create(AddressFlag, AmtFlag, GasFlag, GasPriceFlag, DataFlag, NonceFlag)
	common.IfExit(err)
	logger.Infoln(tx)
	r, err := core.SignAndBroadcast(SignAddrFlag, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	logger.Infof("Signature %X\n", tx.Signature())
	if BroadcastFlag {
		logger.Println("TxID", r)
		logger.Printf("Address %X\n", tx.CreateAddress())
	}
}

func cliCall(cmd *cobra.Command, args []string) {
	tx, err := core.Call(AddressFlag, ToFlag, AmtFlag, GasFlag, GasPriceFlag, DataFlag, NonceFlag)
	common.IfExit(err)
	logger.Infoln(tx)
	r, err := core.SignAndBroadcast(SignAddrFlag, tx, SignFlag, BroadcastFlag, WaitFlag)
	common.IfExit(err)
	logger.Infof("Signature %X\n", tx.Signature())
	if BroadcastFlag {
		logger.Println("TxID", r)
	}
}

func cliName(cmd *cobra.Command, args []string) {
	logger.Errorln("not implemented yet")
}
