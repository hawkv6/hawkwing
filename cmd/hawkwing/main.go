package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/internal/version"
	"github.com/hawkv6/hawkwing/pkg/client"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	versionFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "hawkwing",
	Short: "Hawkwing brings SRv6 policies to your end-host.",
	Long: `Hawkwing brings dynamic SRv6 policies to your end-host.
	Complete documentation is available at https://github.com/hawkv6/hawkwing
	`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print the version number of hawkwing")

	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println(version.GetVersion())
			os.Exit(0)
		}
	}

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

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
