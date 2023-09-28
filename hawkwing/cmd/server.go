package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Hawkwing server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := server.NewServer("host-b")
		if err != nil {
			fmt.Println(err)
		}
		server.Start()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		server.Stop()
		fmt.Println("\nHawkwing stopped")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
