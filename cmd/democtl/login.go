package main

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "github.com/hsri-pf9/distributed-manufacturing-pipeline-simulation-system/api/grpc/proto/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login a user",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		client := proto.NewAuthServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.Login(ctx, &proto.LoginRequest{Email: email, Password: password})
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}

		fmt.Printf("User logged in! UserID: %s, Email: %s, Token: %s\n", resp.UserId, resp.Email, resp.Token)
	},
}

func init() {
	loginCmd.Flags().String("email", "", "User Email")
	loginCmd.Flags().String("password", "", "User Password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
}
