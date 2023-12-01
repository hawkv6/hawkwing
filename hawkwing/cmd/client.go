package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/client"
	"github.com/spf13/cobra"
)

var (
	clientInterface string
	cfgFile         string
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start Hawkwing client",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		mainErrCh := make(chan error)
		client, err := client.NewClient(mainErrCh, clientInterface)
		if err != nil {
			fmt.Println(err)
		}
		client.Start()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-sigChan:
			fmt.Println("\nReceived shutdown signal, exiting...")
		case err := <-mainErrCh:
			fmt.Printf("\nReceived error: %s, exiting...\n", err)
		}

		client.Stop()
		fmt.Println("\nHawkwing stopped")
	},
}

func init() {
	clientCmd.Flags().StringVarP(&clientInterface, "interface", "i", "", "Interface to use for Hawkwing client")
	clientCmd.Flags().StringVar(&cfgFile, "config", "./config.yaml", "config file (default is ./config.yaml)")
	rootCmd.AddCommand(clientCmd)
}

func initConfig() {
	if cfgFile != "" {
		config.GetInstance().SetConfigFile(cfgFile)
	}

	if err := config.Parse(); err != nil {
		log.Fatalln(err)
	}
}
