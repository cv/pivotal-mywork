package main

import (
	"os"
	"github.com/cv/pivotal-tools/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
