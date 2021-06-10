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
	"openanalytics.eu/rdepot/cli/model"
)

func init() {
	packagesListCmd.Flags().StringVar(&nameFilter, "name", "", "filter by name glob pattern")
	packagesListCmd.Flags().StringVarP(&repositoryFilter, "repo", "r", "", "repository to filter with")
	packagesListCmd.Flags().BoolVar(&archivedFilter, "archived", false, "only list packages archived in the repository")
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

			pkgs, err := client.ListPackages(client.DefaultClient(), Config, repositoryFilter)
			if err != nil {
				return err
			}

			if archivedFilter {
				pkgs = model.FilterArchived(pkgs)
			}
			if nameFilter != "" {
				if pkgs, err = model.FilterByName(pkgs, nameFilter); err != nil {
					return err
				}
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
