package main

import (
	"fmt"
	"net/http"
    //"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	status   int
	err      error
}

func initialModel() model {
	return model{
		choices:  []string{"item 1", "item 2", "item 3"},
		selected: make(map[int]struct{}),
		status:   0,
		err:      nil,
	}
}

var url = "https://google.de/"

func (m model) Init() tea.Cmd {
	return checkSomeUrl(url)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := "Todo List\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += fmt.Sprintf("Checking %s ...", url)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	s += "\nPress q to quit.\n"

	return s
}

type statusMsg int
type errMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	case statusMsg:
		m.status = int(msg)
		return m, nil
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	return m, nil
}

func checkSomeUrl(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)
		if err != nil {
			return errMsg{err}
		}
		return statusMsg(res.StatusCode)
	}
}

// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }
