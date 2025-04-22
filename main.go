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
			m.Table.SetHeight(m.Table.Height() - 1)
			return m, nil
		case "shift+left":
			m.Table.SetWidth(m.Table.Width() - 1)
			return m, nil
		case "shift+right":
			m.Table.SetWidth(m.Table.Width() + 1)
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

	helpText := "\nResize: shift+arrow keys | Select: arrow keys | View details: enter | Quit: q or ctrl+c"
	return baseStyle.Render(m.Table.View()) + helpText + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Symbol", Width: 3},
		{Title: "Price", Width: 5},
	}

	rows := []table.Row{
		{"Bitcoin", "BTC", fmt.Sprintf("%.2f", 60000.0)},
		{"Ethereum", "ETH", fmt.Sprintf("%.2f", 4000.0)},
		{"Litecoin", "LTC", fmt.Sprintf("%.2f", 200.0)},
	}

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
		table.WithRows(rows),
		table.WithFocused(true),
	)
	t.SetHeight(10)
	t.SetStyles(s)

	m := model{
		Table: t,
		page:  tableView,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
