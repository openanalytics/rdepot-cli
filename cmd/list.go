package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"openanalytics.eu/rdepot/cli/client"
)

var (
	handlers = map[string]func(client.RDepotConfig) ([]byte, error){
		"packages": client.ListPackages,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:       "list [resource]",
	Short:     "list one or many resources",
	Long:      `List one or many resources`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"packages"},
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
