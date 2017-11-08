package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var RootCmd = &cobra.Command{
	Use:   "pivotal",
	Short: "pivotal is a simple command line for Pivotal Tracker",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Help())
	},
}
