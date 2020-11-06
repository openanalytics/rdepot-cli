package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"openanalytics.eu/rdepot/cli/model"
)

var (
	rootCmd = &cobra.Command{
		Use:   "rdepot",
		Short: "rdepot command line interface",
		Long: `RDepot is a solution of R package repository management.

  More information is available at http://rdepot.io
  Open Analytics 2020`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		SilenceUsage: true,
	}
	Host   string
	Token  string
	output = "json"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "", "http://localhost", "RDepot host")
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "", "", "API token")
}

func Execute() error {
	return rootCmd.Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

type ByteArray []byte

func formatOutput(out model.Output) (string, error) {
	switch output {
	case "json":
		return out.FormatJSON()
	default:
		return "", fmt.Errorf("error type not supported: %s", output)
	}
}

func (o ByteArray) FormatJSON() (string, error) {
	return string(o), nil
}
