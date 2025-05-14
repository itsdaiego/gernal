package ui

import "github.com/charmbracelet/bubbles/table"

func FormatRows(rows []table.Row, columns []table.Column) []table.Row {
	formattedRows := make([]table.Row, len(rows))

	for i, row := range rows {
		formattedRow := make(table.Row, len(row))
		for j, cell := range row {
			if j < len(columns) {
				formattedRow[j] = truncateString(cell, columns[j].Width)
			} else {
				formattedRow[j] = cell
			}
		}
		formattedRows[i] = formattedRow
	}

	return formattedRows
}

func truncateString(s string, width int) string {
	if width <= 3 {
		return s[:width]
	}

	if len(s) <= width {
		return s
	}

	return s[:width-3] + "..."
}
