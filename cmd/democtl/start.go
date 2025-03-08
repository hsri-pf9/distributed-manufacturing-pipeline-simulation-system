package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting gRPC server...")

		// ✅ Corrected variable name and method usage
		proc := exec.Command("./cmd/api-server/api-server") // ✅ Use ./ to indicate an executable
		err := proc.Start()
		if err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		} else {
			fmt.Println("gRPC server started successfully!")
		}
	},
}
