package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	RunE: func(cmd *cobra.Command, args []string) error {
		return build.BuildSite()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
