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
)

func init() {
	packagesDeleteCmd.Flags().StringVar(&nameFilter, "name", "", "filter by name glob pattern")
	packagesDeleteCmd.Flags().StringVarP(&repositoryFilter, "repo", "r", "", "repository to filter with")
	packagesDeleteCmd.Flags().BoolVar(&archivedFilter, "archived", false, "only list packages archived in the repository")
	packagesDeleteCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "do not delete anyhing and just show what would be done")
	packagesCmd.AddCommand(packagesDeleteCmd)
}

var (
	dryRun            bool
	packagesDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete one or many packages",
		Long:  `Delete one or many packages`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if archivedFilter && repositoryFilter == "" {
				return fmt.Errorf(
					"archived filter can only be used when filtering by repository")
			}

			pkgs, err := client.ListPackages(client.DefaultClient(), Config, repositoryFilter, archivedFilter, nameFilter)
			if err != nil {
				return err
			}

			if dryRun {
				for _, pkg := range pkgs {
					fmt.Printf("would be deleted: %s\n", pkg.Summary())
				}
			} else {
				for _, pkg := range pkgs {
					err := client.DeletePackage(client.DefaultClient(), Config, pkg)
					if err != nil {
						return fmt.Errorf("could not delete package (%s): %v", pkg.Summary(), err)
					} else {
						fmt.Printf("deleted %s\n", pkg.Summary())
					}
				}
			}
			return nil
		},
	}
)
