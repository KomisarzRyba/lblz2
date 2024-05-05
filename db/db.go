package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Airtable struct {
	token   string
	baseId  string
	tableId string
}

type Instrument struct {
	ID       string `json:"ID"`
	Type     string `json:"Equipment Type"`
	Brand    string `json:"Brand"`
	Model    string `json:"Model"`
	Location string `json:"Location"`
}

func (i Instrument) Row() table.Row {
	return table.Row{i.Type, i.Brand, i.Model, i.Location}
}

type instrumentWrapper struct {
	Fields Instrument `json:"fields"`
}

type instrumentResponse struct {
	Records []instrumentWrapper `json:"records"`
	Offset  string              `json:"offset"`
}

func NewAirtableFromEnv() (*Airtable, error) {
	token := os.Getenv("AIRTABLE_TOKEN")
	if token == "" {
		return nil, errors.New("token not set")
	}
	baseId := os.Getenv("AIRTABLE_BASE_ID")
	if baseId == "" {
		return nil, errors.New("base id not set")
	}
	tableId := os.Getenv("AIRTABLE_TABLE_ID")
	if tableId == "" {
		return nil, errors.New("table id not set")
	}
	return &Airtable{token, baseId, tableId}, nil
}

type PaginatedInstrumentsMsg struct {
	Err         error
	Instruments []Instrument
	Offset      string
}

func (a *Airtable) FetchInstruments() tea.Cmd {
	return a.FetchPaginatedInstruments("")
}

func (a *Airtable) FetchPaginatedInstruments(offset string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s", a.baseId, a.tableId)
		if offset != "" {
			url += "?offset=" + offset
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return PaginatedInstrumentsMsg{Err: err}
		}
		req.Header.Set("Authorization", "Bearer "+a.token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return PaginatedInstrumentsMsg{Err: err}
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return PaginatedInstrumentsMsg{Err: err}
		}
		var instrRes instrumentResponse
		err = json.Unmarshal(body, &instrRes)
		if err != nil {
			return PaginatedInstrumentsMsg{Err: err}
		}
		instruments := make([]Instrument, len(instrRes.Records))
		for i, r := range instrRes.Records {
			instruments[i] = r.Fields
		}
		return PaginatedInstrumentsMsg{
			Err:         nil,
			Offset:      instrRes.Offset,
			Instruments: instruments,
		}
	}
}
