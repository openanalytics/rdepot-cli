// Copyright 2020 Open Analytics
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func formatMessage(msg client.Message, strict bool) (string, error) {
	if class, err := msg.Class(); err != nil {
		return "", err
	} else {
		switch class {
		case "warning":
			if strict {
				return "", fmt.Errorf(msg.Content())
			} else {
				return fmt.Sprintf("Warning: %s", msg.Content()), nil
			}
		case "error":
			return "", fmt.Errorf(msg.Content())
		default:
			return msg.Content(), nil
		}
	}
}

type ByteArray []byte

func formatOutput(out model.Output) (string, error) {
	switch output {
	case "json":
		if res, err := model.FormatJSON(out); err != nil {
			return "", err
		} else {
			return string(res), nil
		}
	default:
		return "", fmt.Errorf("error type not supported: %s", output)
	}
}

func (o ByteArray) FormatJSON() ([]byte, error) {
	return o, nil
}
