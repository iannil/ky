package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var HttpClient *resty.Client

func init() {
	HttpClient = resty.New()
}

var rootCmd = &cobra.Command{
	Use:   "ky",
	Short: "代劳无聊事",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	rootCmd.AddCommand(cmdWnacg)
	rootCmd.AddCommand(cmdWnacgl)
	rootCmd.AddCommand(cmdNyaa)
	rootCmd.AddCommand(cmdRemoteOK)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
