package app

import (
	"github.com/spf13/cobra"
)

func NewServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE:  run(),
	}
	return cmd
}

func run() func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		opt := &Option{
			Addr: ":8080",
		}
		server := opt.newServer()

		return server.run()
	}
}
