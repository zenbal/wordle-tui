package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	wordle      *Wordle
	width       int
	height      int
	inputs      []WordInput
	cursor      int
	help        bool
	hints       bool
	suggestions bool
	hint        string
}

func NewModel() model {
	inputs := make([]WordInput, 6)
	for i := range inputs {
		inputs[i] = NewWordInput()
	}
	return model{
		wordle:      NewWordle(),
		width:       0,
		height:      0,
		inputs:      inputs,
		cursor:      0,
		help:        false,
		hints:       false,
		suggestions: false,
	}
}

type WordInput []textinput.Model

func NewWordInput() WordInput {
	fields := make([]textinput.Model, 5)
	cursor := cursor.New()
	cursor.SetMode(2)
	for i := range fields {
		ti := textinput.New()
		ti.Prompt = " "
		ti.CharLimit = 1
		ti.Cursor = cursor
		ti.Width = 1
		fields[i] = ti
	}
	return fields
}

const (
	colorGreen     = lipgloss.Color("#538d4e")
	colorYellow    = lipgloss.Color("#b59f3b")
	colorGrey      = lipgloss.Color("#3a3a3c")
	colorWhite     = lipgloss.Color("#ffffff")
	colorBlack     = lipgloss.Color("#121213")
	colorLightGrey = lipgloss.Color("#818384")
)

var (
	defaultInputStyle = lipgloss.NewStyle().
				Padding(1, 1).Background(colorBlack).Foreground(colorWhite)
	greenInputStyle  = defaultInputStyle.Copy().Background(colorGreen).Foreground(colorWhite)
	yellowInputStyle = defaultInputStyle.Copy().Background(colorYellow).Foreground(colorBlack)
	greyInputStyle   = defaultInputStyle.Copy().Background(colorGrey).Foreground(colorWhite)
	inputTextStyle   = lipgloss.NewStyle().Transform(strings.ToUpper)
	helpTextStyle    = lipgloss.NewStyle().Foreground(colorLightGrey)
	titleStyle       = lipgloss.NewStyle().PaddingBottom(1).Bold(true)
)

var inputStyle map[int]lipgloss.Style = map[int]lipgloss.Style{
	0: defaultInputStyle,
	1: greyInputStyle,
	2: yellowInputStyle,
	3: greenInputStyle,
}

func (m model) BoardView() string {
	title := "GUESSES"
	if m.wordle.status == WIN {
		title = "YOU WIN"
	} else if m.wordle.status == LOOSE {
		title = "YOU LOOSE"
	}
	rows := make([]string, 7)
	rows = append(rows, titleStyle.Render(title))
	for i := range m.inputs {
		cols := make([]string, 5)
		for j := range m.inputs[i] {
			feedback := TBD
			if m.wordle.board != nil && m.wordle.board[i] != nil {
				feedback = m.wordle.board[i][j].feedback
			}
			cols = append(cols, inputStyle[int(feedback)].Render(inputTextStyle.Render(m.inputs[i][j].View())))
		}
		col := lipgloss.JoinHorizontal(lipgloss.Center, cols...)
		rows = append(rows, col)
	}

	return lipgloss.NewStyle().MarginRight(2).Render(lipgloss.JoinVertical(lipgloss.Center, rows...))
}

func (m model) AsideView() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.AlphabetView(), m.SuggestionView(), m.HintView(), m.HelpView())
}

func (m model) SuggestionView() string {
	if m.suggestions {
		return fmt.Sprintf("try '%s'\n", m.wordle.suggestNextGuess())
	}
	return ""
}

func (m model) HintView() string {
	if m.hints {
		return m.wordle.message
	}
	return ""
}

func (m model) HelpView() string {
	if m.help {
		return helpTextStyle.Render(lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().MarginRight(2).Render(lipgloss.JoinVertical(lipgloss.Left, "?", "C-c", "C-r", "Return", "C-h", "C-s")),
			lipgloss.JoinVertical(lipgloss.Right, "Help", "Quit", "New Game", "Submit Guess", "Show Hints", "Show Suggestions"),
		))
	}
	return helpTextStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().MarginRight(2).Render("?"),
		"Help",
	))
}

