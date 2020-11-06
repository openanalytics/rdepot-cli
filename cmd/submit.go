package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"openanalytics.eu/rdepot/cli/client"
)

func init() {
	submitCmd.PersistentFlags().StringVarP(&repository, "repo", "r", "", "repository to upload to")
	submitCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "R package archive to upload")
	submitCmd.PersistentFlags().BoolVarP(&replace, "replace", "", true, "replace existing package version")
	rootCmd.AddCommand(submitCmd)
}

var (
	repository string
	filePath   string
	replace    bool

	submitCmd = &cobra.Command{
		Use:   "submit",
		Short: "Submit a package",
		Long:  `Submit a package to RDepot.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := client.RDepotConfig{
				Host:  Host,
				Token: Token,
			}
			res, err := client.SubmitPackage(cfg, filePath, repository, replace)
			if err != nil {
				return err
			}

			out, err := formatOutput(ByteArray(res))
			if err != nil {
				return err
			}

			fmt.Printf(out)
			return nil
		},
	}
)
