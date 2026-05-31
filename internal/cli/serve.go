package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local development server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Serving the static site on port %d...", port)
		return nil
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
	rootCmd.AddCommand(serveCmd)
}