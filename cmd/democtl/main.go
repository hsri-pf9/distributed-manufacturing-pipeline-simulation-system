// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>

// */
// package main
// import (
// 	"log"
// 	// "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/cmd/democtl/cmd"
// 	"github.com/spf13/cobra"
// )

// func main() {
// 	var rootCmd = &cobra.Command{Use: "democtl"}
// 	rootCmd.AddCommand(startCmd)
// 	rootCmd.AddCommand(registerCmd)
// 	rootCmd.AddCommand(loginCmd)
// 	rootCmd.AddCommand(pipelineCmd)

// 	if err := rootCmd.Execute(); err != nil {
// 		log.Fatalf("Error executing command: %v", err)
// 	}
// }

package main

import (
	"log"
	"github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/cmd/democtl/cmd" // ✅ Import the command package
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
