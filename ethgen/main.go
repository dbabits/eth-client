package main

import (
	"github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/spf13/cobra"
)

/*
 Live Ethereum Genesis Params
 ----------------------------
 Difficulty: 0x400000000
 GasLimit: 0x1388
 ExtraData: 0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa
*/

var (
	CsvPathFlag string

	NonceFlag      string
	DifficultyFlag string
	ExtraDataFlag  string
	GasLimitFlag   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ethgen",
		Short: "a tool for generating ethereum genesis files",
		Long:  "a tool for generating ethereum genesis files",
		Run:   cliGen,
	}
	rootCmd.Flags().StringVarP(&CsvPathFlag, "csv", "c", "", "path to .csv where each line is (address, balance)")

	rootCmd.Flags().StringVarP(&NonceFlag, "nonce", "n", "0x0000000000000042", "genesis nonce")
	rootCmd.Flags().StringVarP(&DifficultyFlag, "difficulty", "d", "0x0fffff", "starting mining difficulty")
	rootCmd.Flags().StringVarP(&ExtraDataFlag, "extra-data", "x", "", "extra data for the genesis block")
	rootCmd.Flags().StringVarP(&GasLimitFlag, "gas-limit", "g", "0xffffffffffffffff", "starting gas limit per block")
	rootCmd.Execute()
}
