package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean generated artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := siteRoot(args)

		_ = os.RemoveAll(
			filepath.Join(root, "dist"),
		)

		_ = os.RemoveAll(
			filepath.Join(root, ".whiskey"),
		)

		fmt.Println("Clean complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}