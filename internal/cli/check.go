package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/fingerprint"
	"github.com/sxijyoti/whiskey/internal/graph"
	"github.com/sxijyoti/whiskey/internal/source"
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

		g, err := graph.BuildSiteGraph(
			filepath.Join(
				root,
				"content",
			),
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

		fmt.Println("Fingerprints")
		fmt.Println()

		for _, node := range g.Nodes {

			if node.Type != graph.SourceNode {
				continue
			}

			src, err := source.Resolve(
				node.ID,
			)

			if err != nil {
				fmt.Printf(
					"%s error=%v\n",
					node.ID,
					err,
				)
				continue
			}

			hash, err := fingerprint.FingerprintSource(
				src,
			)

			if err != nil {
				fmt.Printf(
					"%s error=%v\n",
					node.ID,
					err,
				)
				continue
			}

			changed := fingerprint.Changed(
				store,
				node.ID,
				hash,
			)

			fmt.Printf(
				"%s\n",
				node.ID,
			)

			fmt.Printf(
				" changed: %v\n\n",
				changed,
			)

			store[node.ID] = hash
		}

		if err := fingerprint.Save(
			".whiskey/fingerprints.json",
			store,
		); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(
		checkCmd,
	)
}
