package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Table table.Model
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "ctrl+d":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return baseStyle.Render(m.Table.View()) + "\n"

}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}

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

	m := model{t}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
