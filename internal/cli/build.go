package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
)

var buildCmd = &cobra.Command{
	Use:   "build <input.md> <output.html>",
	Short: "Build a markdown file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return build.BuildPage(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
