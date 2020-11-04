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

  More information is available at http://rdepot.io
  Open Analytics 2020`,
    Run: func(cmd *cobra.Command, args []string) {
      // Do Stuff Here
    },
  }
  Host string
  Token string
)

func init() {
  rootCmd.PersistentFlags().StringVarP(&Host, "host", "", "http://localhost", "RDepot host")
  rootCmd.PersistentFlags().StringVarP(&Token, "token", "", "", "API token")
}

func Execute() error {
  return rootCmd.Execute()
}

func er(msg interface{}) {
  fmt.Println("Error:", msg)
  os.Exit(1)
}


