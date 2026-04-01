/*
Copyright © 2026 Sallehin Sallehuddin <rearrange@users.noreply.github.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "adr",
	Short: "ADR CLI - (Opinionated) ADR Tool",
	Long: `ADR CLI is an opinionated tool for managing Architecture Decision Records (ADR).

ADRs are stored as Markdown files under docs/adr/ and tracked via an index file.
The index file will be created as README.md (docs/adr/README.md).
The ADRs will be created following the format "docs/adr/{number}-{title}.md".

Get started by running:
  adr init
  adr new "Record architecture decision"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.adr-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
