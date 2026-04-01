package cmd

import (
	"fmt"

	"github.com/rearrange/opreiadrcligo/internal/core"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new Architecture Decision Record",
	Long: `Creates a new ADR file in docs/adr/ with a sequential 4-digit number
and the title slugified as the filename, e.g. 0003-use-postgresql.md.

The ADR index (docs/adr/README.md) is updated automatically.
Run 'adr init' before creating your first ADR.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := core.New(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created ADR   : %s\n", result.FilePath)
		fmt.Printf("✓ Updated index : %s\n", core.IndexFile)
		fmt.Printf("\nTitle  : %s\n", result.Title)
		fmt.Printf("Number : %04d\n", result.Number)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
