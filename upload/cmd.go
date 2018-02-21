package upload

import (
	"github.com/nlepage/hashcode-cli/cmd"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads answers to the judge system",
	RunE: func(_ *cobra.Command, args []string) error {
		return upload(args)
	},
}

func init() {
	cmd.Cmd.AddCommand(uploadCmd)
}
