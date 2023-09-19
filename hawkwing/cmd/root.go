package cmd

import (
	"fmt"
	"os"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	log     = logrus.New()
)

var rootCmd = &cobra.Command{
	Use:   "hawkwing",
	Short: "Hawkwing brings SRv6 policies to your end-host.",
	Long: `Hawkwing brings dynamic SRv6 policies to your end-host.
	Complete documentation is available at:
	https://github.com/hawkv6/hawkwing
	`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	client, err := client.NewClient("host-a")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	client.Start()
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// 	<-c

	// 	client.Stop()
	// 	fmt.Println("\nHawkwing stopped")
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	checkIsRoot()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		config.GetInstance().SetConfigFile(cfgFile)
	}

	if err := config.Parse(); err != nil {
		log.Fatalln(err)
	}
}

// TODO implement it proberly
func checkIsRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("Hawkwing must be run as root")
		os.Exit(1)
	}
}
