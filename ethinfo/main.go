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

	var broadcastCmd = &cobra.Command{
		Use:   "broadcast",
		Short: "ethinfo broadcast <tx hex>",
		Long:  "broadcast a hex encoded rlp serialized transaction",
		Run:   cliBroadcast,
	}

	var rootCmd = &cobra.Command{
		Use:   "ethinfo",
		Short: "a tool for talking to ethereum chains",
		Long:  "a tool for talking to ethereum chains",
	}
	rootCmd.PersistentFlags().StringVarP(&HostAddrFlag, "node-addr", "", HOST, "<ip>:<port> of the node we're talking to")

	rootCmd.PersistentPreRun = before

	rootCmd.AddCommand(versionCmd, statusCmd, accountCmd, storageCmd, broadcastCmd)
	rootCmd.Execute()
}

func before(cmd *cobra.Command, args []string) {
	HostAddrFlag = "http://" + HostAddrFlag
	client = utils.NewClient(HostAddrFlag)

}
