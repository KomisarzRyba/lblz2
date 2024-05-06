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

func QrFilePath(code string) string {
	return filepath.Join(os.TempDir(), code+".png")
}

func CreateQr(recordId, instrumentId string) tea.Cmd {
	return func() tea.Msg {
		code := NewCode(recordId, instrumentId)
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

func NewCode(recordId, instrumentId string) string {
	return strings.Join([]string{recordId, strings.ReplaceAll(instrumentId, " ", "_")}, "_")
}

func newQr(code string) (barcode.Barcode, error) {
	qr, err := qr.Encode(code, qr.H, qr.Auto)
	if err != nil {
		return nil, err
	}
	qr, err = barcode.Scale(qr, 200, 200)
	if err != nil {
		return nil, err
	}
	return qr, nil
}
