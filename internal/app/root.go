package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "depsguard",
	Short: "depsguard enforces dependency explanations",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(explainCmd)
	rootCmd.AddCommand(auditCmd)
}
