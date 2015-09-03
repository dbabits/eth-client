package main

import (
	"fmt"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	HOST_IP   = "0.0.0.0"
	HOST_PORT = "8545"
	HOST      = fmt.Sprintf("http://%s:%s", HOST_IP, HOST_PORT)

	SIGN_IP   = "0.0.0.0"
	SIGN_PORT = "4767"
	SIGN      = fmt.Sprintf("http://%s:%s", SIGN_IP, SIGN_PORT)

	// TODO: overwrite by env and flags
)

var (

	// all transactions take
	NonceFlag    int64
	AmtFlag      string
	GasFlag      string
	GasPriceFlag string

	// sign/broadcast/wait
	AddressFlag   string
	SignFlag      bool
	BroadcastFlag bool
	WaitFlag      bool

	// addresses
	HostAddrFlag string
	SignAddrFlag string

	// specifics
	ToFlag   string
	DataFlag string
)

func addCommonFlags(cmds []*cobra.Command) {
	for _, c := range cmds {
		c.Flags().Int64VarP(&NonceFlag, "nonce", "n", 0, "nonce for transaction")
		c.Flags().StringVarP(&AmtFlag, "amt", "a", "a", "amount to send")
		c.Flags().StringVarP(&GasFlag, "gas", "g", "", "amount of gas to provide")
		c.Flags().StringVarP(&GasPriceFlag, "price", "p", "", "price we're willing to pay per gas")

		c.Flags().StringVarP(&AddressFlag, "addr", "", "", "address to use for signing")
		c.Flags().BoolVarP(&SignFlag, "sign", "s", false, "sign the transaction")
		c.Flags().BoolVarP(&BroadcastFlag, "broadcast", "b", false, "broadcast the tx to the chain")
		c.Flags().BoolVarP(&WaitFlag, "wait", "w", false, "wait for the tx to be mined into a block")
	}
}

func main() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "check ethtx version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("0.0.1")
		},
	}

	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "ethtx send --amt <amt> --to <to>",
		Long:  "craft a simple send transaction",
		Run:   cliSend,
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "ethtx create --code <code>",
		Long:  "create a new contract",
		Run:   cliCreate,
	}

	var callCmd = &cobra.Command{
		Use:   "call",
		Short: "ethtx call --to <to> --data <data>",
		Long:  "call a contract",
		Run:   cliCall,
	}

	// custom flags
	sendCmd.Flags().StringVarP(&ToFlag, "to", "t", "", "destination address")
	callCmd.Flags().StringVarP(&ToFlag, "to", "t", "", "destination address")
	callCmd.Flags().StringVarP(&DataFlag, "data", "d", "", "data to send to the contract")
	createCmd.Flags().StringVarP(&DataFlag, "code", "c", "", "code for the new contract")

	commands := []*cobra.Command{sendCmd, createCmd, callCmd}
	addCommonFlags(commands)

	var rootCmd = &cobra.Command{
		Use:   "ethtx",
		Short: "a tool for sending transactions to ethereum chains",
		Long:  "a tool for sending transactions to ethereum chains",
	}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(commands...)
	rootCmd.Execute()
}
