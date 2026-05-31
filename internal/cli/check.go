package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate site configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking site...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}