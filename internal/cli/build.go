package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
)

var full bool

var buildCmd = &cobra.Command{
	Use:   "build [site-root]",
	Short: "Build the site",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		root := "site"

		if len(args) == 1 {
			root = args[0]
		}

		if full {
			return build.BuildSite(
				root,
			)
		}

		return build.IncrementalBuild(
			root,
		)
	},
}

func init() {

	buildCmd.Flags().BoolVar(
		&full,
		"full",
		false,
		"force full rebuild",
	)

	rootCmd.AddCommand(
		buildCmd,
	)
}