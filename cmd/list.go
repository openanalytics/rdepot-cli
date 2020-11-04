package cmd

import (
  "fmt"

  "github.com/spf13/cobra"

  "openanalytics.eu/rdepot/cli/client"
)

var (
  handlers = map[string]func(client.RDepotConfig) ([]byte, error){
    "packages": client.PackagesList,
  }
)

func init() {
  rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
  Use:   "list [resource]",
  Short: "list one or many resources",
  Long:  `List one or many resources`,
  Args: cobra.ExactValidArgs(1),
  ValidArgs: []string{"packages"},
  Run: func(cmd *cobra.Command, args []string) {
    cfg := client.RDepotConfig {
      Host: Host,
      Token: Token,
    }

    fmt.Println(cfg)
    res, err := client.PackagesList(cfg)
    switch {
    case err != nil:
      fmt.Println("Error!")
      fmt.Println(err)
    default:
      fmt.Printf(string(res))
    }
  },
}

