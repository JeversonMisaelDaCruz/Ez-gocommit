package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeversonmisael/ez-gocommit/internal/ai"
)

type Result struct {
	Message   string
	Body      string
	Cancelled bool
}

var (
	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63"))

	styleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	styleUnselected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	styleHighConf = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true)

	styleMedConf = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	styleLowConf = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	styleReasoning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)

	styleEditLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true)

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

type mode int

const (
	modeSelect mode = iota
	modeEdit
)

type model struct {
	suggestions []ai.Suggestion
	cursor      int
	mode        mode
	editBuffer  string
	editCursor  int
	result      *Result
}

func newModel(suggestions []ai.Suggestion) model {
	return model{
		suggestions: suggestions,
		cursor:      0,
		mode:        modeSelect,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeSelect:
			return m.updateSelect(msg)
		case modeEdit:
			return m.updateEdit(msg)
		}
	}
	return m, nil
}

func (m model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.suggestions)-1 {
			m.cursor++
		}
	case "1":
		if len(m.suggestions) >= 1 {
			m.cursor = 0
		}
	case "2":
		if len(m.suggestions) >= 2 {
			m.cursor = 1
		}
	case "3":
		if len(m.suggestions) >= 3 {
			m.cursor = 2
		}
	case "e":
		m.mode = modeEdit
		m.editBuffer = m.suggestions[m.cursor].Message
		m.editCursor = len(m.editBuffer)
	case "enter":
		selected := m.suggestions[m.cursor]
		m.result = &Result{
			Message: selected.Message,
			Body:    selected.Body,
		}
		return m, tea.Quit
	case "q", "ctrl+c", "esc":
		m.result = &Result{Cancelled: true}
		return m, tea.Quit
	}
	return m, nil
}

func (m model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.editBuffer) == "" {
			return m, nil
		}
		m.result = &Result{
			Message: strings.TrimSpace(m.editBuffer),
			Body:    m.suggestions[m.cursor].Body,
		}
		return m, tea.Quit
	case "esc":
		m.mode = modeSelect
		m.editBuffer = ""
	case "ctrl+c":
		m.result = &Result{Cancelled: true}
		return m, tea.Quit
	case "backspace":
		if m.editCursor > 0 {
			m.editBuffer = m.editBuffer[:m.editCursor-1] + m.editBuffer[m.editCursor:]
			m.editCursor--
		}
	case "left":
		if m.editCursor > 0 {
			m.editCursor--
		}
	case "right":
		if m.editCursor < len(m.editBuffer) {
			m.editCursor++
		}
	case "ctrl+a", "home":
		m.editCursor = 0
	case "ctrl+e", "end":
		m.editCursor = len(m.editBuffer)
	default:
		if len(msg.Runes) > 0 {
			ch := string(msg.Runes)
			m.editBuffer = m.editBuffer[:m.editCursor] + ch + m.editBuffer[m.editCursor:]
			m.editCursor += len(ch)
		}
	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("  Ez-gocommit — Select a commit message") + "\n\n")

	for i, s := range m.suggestions {
		isSelected := i == m.cursor

		var confBadge string
		switch strings.ToLower(s.Confidence) {
		case "high":
			confBadge = styleHighConf.Render("●● HIGH  ")
		case "medium":
			confBadge = styleMedConf.Render("●○ MED   ")
		default:
			confBadge = styleLowConf.Render("○○ LOW   ")
		}

		rank := fmt.Sprintf("[%d]", s.Rank)

		if isSelected {
			cursor := styleSelected.Render("▶")
			msgStr := styleSelected.Render(s.Message)
			line := fmt.Sprintf(" %s %s %s  %s", cursor, rank, confBadge, msgStr)
			sb.WriteString(line + "\n")
		} else {
			cursor := "  "
			msgStr := styleUnselected.Render(s.Message)
			line := fmt.Sprintf(" %s %s %s  %s", cursor, rank, confBadge, msgStr)
			sb.WriteString(styleUnselected.Render(line) + "\n")
		}
	}

	if m.cursor < len(m.suggestions) {
		reasoning := m.suggestions[m.cursor].Reasoning
		if reasoning != "" {
			sb.WriteString("\n")
			sb.WriteString(styleReasoning.Render("  💬 "+reasoning) + "\n")
		}
		body := m.suggestions[m.cursor].Body
		if body != "" {
			sb.WriteString(styleReasoning.Render("  📝 Body: "+truncateStr(body, 80)) + "\n")
		}
	}

	sb.WriteString("\n")

	if m.mode == modeEdit {
		sb.WriteString(styleEditLabel.Render("  Edit message:") + "\n")
		before := m.editBuffer[:m.editCursor]
		after := m.editBuffer[m.editCursor:]
		editLine := "  " + before + "│" + after
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(editLine) + "\n\n")
		sb.WriteString(styleHelp.Render("  Enter confirm • Esc cancel edit • Ctrl+C abort") + "\n")
	} else {
		sb.WriteString(styleHelp.Render("  ↑↓/jk navigate • 1-3 jump • Enter confirm • e edit • q abort") + "\n")
	}

	return styleBorder.Render(sb.String())
}

func Run(suggestions []ai.Suggestion) (*Result, error) {
	m := newModel(suggestions)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("TUI error: %w", err)
	}
	fm := finalModel.(model)
	if fm.result == nil {
		return &Result{Cancelled: true}, nil
	}
	return fm.result, nil
}

func truncateStr(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
