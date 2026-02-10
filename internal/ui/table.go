package ui

import (
	"fmt"
	"tz-snipe/internal/core"

	"github.com/charmbracelet/bubbles/table"
)

func newTable(preds []core.Prediction) table.Model {
	cols := []table.Column{
		{Title: "Country", Width: 20},
		{Title: "Probability", Width: 30},
	}

	var rows []table.Row

	// format data & build rows
	for i := 0; i < len(preds); i++ {
		p := preds[i]
		rows = append(rows, table.Row{p.Country, fmt.Sprintf("%.4f", p.Probability*100)})
	}

	table := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10), // scrollable after 10
	)

	return table
}
