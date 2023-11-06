package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/pkg/client"
	"github.com/spf13/cobra"
)

var clientInterface string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start Hawkwing client",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClient(clientInterface)
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
	clientCmd.Flags().StringVarP(&clientInterface, "interface", "i", "", "Interface to use for Hawkwing client")
	rootCmd.AddCommand(clientCmd)
}
