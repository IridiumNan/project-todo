package viewer

// PlainPrintViewer a viewer which just print the context on terminal
type PlainPrintViewer struct{}

// PathView for PlainPrintViewer
// Just read file the print it on terminal
func (ppv *PlainPrintViewer) PathView(path string) error

// BatPrintViewer use bat command to view context
type BatPrintViewer struct{}

func (bpv *BatPrintViewer) PathView(path string) error
