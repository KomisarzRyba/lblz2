package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type Airtable struct {
	token   string
	baseId  string
	tableId string
}

type Instrument struct {
	ID       string   `json:"ID"`
	Type     string   `json:"Equipment Type"`
	Brand    string   `json:"Brand"`
	Model    string   `json:"Model"`
	Location string   `json:"Location"`
	Color    []string `json:"Color / Material"`
	Count    int      `json:"Count"`
	Barcode  struct {
		Text string `json:"text"`
	} `json:"Barcode"`
}

type Record struct {
	ID     string     `json:"id"`
	Fields Instrument `json:"fields"`
}

func (r Record) Row() table.Row {
	row := table.NewRow(table.RowData{
		"record_id": r.ID,
		"id":        r.Fields.ID,
		"barcode":   r.Fields.Barcode,
		"color":     r.Fields.Color,
		"count":     r.Fields.Count,
		"type":      r.Fields.Type,
		"brand":     r.Fields.Brand,
		"model":     r.Fields.Model,
		"location":  r.Fields.Location,
	})
	if r.Fields.Barcode.Text != "" {
		row.Data["has_qr"] = "✓"
	}
	return row
}

func RecordFromRow(row table.RowData) Record {
	return Record{
		ID: row["record_id"].(string),
		Fields: Instrument{
			ID:       row["id"].(string),
			Type:     row["type"].(string),
			Brand:    row["brand"].(string),
			Model:    row["model"].(string),
			Location: row["location"].(string),
			Color:    row["color"].([]string),
			Count:    row["count"].(int),
			Barcode: struct {
				Text string `json:"text"`
			}{Text: row["barcode"].(struct {
				Text string `json:"text"`
			}).Text},
		},
	}
}

type instrumentResponse struct {
	Records []Record `json:"records"`
	Offset  string   `json:"offset"`
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
	Err     error
	Records []Record
	Offset  string
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
		return PaginatedInstrumentsMsg{
			Err:     nil,
			Offset:  instrRes.Offset,
			Records: instrRes.Records,
		}
	}
}

type UpdateBarcodeFieldMsg struct {
	Err            error
	UpdatedBarcode string
}

func (a *Airtable) UpdateBarcodeField(recordId, code string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s/%s", a.baseId, a.tableId, recordId)
		payload, err := json.Marshal(map[string]interface{}{
			"fields": map[string]interface{}{
				"Barcode": map[string]interface{}{
					"text": code,
				},
			},
		})
		if err != nil {
			return UpdateBarcodeFieldMsg{Err: err}
		}
		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
		req.Header.Set("Authorization", "Bearer "+a.token)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return UpdateBarcodeFieldMsg{Err: err}
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return UpdateBarcodeFieldMsg{Err: err}
		}
		var updatedFieldResponse struct{ Barcode struct{ text string } }
		err = json.Unmarshal(body, &updatedFieldResponse)
		if err != nil {
			return UpdateBarcodeFieldMsg{Err: err}
		}
		return UpdateBarcodeFieldMsg{UpdatedBarcode: updatedFieldResponse.Barcode.text}
	}
}
