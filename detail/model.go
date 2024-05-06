package detail

import (
	"fmt"
	"strings"

	"github.com/KomisarzRyba/lblz2/db"
	"github.com/KomisarzRyba/lblz2/keymap"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	record db.Record
	keys   keymap.DetailKeyMap
	help   help.Model
}

func NewModel(record db.Record) *Model {
	return &Model{
		record: record,
		keys:   keymap.NewDetailKeymap(),
		help:   help.New(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

type DetailCloseMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return DetailCloseMsg{} }
		case key.Matches(msg, m.keys.Print):
			return m, tea.Println("printing " + m.record.ID)
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.help.View(m.keys))
	s.WriteString("\n")
	s.WriteString(
		lipgloss.PlaceHorizontal(80, lipgloss.Center,
			fmt.Sprintf("%s %s %s", m.record.Fields.Brand, m.record.Fields.Model, m.record.Fields.Type),
		),
	)

	return s.String()
}