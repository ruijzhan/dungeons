package app

import (
	"log"

	"github.com/spf13/cobra"
)

func NewApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Run:   run(),
	}
	return cmd
}

func run() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		server := NewServer(ServerOption{Addr: ":8080"})
		err := server.Run()
		if err != nil {
			log.Fatalf("server run failed: %v", err)
		}
	}
}
