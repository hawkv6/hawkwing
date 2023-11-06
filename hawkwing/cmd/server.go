package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/pkg/server"
	"github.com/spf13/cobra"
)

var serverInterface string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Hawkwing server",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := server.NewServer(serverInterface)
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
	serverCmd.Flags().StringVarP(&serverInterface, "interface", "i", "", "Interface to use for Hawkwing server")
	rootCmd.AddCommand(serverCmd)
}
