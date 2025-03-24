package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

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
	return "Hello, World! Press q to quit."
}

func main() {
	fmt.Println("Hello, World!")

	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return
	}
}
