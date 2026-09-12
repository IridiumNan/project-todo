package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/IridiumNan/project-todo/internal/models"
)

func OpenWithEnvEditor(filePath string, defaultEditor string) error {
	editor := os.Getenv("EDITOR")
	if editor == models.EmptyStr {
		slog.Warn("env variable EDITOR NOT SET, use default editor", "default_editor", defaultEditor)
		editor = defaultEditor
	}
	return OpenWithProgram(editor, filePath)
}

func OpenWithProgram(program string, filePath string) error {
	if program == models.EmptyStr {
		return fmt.Errorf("empty cmd, fail to open file")
	}

	cmd := exec.Command(program, filePath)
	slog.Debug("exec: open with editor", "cmd", cmd.String())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error when exec command, cmd string: %s, err: %s", cmd.String(), err.Error())
	}

	return nil
}
