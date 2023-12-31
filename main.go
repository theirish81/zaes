package main

import (
	"github.com/spf13/cobra"
	"zaes/lib"
)

// rootCmd is the root command
var rootCmd = &cobra.Command{
	Use:  "zaes [command]",
	Long: "Zaes is security utility that allows you to encrypt and securely erase files and directories",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(lib.EncryptCmd, lib.DecryptCmd, lib.WipeCmd)
	_ = rootCmd.Execute()
}
