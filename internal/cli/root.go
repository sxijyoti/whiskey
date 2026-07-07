package cli

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "whiskey",
	Short: "static site generator",
	Long: "Whiskey is a static site generator written in Go.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func siteRoot(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return "site"
}