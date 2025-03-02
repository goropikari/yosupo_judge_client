package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yosupocl",
	Short: "yosupo judge client",
	Long:  `yosupo judge api client`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(problemInfoCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(downloadTestCasesCmd)
}
