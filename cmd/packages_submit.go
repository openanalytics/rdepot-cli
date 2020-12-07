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
	"github.com/spf13/cobra"

	"openanalytics.eu/rdepot/cli/client"
)

func init() {
	packagesSubmitCmd.PersistentFlags().StringVarP(&repository, "repo", "r", "", "repository to upload to")
	packagesSubmitCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "R package archive to upload")
	packagesSubmitCmd.PersistentFlags().BoolVarP(&replace, "replace", "", true, "replace existing package version")
	packagesSubmitCmd.PersistentFlags().BoolVarP(&strict, "strict", "", true, "convert warnings into errors")
	packagesSubmitCmd.PersistentFlags().BoolVarP(&generateManuals, "generate-manuals", "", true, "generate a manual for the submitted package")
	packagesCmd.AddCommand(packagesSubmitCmd)
}

var (
	repository      string
	filePath        string
	strict          bool
	replace         bool
	generateManuals bool

	packagesSubmitCmd = &cobra.Command{
		Use:   "submit",
		Short: "Submit a package",
		Long:  `Submit a package to RDepot.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := client.SubmitPackage(client.DefaultClient(), Config, filePath, repository, replace, generateManuals)
			if err != nil {
				return err
			}

			msg, err := formatMessage(res, strict)
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("Package %s: %s", filePath, msg))
			return nil
		},
	}
)
