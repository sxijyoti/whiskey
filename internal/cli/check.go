package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/graph"
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

		fmt.Println()
		fmt.Println(
			"Dependency Graph",
		)
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(
		checkCmd,
	)
}