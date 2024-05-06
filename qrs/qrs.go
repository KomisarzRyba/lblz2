package qrs

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateQrMsg struct {
	Err  error
	Code string
}

func CreateQr(recordId, instrumentId string) tea.Cmd {
	return func() tea.Msg {
		code := strings.Join([]string{recordId, strings.ReplaceAll(instrumentId, " ", "_")}, "_")
		qr, err := newQr(code)
		if err != nil {
			return CreateQrMsg{Err: err}
		}
		file, err := os.Create(filepath.Join(os.TempDir(), code+".png"))
		if err != nil {
			return CreateQrMsg{Err: err}
		}
		defer file.Close()
		if err := png.Encode(file, qr); err != nil {
			return CreateQrMsg{Err: err}
		}
		return CreateQrMsg{Err: nil, Code: code}
	}
}

func newQr(from string) (barcode.Barcode, error) {
	qr, err := qr.Encode(from, qr.H, qr.Auto)
	if err != nil {
		return nil, err
	}
	qr, err = barcode.Scale(qr, 200, 200)
	if err != nil {
		return nil, err
	}
	return qr, nil
}
