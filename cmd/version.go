package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of rdepot-cli",
  Long:  `The version of rdepot-cli is increased in tandem with RDepot`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("rdepot-cli v1.4.1")
  },
}

