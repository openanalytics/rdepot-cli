package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"openanalytics.eu/rdepot/cli/client"
)

func init() {
	packagesCmd.AddCommand(packagesListCmd)
}

var packagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List one or many packages",
	Long:  `List one or many packages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := client.RDepotConfig{
			Host:  Host,
			Token: Token,
		}

		if res, err := client.ListPackages(cfg); err != nil {
			return err
		} else {
			fmt.Printf(string(res))
			return nil
		}
	},
}
