package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "site"

		if len(args) == 1 {
			root = args[0]
		}

		return build.BuildSite(root)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
