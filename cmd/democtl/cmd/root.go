package cmd

import (
	// "os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "democtl",
	Short: "CLI for managing pipelines",
	Long:  `CLI tool for interacting with the gRPC server`,
}

// Execute runs the root command and returns an error
func Execute() error { // ✅ Fix: Return error
	return rootCmd.Execute()
}


func init() {
	// ✅ Register subcommands here
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(pipelineCmd)
}
