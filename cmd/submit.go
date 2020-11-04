package cmd

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			cfg := client.RDepotConfig{
				Host:  Host,
				Token: Token,
			}
			client.SubmitPackage(cfg, filePath, repository, replace)
		},
	}
)
