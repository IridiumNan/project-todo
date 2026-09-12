package viewer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ContextViewer for view context file
// All context is store as a single file
// And the ContextViewer is used to display the context with different way
// It can be open with $EDITOR, or formatted print, colorful print or other method
type ContextViewer interface {
	// PathView receive the stored context file path
	// Then display it by specific method
	PathView(path string) error

	// StrView receive the content of whole context file
	// Usually formats or colors it for display
	// StrView(content string)
}

// EditorViewer which using terminal editor to view the context file
// It doesn't handle docx, doc or png files
type EditorViewer struct {
	// ViewCommands should be a shell command which receive a file path then display it's markdown content
	// For instance, nvim %s, marktext %s, helix %s ...
	// It support customization on toml configuration file
	ViewCommands []string
}

// PathView for EditorViewer
// use the editor to open context file
// It will try the ViewCommands one by one
func (ev *EditorViewer) PathView(path string) error {
	var err error
	for _, viewCmd := range ev.ViewCommands {
		shellCmd := fmt.Sprintf(viewCmd, path)
		if err = ev.tryShellCmd(shellCmd); err == nil {
			return nil
		}
	}

	return fmt.Errorf("error when open with editor viewer: %s", err.Error())
}

// tryShellCmd try to exec shellCmd
// This is for editor viewer opening the context file
func (ev *EditorViewer) tryShellCmd(shellCmd string) error {
	cmdParts := strings.Split(shellCmd, " ")

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("view command: %s, err: %s", shellCmd, err.Error())
	}

	return nil
}
