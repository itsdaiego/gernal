package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	tableView page = iota
	detailView
)

type model struct {
	Table       table.Model
	page        page
	selectedRow int
	minColWidth int
}

type Coin struct {
	Name   string
	Symbol string
	Price  float64
}

func (m model) Init() tea.Cmd {
	return nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var detailStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(1, 2).
	Width(50)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "ctrl+d":
			return m, tea.Quit
		case "left", "right", "up", "down":
			if m.page == tableView {
				m.Table, cmd = m.Table.Update(msg)
				return m, cmd
			}
		case "enter":
			if m.page == tableView {
				m.page = detailView
				m.selectedRow = m.Table.Cursor()
				return m, nil
			}
		case "esc":
			if m.page == detailView {
				m.page = tableView
				return m, nil
			}
		case "shift+up":
			m.Table.SetHeight(m.Table.Height() + 1)
			return m, nil
		case "shift+down":
			if m.Table.Height() > 3 {
				m.Table.SetHeight(m.Table.Height() - 1)
			}
			return m, nil
		case "shift+left":
			// Resize columns proportionally
			cols := m.Table.Columns()
			for i := range cols {
				if cols[i].Width > m.minColWidth {
					cols[i].Width--
					m.Table.SetColumns(cols)
					break
				}
			}
			return m, nil
		case "shift+right":
			// Increase width of smallest column first
			cols := m.Table.Columns()
			minIdx := 0
			for i := 1; i < len(cols); i++ {
				if cols[i].Width < cols[minIdx].Width {
					minIdx = i
				}
			}
			cols[minIdx].Width++
			m.Table.SetColumns(cols)
			return m, nil
		case "1", "2", "3":
			// Resize specific columns
			colIdx := int(msg.String()[0] - '1')
			cols := m.Table.Columns()
			if colIdx >= 0 && colIdx < len(cols) {
				cols[colIdx].Width++
				m.Table.SetColumns(cols)
			}
			return m, nil
		case "!", "@", "#":
			// Decrease specific columns
			keyMap := map[string]int{"!": 0, "@": 1, "#": 2}
			colIdx := keyMap[msg.String()]
			cols := m.Table.Columns()
			if colIdx >= 0 && colIdx < len(cols) && cols[colIdx].Width > m.minColWidth {
				cols[colIdx].Width--
				m.Table.SetColumns(cols)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.page == detailView {
		rows := m.Table.Rows()

		if m.selectedRow >= 0 && m.selectedRow < len(rows) {
			selectedCoin := rows[m.selectedRow]
			detail := fmt.Sprintf("Detailed Information\n\nName: %s\nSymbol: %s\nPrice: $%s\n\nPress ESC to go back",
				selectedCoin[0], selectedCoin[1], selectedCoin[2])
			return detailStyle.Render(detail)
		}
	}

	helpText := "\nResize table: shift+arrow keys | Resize columns: 1,2,3 to increase, !,@,# to decrease"
	helpText += "\nSelect: arrow keys | View details: enter | Quit: q or ctrl+c"
	return baseStyle.Render(m.Table.View()) + helpText + "\n"
}

// truncateString ensures text fits within column width with ellipsis
func truncateString(s string, width int) string {
	if width <= 3 {
		return s[:width]
	}
	
	if len(s) <= width {
		return s
	}
	
	return s[:width-3] + "..."
}

// formatRows ensures all data fits within column widths
func formatRows(rows []table.Row, columns []table.Column) []table.Row {
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

func main() {
	minColWidth := 3
	
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Symbol", Width: 6},
		{Title: "Price", Width: 8},
	}

	rows := []table.Row{
		{"Bitcoin", "BTC", fmt.Sprintf("%.2f", 60000.0)},
		{"Ethereum", "ETH", fmt.Sprintf("%.2f", 4000.0)},
		{"Litecoin", "LTC", fmt.Sprintf("%.2f", 200.0)},
		{"Cardano", "ADA", fmt.Sprintf("%.2f", 1.20)},
		{"Polkadot", "DOT", fmt.Sprintf("%.2f", 30.50)},
	}

	// Format rows to fit column widths
	formattedRows := formatRows(rows, columns)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(formattedRows),
		table.WithFocused(true),
	)
	t.SetHeight(10)
	t.SetStyles(s)

	m := model{
		Table:       t,
		page:        tableView,
		minColWidth: minColWidth,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
