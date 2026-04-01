/*
Copyright © 2026 Sallehin Sallehuddin <rearrange@users.noreply.github.com>
*/
package cmd

import (
	"fmt"

	"github.com/rearrange/opreiadrcligo/internal/core"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the ADR directory and index",
	Long: `Creates the docs/adr/ directory and an index file (README.md)
to track all Architecture Decision Records.

This command must be run before creating any ADRs. It is safe to run only once —
it will exit with an error if the directory or index already exists.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := core.Init()
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created directory  : %s\n", result.Dir)
		fmt.Printf("✓ Created index file : %s\n", result.IndexFile)
		fmt.Println("\nYou can now create your first ADR:")
		fmt.Println(`  adr new "Record your first decision"`)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
