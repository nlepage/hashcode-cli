package main

import (
	"os"

	"github.com/nlepage/hashcode-cli/cmd"

	_ "github.com/nlepage/hashcode-cli/upload"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
