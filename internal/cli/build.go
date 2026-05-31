package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building the site...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}