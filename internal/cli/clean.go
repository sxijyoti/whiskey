package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean generated artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = os.RemoveAll("dist")
		_ = os.RemoveAll(".whiskey")

		fmt.Println("Clean complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}