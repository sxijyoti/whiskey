package cli

import (
	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/source"
)

var syncCmd = &cobra.Command{
	Use:   "sync [site-root]",
	Short: "Synchronize remote sources",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		root := siteRoot(args)

		if err := build.EnsureWorkspace(
			root,
		); err != nil {
			return err
		}

		cfg, err := config.Load(
			root,
		)
		if err != nil {
			return err
		}

		g, err := graph.BuildSiteGraph(
			root,
			cfg.Theme,
		)
		if err != nil {
			return err
		}

		manifest, err := source.LoadManifest(
			root,
		)
		if err != nil {
			return err
		}

		_, err = build.MaterializeSources(
			root,
			g,
			manifest,
		)

		return err
	},
}

func init() {
	rootCmd.AddCommand(
		syncCmd,
	)
}
