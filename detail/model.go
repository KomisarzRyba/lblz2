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
	"github.com/taivy/go.btcqr"
)

type Model struct {
	record db.Record
	keys   keymap.DetailKeyMap
	help   help.Model
}

func NewModel(record db.Record) *Model {
	model := &Model{
		record: record,
		keys:   keymap.NewDetailKeymap(),
		help:   help.New(),
	}
	if record.Fields.Barcode.Text != "" {
		model.keys = model.keys.WithReprint()
	}
	return model
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
	s.WriteString(
		lipgloss.PlaceHorizontal(80, lipgloss.Center,
			lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Render(
				fmt.Sprintf("%s %s %s\n%s\nCount: %d | Color: %v",
					m.record.Fields.Brand,
					m.record.Fields.Model,
					m.record.Fields.Type,
					m.record.Fields.Location,
					m.record.Fields.Count,
					m.record.Fields.Color,
				),
			),
		),
	)
	s.WriteString("\n")
	if code := m.record.Fields.Barcode.Text; code != "" {
		qr, err := qrt.Generate(code)
		if err == nil {
			s.WriteString(
				lipgloss.PlaceHorizontal(80, lipgloss.Center, qr),
			)
		}
	} else {
		s.WriteString(
			lipgloss.Place(80, 8, lipgloss.Center, lipgloss.Center,
				fmt.Sprintf("%s is not labeled", m.record.Fields.Type),
			),
		)
	}
	s.WriteString("\n")
	s.WriteString(
		lipgloss.PlaceHorizontal(80, lipgloss.Center,
			m.help.View(m.keys),
		),
	)

	return s.String()
}
