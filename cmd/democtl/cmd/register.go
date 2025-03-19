// package cmd

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
// 	"github.com/spf13/cobra"
// 	"google.golang.org/grpc"
// )

// var registerCmd = &cobra.Command{
// 	Use:   "register",
// 	Short: "Register a new user",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		email, _ := cmd.Flags().GetString("email")
// 		password, _ := cmd.Flags().GetString("password")

// 		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 		if err != nil {
// 			log.Fatalf("Failed to connect: %v", err)
// 		}
// 		defer conn.Close()

// 		client := proto.NewAuthServiceClient(conn)
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 		defer cancel()

// 		resp, err := client.Register(ctx, &proto.RegisterRequest{Email: email, Password: password})
// 		if err != nil {
// 			log.Fatalf("Registration failed: %v", err)
// 		}

// 		fmt.Printf("User registered successfully! UserID: %s, Email: %s\n", resp.UserId, resp.Email)
// 	},
// }

// func init() {
// 	registerCmd.Flags().String("email", "", "User Email")
// 	registerCmd.Flags().String("password", "", "User Password")
// 	registerCmd.MarkFlagRequired("email")
// 	registerCmd.MarkFlagRequired("password")
// }

package cmd

import (
	"fmt"
	"log"

	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		// 🔹 Get gRPC connection
		conn, ctx, cancel := GetGRPCConnection()
		defer conn.Close()
		defer cancel()

		client := proto.NewAuthServiceClient(conn)

		resp, err := client.Register(ctx, &proto.RegisterRequest{Email: email, Password: password})
		if err != nil {
			log.Fatalf("❌ Registration failed: %v", err)
		}

		fmt.Printf("✅ User registered! UserID: %s, Email: %s\n", resp.UserId, resp.Email)
	},
}

func init() {
	registerCmd.Flags().String("email", "", "User Email")
	registerCmd.Flags().String("password", "", "User Password")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")
}
