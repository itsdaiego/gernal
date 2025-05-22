package main

import (
	"fmt"
	api "main/internal/api"
	ui "main/internal/ui"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	tableView  page = iota
	detailView page = iota
)

type model struct {
	Table           table.Model
	page            page
	selectedRow     int
	minColWidth     int
	minHeight       int
	heightIncrement int
	widthIncrement  int
}

func (m model) Init() tea.Cmd {
	return nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Width(50).
	Align(lipgloss.Center).
	Padding(1, 2)

var detailStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Width(50).
	Padding(1, 2)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "ctrl+d":
			return m, tea.Quit
		case "up", "down":
			if m.page == tableView {
				m.Table, cmd = m.Table.Update(msg)
				return m, cmd
			}
		case "enter", "right":
			if m.page == tableView {
				m.Table.SetHeight(m.Table.Height())
				m.Table.SetWidth(m.Table.Width())
				m.page = detailView
				m.selectedRow = m.Table.Cursor()
				return m, nil
			}
		case "esc", "left":
			if m.page == detailView {
				m.page = tableView
				return m, nil
			}
		case "shift+up":
			if m.Table.Height() > m.minHeight {
				m.Table.SetHeight(m.Table.Height() - m.heightIncrement)
			}
			return m, nil
		case "shift+down":
			m.Table.SetHeight(m.Table.Height() + m.heightIncrement)
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

	return baseStyle.Render(m.Table.View())
}

func main() {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Symbol", Width: 6},
		{Title: "Price", Width: 8},
	}

	rows, err := api.FetchCoins()

	if err != nil {
		fmt.Println("Error fetching coins:", err)
		os.Exit(1)
	}

	t := ui.RenderTable(columns, rows)

	m := model{
		Table:           t,
		page:            tableView,
		minColWidth:     5,
		minHeight:       5,
		heightIncrement: 5,
		widthIncrement:  5,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
