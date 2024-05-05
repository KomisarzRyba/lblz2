package main

import (
	"fmt"
	"os"

	"github.com/KomisarzRyba/lblz2/db"
	"github.com/KomisarzRyba/lblz2/keymap"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	airtable   *db.Airtable
	isSelected bool
	table      table.Model
	tableKeys  keymap.TableKeyMap
	help       help.Model
}

func newModel() model {
	airtable, err := db.NewAirtableFromEnv()
	if err != nil {
		fmt.Printf("Oopsie, there's been an error: %v\n", err)
		os.Exit(1)
	}
	return model{
		airtable:   airtable,
		isSelected: false,
		table: table.New(table.WithColumns([]table.Column{
			{Title: "Type", Width: 12},
			{Title: "Brand", Width: 12},
			{Title: "Model", Width: 12},
			{Title: "Location", Width: 24},
		}), table.WithFocused(true)),
		tableKeys: keymap.NewTableKeyMap(),
		help:      help.New(),
	}
}

func (m model) Init() tea.Cmd { return m.airtable.FetchInstruments() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if !m.isSelected {
			switch {
			case key.Matches(msg, m.tableKeys.Select):
				return m, tea.Println(m.table.SelectedRow())
			case key.Matches(msg, m.tableKeys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.tableKeys.Help):
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	case db.PaginatedInstrumentsMsg:
		if err := msg.Err; err != nil {
			return m, tea.Println(msg.Err)
		}
		newRows := make([]table.Row, len(msg.Instruments))
		for i, instr := range msg.Instruments {
			newRows[i] = instr.Row()
		}
		m.table.SetRows(append(m.table.Rows(), newRows...))
		if msg.Offset != "" {
			return m, m.airtable.FetchPaginatedInstruments(msg.Offset)
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s", m.help.View(m.tableKeys), m.table.View())
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Oopsie, there's been an error: %v", err)
		os.Exit(1)
	}
}
