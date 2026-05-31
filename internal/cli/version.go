package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information about whiskey",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Whiskey %s\n", Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}