package app

import (
	"log"

	"github.com/spf13/cobra"
)

func NewApp() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Run:   run(),
	}
	return serverCmd
}

func run() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Create server with default options
		server := NewServer(ServerOption{
			HttpListenAddr: ":8080",
			GrpcListenAddr: ":8081",
		})

		// Start the server
		if err := server.Run(); err != nil {
			log.Fatalf("server failed to run: %v", err)
		}
	}
}