func (m model) AlphabetView() string {
	alphabet := [][]rune{
		{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'},
		{'a', 's', 'd', 'f', 'h', 'j', 'k', 'l'},
		{'z', 'x', 'c', 'v', 'b', 'n', 'm'},
	}
	view := make([][]string, 3)
	for i := range view {
		view[i] = make([]string, 10)
	}
	for row := range alphabet {
	outer:
		for col := range alphabet[row] {
			char := alphabet[row][col]
			for _, char_idx := range m.wordle.assign {
				if char_idx == alphabet_idx(byte(char)) {
					view[row][col] = greenInputStyle.Padding(1, 1).Render(strings.ToUpper(string(char)))
					continue outer
				}
			}
			for char_idx, include := range m.wordle.include {
				if char_idx == alphabet_idx(byte(char)) && include {
					view[row][col] = yellowInputStyle.Padding(1, 1).Render(strings.ToUpper(string(char)))
					continue outer
				}
				if char_idx == alphabet_idx(byte(char)) && !include {
					view[row][col] = greyInputStyle.Padding(1, 1).Render(strings.ToUpper(string(char)))
					continue outer
				}
			}
			view[row][col] = defaultInputStyle.Padding(1, 1).Render(strings.ToUpper(string(char)))
		}
	}
	view_joined_rows := make([]string, 3)
	for i := range view_joined_rows {
		view_joined_rows[i] = lipgloss.JoinHorizontal(lipgloss.Bottom, view[i]...)
	}
	return lipgloss.NewStyle().MarginBottom(2).Render(lipgloss.JoinVertical(lipgloss.Center, view_joined_rows...))
}

func (m *model) newGame() {
	m.wordle = NewWordle()
	inputs := make([]WordInput, 6)
	for i := range inputs {
		inputs[i] = NewWordInput()
	}
	m.inputs = inputs
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.BoardView(), m.AsideView()),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			m.newGame()
			return m, cmd
		case tea.KeyBackspace.String():
			m.handleKeyBackspace(msg)
		case tea.KeyEnter.String():
			m.handleKeyEnter()
		case "?":
			m.help = !m.help
		case "ctrl+h":
			m.hints = !m.hints
		case "ctrl+s":
			m.suggestions = !m.suggestions
		default:
            if m.wordle.status != ONGOING {
                m.newGame()
                return m, cmd
            }
			if !(msg.String() >= "a" && msg.String() <= "z") {
				return m, cmd
			}
			m.handleKeyAlphabet(msg)
		}
	}
	return m, cmd
}

func (m *model) handleKeyBackspace(msg tea.KeyMsg) tea.Cmd {
	current_input := &m.inputs[m.wordle.attempt][m.cursor]
	current_input.Focus()
	var cmd tea.Cmd
	if current_input.Value() == "" && m.cursor > 0 {
		m.cursor--
		current_input = &m.inputs[m.wordle.attempt][m.cursor]
	}
	*current_input, cmd = current_input.Update(msg)
	if m.cursor > 0 {
		m.cursor--
	}
	return cmd
}

func (m *model) handleKeyEnter() tea.Cmd {
	var cmd tea.Cmd
	if m.cursor != 4 {
		return cmd
	}
	word := ""
	for i := range m.inputs[m.wordle.attempt] {
		word += m.inputs[m.wordle.attempt][i].Value()
	}

	guess, err := NewGuess(word)
	if err == nil {
		m.wordle.validate(guess)
	}

	if err := m.wordle.guess(strings.ToLower(word)); err != nil {
		return cmd
	}
	m.cursor = 0
	return cmd
}

func (m *model) handleKeyAlphabet(msg tea.KeyMsg) tea.Cmd {
	current_input := &m.inputs[m.wordle.attempt][m.cursor]
	current_input.Focus()
	var cmd tea.Cmd
	if current_input.Value() != "" && m.cursor < 4 {
		m.cursor++
		current_input = &m.inputs[m.wordle.attempt][m.cursor]
	}
	*current_input, cmd = current_input.Update(msg)
	if m.cursor < 4 {
		m.cursor++
	}
	return cmd
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
