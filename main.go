package main

import (
	"fmt"
	"os"

	"github.com/KomisarzRyba/lblz2/db"
	"github.com/KomisarzRyba/lblz2/detail"
	"github.com/KomisarzRyba/lblz2/keymap"
	"github.com/KomisarzRyba/lblz2/qrs"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type model struct {
	airtable  *db.Airtable
	table     table.Model
	tableKeys keymap.TableKeyMap
	help      help.Model
	detail    *detail.Model
}

func newModel() model {
	airtable, err := db.NewAirtableFromEnv()
	if err != nil {
		fmt.Printf("Oopsie, there's been an error: %v\n", err)
		os.Exit(1)
	}
	return model{
		airtable: airtable,
		table: table.New([]table.Column{
			table.NewFlexColumn("type", "Type", 1).WithFiltered(true).WithStyle(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#f4dbd6")),
			),
			table.NewColumn("brand", "Brand", 10).WithFiltered(true).WithStyle(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#f0c6c6")),
			),
			table.NewFlexColumn("model", "Model", 2).WithFiltered(true).WithStyle(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#f5bde6")),
			),
			table.NewColumn("location", "Location", 18).WithFiltered(true).WithStyle(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#f5a97f")),
			),
			table.NewColumn("has_qr", "Label", 5).WithStyle(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#a6da95")),
			),
		}).Filtered(true).Focused(true).WithPageSize(12).WithTargetWidth(80).WithMissingDataIndicatorStyled(
			table.NewStyledCell(
				"x", lipgloss.NewStyle().Foreground(
					lipgloss.Color("#ed8796"),
				),
			),
		).WithBaseStyle(
			lipgloss.NewStyle().BorderForeground(lipgloss.Color("#b7bdf8")),
		),
		tableKeys: keymap.NewTableKeyMap(),
		help:      help.New(),
		detail:    nil,
	}
}

func (m model) Init() tea.Cmd { return m.airtable.FetchInstruments() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.detail != nil {
		switch msg.(type) {
		case detail.DetailCloseMsg:
			m.detail = nil
			return m, nil
		}
		newDetail, cmd := m.detail.Update(msg)
		if d, ok := newDetail.(detail.Model); ok {
			m.detail = &d
		}
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tableKeys.Select):
			m.detail = detail.NewModel(db.RecordFromRow(m.table.HighlightedRow().Data), m.airtable)
		case key.Matches(msg, m.tableKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.tableKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	case db.PaginatedInstrumentsMsg:
		if msg.Err != nil {
			return m, tea.Println(msg.Err)
		}
		newRows := make([]table.Row, len(msg.Records))
		for i, record := range msg.Records {
			newRows[i] = record.Row()
		}
		m.table = m.table.WithRows(append(m.table.GetVisibleRows(), newRows...))
		if msg.Offset != "" {
			return m, m.airtable.FetchPaginatedInstruments(msg.Offset)
		}
	case qrs.CreateQrMsg:
		if msg.Err != nil {
			return m, tea.Println(msg.Err)
		}
		return m, tea.Println(msg.Code)
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.detail != nil {
		return m.detail.View()
	}
	return fmt.Sprintf("%s\n%s", m.help.View(m.tableKeys), m.table.View())
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oopsie, there's been an error: %v", err)
		os.Exit(1)
	}
}
