package main

import (
	"fmt"
	"os"

	"github.com/hawkv6/hawkwing/internal/version"
	"github.com/hawkv6/hawkwing/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println("Hello World")
		client.NewClient("host-a")
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
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
