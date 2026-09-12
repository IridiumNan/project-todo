package viewer

// DocViewer use commands to open doc file
// libreoffice or wps to open the context file with format docx, doc, pdf etc
// TODO: Add it's configuration into global config file
type DocViewer struct {
	ViewCommands []string
}
