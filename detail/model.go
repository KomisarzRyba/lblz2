package detail

import (
	"fmt"
	"strings"

	"github.com/KomisarzRyba/lblz2/db"
	"github.com/KomisarzRyba/lblz2/keymap"
	"github.com/KomisarzRyba/lblz2/qrs"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taivy/go.btcqr"
)

type Model struct {
	record   db.Record
	keys     keymap.DetailKeyMap
	help     help.Model
	airtable *db.Airtable
}

func NewModel(record db.Record, airtable *db.Airtable) *Model {
	model := &Model{
		record:   record,
		keys:     keymap.NewDetailKeymap(),
		help:     help.New(),
		airtable: airtable,
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
			return m, qrs.CreateQr(m.record.ID, m.record.Fields.ID)
		}
	case qrs.CreateQrMsg:
		if msg.Err != nil {
			return m, tea.Println(msg.Err)
		}
		return m, tea.Batch(
			tea.Println("QR generated successfully, updating records"),
			m.airtable.UpdateBarcodeField(m.record.ID, msg.Code),
		)
	case db.UpdateBarcodeFieldMsg:
		if msg.Err != nil {
			return m, tea.Println(msg.Err)
		}
		return m, tea.Batch(
			tea.Println("Record successfully updated, printing"),
			qrs.RequestPrint(msg.UpdatedBarcode),
		)
	case qrs.RequestPrintMsg:
		if msg.Err != nil {
			return m, tea.Println(msg.Err)
		}
		return m, func() tea.Msg { return DetailCloseMsg{} }
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
	s.WriteString("\n\n")
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
