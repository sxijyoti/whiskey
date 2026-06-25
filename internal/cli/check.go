package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/config"
	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/planner"
)

var checkCmd = &cobra.Command{
	Use:   "check [site-root]",
	Short: "Inspect site dependencies",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(
		cmd *cobra.Command,
		args []string,
	) error {

		root := "site"

		if len(args) == 1 {
			root = args[0]
		}

		cfg, err := config.Load(root)
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

		store, err := fingerprint.Load(
			".whiskey/fingerprints.json",
		)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Dependency Graph")
		fmt.Println()

		for _, edge := range g.Edges {

			fmt.Printf(
				"%s\n",
				edge.From,
			)

			fmt.Printf(
				" └── %s\n\n",
				edge.To,
			)
		}

		changed, err := fingerprint.ChangedSources(
			g,
			store,
		)
		if err != nil {
			return err
		}

		fmt.Println("Changed Sources")
		fmt.Println()

		if len(changed) == 0 {

			fmt.Println("(none)")
			fmt.Println()

		} else {

			for _, src := range changed {
				fmt.Println(src)
			}

			fmt.Println()
		}

		localDirty, err := planner.LocalDirtyPages(
			filepath.Join(
				root,
				"content",
			),
			store,
		)
		if err != nil {
			return err
		}

		dirty := planner.IncrementalDirtySet(
			g,
			localDirty,
			changed,
		)

		fmt.Println("Dirty Pages")
		fmt.Println()

		if len(dirty) == 0 {

			fmt.Println("(none)")
			fmt.Println()

		} else {

			for _, page := range dirty {
				fmt.Println(page)
			}

			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(
		checkCmd,
	)
}
