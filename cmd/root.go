package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "database-connection-string-converter",
	Short: "A simple CLI tool to convert database connection strings",
}

func Execute() {
	_ = rootCmd.Execute()
}

func init() {
}
