package viewer

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
	// For instance, nvim %s, marktext %s, bat %s ...
	// It support customization on toml configuration file
	ViewCommands []string
}

// PathView for EditorViewer
// use the editor to open context file
// It will try the ViewCommands one by one
func (mv *EditorViewer) PathView(path string) error
