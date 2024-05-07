// Copyright 2020-2024 Open Analytics
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
	"openanalytics.eu/rdepot/cli/model"
)

func init() {
	packagesListCmd.Flags().StringVar(&nameFilter, "name", "", "filter by name glob pattern")
	packagesListCmd.Flags().StringVarP(&repositoryFilter, "repo", "r", "", "repository to filter with")
	packagesListCmd.Flags().BoolVar(&archivedFilter, "archived", false, "return packages that do not have the latest version in a repository")
	packagesCmd.AddCommand(packagesListCmd)
}

var (
	nameFilter       string
	repositoryFilter string
	archivedFilter   bool

	packagesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List one or many packages",
		Long:  `List one or many packages`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if archivedFilter && repositoryFilter == "" {
				return fmt.Errorf(
					"archived filter can only be used when filtering by repository")
			}
			var pkgs model.Output
			var err error
			switch Config.Technology {
			case "r":
				pkgs, err = client.ListGenericPackages[model.RPackage](client.DefaultClient(), Config, repositoryFilter, archivedFilter, nameFilter)
			case "python":
				pkgs, err = client.ListGenericPackages[model.PythonPackage](client.DefaultClient(), Config, repositoryFilter, archivedFilter, nameFilter)
			case "all":
				pkgs, err = client.ListGenericPackages[model.Package](client.DefaultClient(), Config, repositoryFilter, archivedFilter, nameFilter)
			default:
				return fmt.Errorf("undefined technology %s", Config.Technology)
			}
			if err != nil {
				return err
			}

			if out, err := formatOutput(pkgs); err != nil {
				return err
			} else {
				fmt.Print(out)
				return nil
			}
		},
	}
)
