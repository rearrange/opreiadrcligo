package cmd

import (
	"fmt"
	"strings"

	"github.com/rearrange/opreiadrcligo/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Architecture Decision Records",
	Long: `Lists all ADRs in docs/adr/ as a formatted table.

The Status column reflects the current value inside each file, so a decision
that has been updated from "Draft" to "Accepted" is shown correctly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		adrs, err := core.List()
		if err != nil {
			return err
		}

		if len(adrs) == 0 {
			fmt.Println("No ADRs found. Run 'adr new \"Your decision title\"' to create one.")
			return nil
		}

		printTable(adrs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// printTable renders adrs as a padded Markdown-style table whose column widths
// are derived from the longest value in each column.
func printTable(adrs []core.ADREntry) {
	const (
		hNum    = "#"
		hTitle  = "Title"
		hDate   = "Date"
		hStatus = "Status"
	)

	// Minimum widths are the header label lengths.
	wNum, wTitle, wDate, wStatus := len(hNum), len(hTitle), len(hDate), len(hStatus)
	for _, e := range adrs {
		if n := len(fmt.Sprintf("%04d", e.Number)); n > wNum {
			wNum = n
		}
		if n := len(e.Title); n > wTitle {
			wTitle = n
		}
		if n := len(e.Date); n > wDate {
			wDate = n
		}
		if n := len(e.Status); n > wStatus {
			wStatus = n
		}
	}

	row := func(num, title, date, status string) {
		fmt.Printf("| %-*s | %-*s | %-*s | %-*s |\n",
			wNum, num, wTitle, title, wDate, date, wStatus, status)
	}
	sep := func() {
		fmt.Printf("|-%s-|-%s-|-%s-|-%s-|\n",
			strings.Repeat("-", wNum),
			strings.Repeat("-", wTitle),
			strings.Repeat("-", wDate),
			strings.Repeat("-", wStatus))
	}

	row(hNum, hTitle, hDate, hStatus)
	sep()
	for _, e := range adrs {
		row(fmt.Sprintf("%04d", e.Number), e.Title, e.Date, e.Status)
	}
}
