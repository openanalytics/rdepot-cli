package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
)

var (
  rootCmd = &cobra.Command{
    Use:   "rdepot",
    Short: "rdepot command line interface",
    Long: `RDepot is a solution of R package repository management.
                  Created and maintained by Open Analytics.
                  More information is available at http://rdepot.io`,
    Run: func(cmd *cobra.Command, args []string) {
      // Do Stuff Here
    },
  }
)

func Execute() error {
  return rootCmd.Execute()
}

func er(msg interface{}) {
  fmt.Println("Error:", msg)
  os.Exit(1)
}


