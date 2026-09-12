package viewer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IridiumNan/project-todo/internal/utils"
)

// PlainPrintViewer a viewer which just print the context on terminal
type PlainPrintViewer struct{}

// PathView for PlainPrintViewer
// Just read file the print it on terminal
func (ppv *PlainPrintViewer) PathView(path string) error {
	ctxByte, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("plain Print Viewer: error when read context file, err: %s", err.Error())
	}

	prefix := "\n\nCONTEXT\n\n"

	fmt.Println(prefix, string(ctxByte))

	return nil
}

// BatPrintViewer use bat command to view context or batcat (on ubuntu)
type BatPrintViewer struct{}

// PathView for [BatPrintViewer]
// Try both bat and batcat command
func (bpv *BatPrintViewer) PathView(path string) error {
	// try bat command first
	// if failed, try batcat (for ubuntu)

	err := utils.OpenWithProgram("bat", path)
	if err != nil {
		slog.Warn("while try to view context with bat command. Try to use batcat", "err", err)

		err = utils.OpenWithProgram("batcat", path)
		if err != nil {
			return fmt.Errorf("failed to view context with bat command, both batcat failed", "err", err)
		}

	}
	return nil
}
