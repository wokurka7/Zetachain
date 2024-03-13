package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
	"github.com/zeta-chain/zetacore/cmd/zetatool/filterdeposit"
)

var rootCmd = &cobra.Command{
	Use:   "zetatool",
	Short: "utility tool for zeta-chain",
}

func init() {
	rootCmd.AddCommand(filterdeposit.Cmd)
	rootCmd.PersistentFlags().String(config.Flag, "", "custom config file: --config filename.json")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}