package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of pivotal",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pivotal-tools -- HEAD")
	},
}
