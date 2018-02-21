package cmd

import (
	"github.com/nlepage/hashcode-cli/config"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "hashcode",
	Short: "CLI for Hashcode Judge System",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		return config.ReadConfig()
	},
}
