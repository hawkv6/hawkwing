package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/pkg/client"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start Hawkwing client",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient("host-a")
		if err != nil {
			fmt.Println(err)
		}
		client.Start()
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		client.Stop()
		fmt.Println("\nHawkwing stopped")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
