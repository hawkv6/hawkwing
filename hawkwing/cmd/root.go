package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "hawkwing",
	Short: "Hawkwing brings SRv6 policies to your end-host.",
	Long: `
Hawkwing brings dynamic SRv6 policies to your end-host.

Start HawkWing in client-mode:
	hawkwing client --interface <interface> --config <config-file>

Start HawkWing in server-mode:
	hawkwing server --interface <interface>

Complete documentation is available at:
https://github.com/hawkv6/hawkwing
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	checkIsRoot()
}

func checkIsRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("Hawkwing must be run as root")
		os.Exit(1)
	}
}
