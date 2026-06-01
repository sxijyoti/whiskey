package cli

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/sxijyoti/whiskey/internal/build"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve [site-root]",
	Short: "Build and serve a Whiskey site",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		root := "site"

		if len(args) == 1 {
			root = args[0]
		}

		if err := build.BuildSite(root); err != nil {
			return err
		}

		fmt.Printf(
			"Serving %s at http://localhost:%d\n",
			root,
			port,
		)

		fs := http.FileServer(http.Dir("dist"))

		return http.ListenAndServe(
			fmt.Sprintf(":%d", port),
			fs,
		)
	},
}

func init() {
	serveCmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		8080,
		"server port",
	)

	rootCmd.AddCommand(serveCmd)
}
