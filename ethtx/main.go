package main

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eth-client/ethtx/core"
	"github.com/eris-ltd/eth-client/utils"

	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	HOST_IP   = "0.0.0.0"
	HOST_PORT = "8545"
	HOST      = fmt.Sprintf("%s:%s", HOST_IP, HOST_PORT)

	SIGN_IP   = "0.0.0.0"
	SIGN_PORT = "4767"
	SIGN      = fmt.Sprintf("%s:%s", SIGN_IP, SIGN_PORT)

	ADDR = ""

	client *utils.Client
)

// override the hardcoded defaults with env variables if they're set
func init() {
	signAddr := os.Getenv("ETHTX_SIGN_ADDR")
	if signAddr != "" {
		SIGN = signAddr
	}

	nodeAddr := os.Getenv("ETHTX_NODE_ADDR")
	if nodeAddr != "" {
		HOST = nodeAddr
	}

	addr := os.Getenv("ETHTX_ADDR")
	if addr != "" {
		ADDR = addr
	}
}

var (
	// logging
	LogLevelFlag int

	// all transactions take
	NonceFlag    uint64
	AmtFlag      string
	GasFlag      string
	GasPriceFlag string

	// sign/broadcast/wait
	AddressFlag   string
	BinaryFlag    bool
	SignFlag      bool
	BroadcastFlag bool
	WaitFlag      bool

	// http addresses
	HostAddrFlag string
	SignAddrFlag string

	// specifics
	ToFlag   string
	DataFlag string
)

func addCommonFlags(cmds []*cobra.Command) {
	for _, c := range cmds {
		c.Flags().Uint64VarP(&NonceFlag, "nonce", "n", 0, "nonce for transaction")
		c.Flags().StringVarP(&AmtFlag, "amt", "a", "", "amount to send")
		c.Flags().StringVarP(&GasFlag, "gas", "g", "", "amount of gas to provide")
		c.Flags().StringVarP(&GasPriceFlag, "price", "p", "", "price we're willing to pay per gas")
	}
}

func main() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "check ethtx version",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Println("0.0.2")
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

	// COMMANDS
	commands := []*cobra.Command{sendCmd, createCmd, callCmd}
	addCommonFlags(commands)

	var rootCmd = &cobra.Command{
		Use:   "ethtx",
		Short: "a tool for sending transactions to ethereum chains",
		Long:  "a tool for sending transactions to ethereum chains",
	}
	rootCmd.PersistentFlags().IntVarP(&LogLevelFlag, "log", "l", 0, "set the log level")
	rootCmd.PersistentFlags().StringVarP(&SignAddrFlag, "sign-addr", "", SIGN, "address to use for signing")
	rootCmd.PersistentFlags().StringVarP(&HostAddrFlag, "node-addr", "", HOST, "address to use for signing")
	rootCmd.PersistentFlags().StringVarP(&AddressFlag, "addr", "", ADDR, "address to use for signing")
	rootCmd.PersistentFlags().BoolVarP(&BinaryFlag, "binary", "", false, "print the tx's rlp serialized bytes (eg. to broadcast later)")
	rootCmd.PersistentFlags().BoolVarP(&SignFlag, "sign", "s", false, "sign the transaction")
	rootCmd.PersistentFlags().BoolVarP(&BroadcastFlag, "broadcast", "b", false, "broadcast the tx to the chain")
	rootCmd.PersistentFlags().BoolVarP(&WaitFlag, "wait", "w", false, "wait for the tx to be mined into a block")

	rootCmd.PersistentPreRun = before
	rootCmd.PersistentPostRun = after

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(commands...)
	rootCmd.Execute()
}

func before(cmd *cobra.Command, args []string) {
	SignAddrFlag = "http://" + SignAddrFlag
	HostAddrFlag = "http://" + HostAddrFlag
	core.EthClient = utils.NewClient(HostAddrFlag)

	log.SetLoggers(log.LogLevel(LogLevelFlag), os.Stdout, os.Stderr)
}

func after(cmd *cobra.Command, args []string) {
	log.Flush()
}
