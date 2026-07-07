package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/source"
)

var full bool
var offline bool

var buildCmd = &cobra.Command{
	Use:   "build [site-root]",
	Short: "Build the site",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		root := siteRoot(args)

		source.Offline = offline

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

	buildCmd.Flags().BoolVar(
		&offline,
		"offline",
		false,
		"build using cached remote sources only",
	)

	rootCmd.AddCommand(
		buildCmd,
	)
}