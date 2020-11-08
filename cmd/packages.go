package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(packagesCmd)
}

var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "Perform package actions",
	Long:  `Perform package actions`,
	Run:   func(cmd *cobra.Command, args []string) {},
}
