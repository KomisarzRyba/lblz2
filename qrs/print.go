package qrs

import (
	"log"
	"os/exec"
)

// TODO: make into tea.Cmd
func CreatePrintJob(filepath string) {
	cmd := exec.Command("lpr", filepath)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing command: %s", err)
	}
}
