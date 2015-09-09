package main

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eth-client/utils"

	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	HOST_IP   = "0.0.0.0"
	HOST_PORT = "8545"
	HOST      = fmt.Sprintf("%s:%s", HOST_IP, HOST_PORT)

	client *utils.Client
)

// override the hardcoded defaults with env variables if they're set
func init() {
	nodeAddr := os.Getenv("ETHTX_NODE_ADDR")
	if nodeAddr != "" {
		HOST = nodeAddr
	}
}

var (
	HostAddrFlag string

	// flags for `call` and `estimate`
	ToFlag    string
	FromFlag  string
	AmtFlag   string
	GasFlag   string
	PriceFlag string
	DataFlag  string
)

func main() {

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "check ethinfo version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("0.0.1")
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "ethinfo status",
		Long:  "print the node status",
		Run:   cliStatus,
	}

	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "ethinfo account <address>",
		Long:  "print an account",
		Run:   cliAccount,
	}

	var storageCmd = &cobra.Command{
		Use:   "storage",
		Short: "ethinfo storage <address> [storage key]",
		Long:  "print an account's storage or a single storage key",
		Run:   cliStorage,
	}

	var receiptCmd = &cobra.Command{
		Use:   "receipt",
		Short: "ethinfo reciept <tx hash>",
		Long:  "fetch the receipt for a transaction",
		Run:   cliReceipt,
	}

	var broadcastCmd = &cobra.Command{
		Use:   "broadcast",
		Short: "ethinfo broadcast <tx hex>",
		Long:  "broadcast a hex encoded rlp serialized transaction",
		Run:   cliBroadcast,
	}

	var estimateCmd = &cobra.Command{
		Use:   "estimate",
		Short: "ethinfo estimate [flags]",
		Long:  "estimate the gas required to run a transaction",
		Run:   cliEstimate,
	}
	estimateCmd.Flags().StringVarP(&ToFlag, "to", "t", "", "contract to call")
	estimateCmd.Flags().StringVarP(&FromFlag, "from", "f", "", "address to send from")
	estimateCmd.Flags().StringVarP(&AmtFlag, "amt", "a", "", "amt to send")
	estimateCmd.Flags().StringVarP(&GasFlag, "gas", "g", "", "gas to allocate for the call")
	estimateCmd.Flags().StringVarP(&PriceFlag, "price", "p", "", "price per unit of gas")
	estimateCmd.Flags().StringVarP(&DataFlag, "data", "d", "", "data to send the contract")

	var callCmd = &cobra.Command{
		Use:   "call",
		Short: "ethinfo call [flags]",
		Long:  "simulate calling a contract",
		Run:   cliCall,
	}
	callCmd.Flags().StringVarP(&ToFlag, "to", "t", "", "contract to call")
	callCmd.Flags().StringVarP(&FromFlag, "from", "f", "", "address to send from")
	callCmd.Flags().StringVarP(&AmtFlag, "amt", "a", "", "amt to send")
	callCmd.Flags().StringVarP(&GasFlag, "gas", "g", "", "gas to allocate for the call")
	callCmd.Flags().StringVarP(&PriceFlag, "price", "p", "", "price per unit of gas")
	callCmd.Flags().StringVarP(&DataFlag, "data", "d", "", "data to send the contract")

	var blocksCmd = &cobra.Command{
		Use:   "block",
		Short: "ethinfo block <number or hash>",
		Long:  "fetch a block by number or hash",
		Run:   cliBlocks,
	}

	var rootCmd = &cobra.Command{
		Use:   "ethinfo",
		Short: "a tool for talking to ethereum chains",
		Long:  "a tool for talking to ethereum chains",
	}
	rootCmd.PersistentFlags().StringVarP(&HostAddrFlag, "node-addr", "", HOST, "<ip>:<port> of the node we're talking to")

	rootCmd.PersistentPreRun = before

	rootCmd.AddCommand(
		versionCmd,
		statusCmd,
		accountCmd,
		storageCmd,
		broadcastCmd,
		receiptCmd,
		estimateCmd,
		callCmd,
		blocksCmd)
	rootCmd.Execute()
}

func before(cmd *cobra.Command, args []string) {
	HostAddrFlag = "http://" + HostAddrFlag
	client = utils.NewClient(HostAddrFlag)

}
