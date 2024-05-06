package qrs

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type RequestPrintMsg struct {
	Err error
}

func RequestPrint(code string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("lpr", QrFilePath(code))
		if err := cmd.Run(); err != nil {
			return RequestPrintMsg{Err: err}
		}
		return RequestPrintMsg{}
	}
}
