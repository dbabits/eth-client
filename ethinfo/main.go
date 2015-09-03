package main

import (
	"fmt"

	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	HOST_IP   = "0.0.0.0"
	HOST_PORT = "8545"
	HOST      = fmt.Sprintf("http://%s:%s", HOST_IP, HOST_PORT)
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

	var rootCmd = &cobra.Command{
		Use:   "ethinfo",
		Short: "a tool for talking to ethereum chains",
		Long:  "a tool for talking to ethereum chains",
	}
	rootCmd.AddCommand(versionCmd, statusCmd)
	rootCmd.Execute()
}
