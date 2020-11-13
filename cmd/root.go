package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"openanalytics.eu/rdepot/cli/client"
	"openanalytics.eu/rdepot/cli/model"
)

var (
	rootCmd = &cobra.Command{
		Use:   "rdepot",
		Short: "rdepot command line interface",
		Long: `RDepot is a solution of R package repository management.

  More information is available at http://rdepot.io
  Open Analytics 2020`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			Config = client.RDepotConfig{
				Host:  viper.GetString("host"),
				Token: viper.GetString("token"),
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		SilenceUsage: true,
	}
	Host   string
	Token  string
	output = "json"

	Config client.RDepotConfig
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "", "http://localhost", "RDepot host")
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "", "", "API token")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.SetEnvPrefix("RDEPOT")
	viper.BindEnv("token")
	viper.BindEnv("host")
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
