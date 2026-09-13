package viewer

import (
	"fmt"
	"log/slog"

	"github.com/IridiumNan/project-todo/internal/utils"
)

// ContextViewer for view context file
// All context is store as a single file
// And the ContextViewer is used to display the context with different way
// It can be open with $EDITOR, or formatted print, colorful print or other method
type ContextViewer interface {
	// PathView receive the stored context file path
	// Then display it by specific method
	// This function will open the context file read-only
	PathView(ctxPath string) error

	// PathView open context file directly
	// User can edit this context file
	PathEdit(ctxPath string) error

	// StrView receive the content of whole context file
	// Usually formats or colors it for display
	// StrView(content string)
}

// EditorViewer which using terminal editor to view the context file
// It doesn't handle docx, doc or png files
type EditorViewer struct {
	// ViewCommands should be a shell command which receive a file path then display it's markdown content
	// For instance, nvim, marktext, helix ...
	// It support customization on toml configuration file
	ViewCommands []string
}

// PathView for EditorViewer
// use the editor to open context file
// It will try the ViewCommands one by one
func (ev *EditorViewer) PathView(ctxPath string) error {
	return ev.pathOpen(ctxPath, utils.ModeRead)
}

func (ev *EditorViewer) PathEdit(ctxPath string) error {
	return ev.pathOpen(ctxPath, utils.ModeEdit)
}

func (ev *EditorViewer) pathOpen(ctxPath string, mode utils.OpenMode) error {
	var err error
	for _, viewCmd := range ev.ViewCommands {
		switch mode {
		case utils.ModeEdit:
			err = utils.OpenWithProgram(viewCmd, ctxPath, utils.NoFlag)
		case utils.ModeRead:
			err = utils.OpenWithProgram(viewCmd, ctxPath, utils.GetReadOnlyFlag(viewCmd))
		}
		// if err == nil, view complete then just return
		if err == nil {
			return nil
		}
		slog.Warn("while viewing context", "err", err)
	}

	return fmt.Errorf("error when open with editor viewer: %s", err.Error())
}

func NewEditorViewer(commands []string) *EditorViewer {
	return &EditorViewer{
		ViewCommands: commands,
	}
}
