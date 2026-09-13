package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/IridiumNan/project-todo/internal/models"
)

type OpenMode int

const (
	ModeRead = iota
	ModeEdit
)

const NoFlag = ""

func OpenWithEnvEditor(filePath string, defaultEditor string, mode OpenMode) error {
	editor := os.Getenv("EDITOR")
	if editor == models.EmptyStr {
		slog.Warn("env variable EDITOR NOT SET, use default editor", "default_editor", defaultEditor)
		editor = defaultEditor
	}
	switch mode {
	case ModeRead:
		return OpenWithProgram(editor, filePath, GetReadOnlyFlag(editor))
	default:
		return OpenWithProgram(editor, filePath, NoFlag)
	}
}

// OpenWithProgram if you want to get the read-only flag, use function [GetReadOnlyFlag]
// If don't need any flag, use [NoFlag] as flag
func OpenWithProgram(program string, filePath string, flag string) error {
	if program == models.EmptyStr {
		return fmt.Errorf("empty cmd, fail to open file")
	}

	var cmd *exec.Cmd
	if flag == NoFlag {
		cmd = exec.Command(program, filePath)
	} else {
		cmd = exec.Command(program, flag, filePath)
	}
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

// GetReadOnlyFlag return the editor's read only flag
func GetReadOnlyFlag(editor string) (flag string) {
	switch editor {
	case "nano":
		return "-v"
	case "vim", "nvim":
		return "-R"
	}

	return NoFlag
}

// TODO: add the read-only mode and edit mode for different editor
// vim and nvim with -R flag (read-only mode)
// nano with -v flag
