// This file just for demo
import "time"

type Energy int

const (
	EnergyLow Energy = iota
	EnergyMedium
	EnergyHigh
)

type WorkStatus int

const (
	StatusTODO WorkStatus = iota
	StatusDOING
	StatusDONE
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
// It doesn't hanle docx, doc or png files
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

// PlainPrintViewer a viewer which just print the context on terminal
type PlainPrintViewer struct{}

// PathView for PlainPrintViewer
// Just read file the print it on terminal
func (ppv *PlainPrintViewer) PathView(path string) error

// BatPrintViewer use bat command to view context
type BatPrintViewer struct{}

func (bpv *BatPrintViewer) PathView(path string) error

// ViewerType for sore type into file
type ViewerType int

const (
	ViewerEditor ViewerType = iota
	ViewerPlainPrint
	ViewerBatPrint
)

type Work struct {
	// Use timestamp hash as id
	ID string

	// EnergyRequirement mark the suitable status for handling this work
	EnergyRequirement Energy

	Status WorkStatus

	// Context provide the useful information about this work
	// It will be store on the file Path
	ContextPath string

	// Viewer for displaying the Context file content
	// Current available types
	// [EditorViewer] [PlainPrintViewer] [BatPrintViewer]
	Viewer ViewerType

	// Time record for logging and work analysis
	CreateTime time.Time

	StartTime time.Time

	EndTime time.Time

	// BlockedWorksID contains the ID of works which is blocked by this work
	// When this work is done, system will update these works with BlockedTimes -= 1
	BlockedWorksID []string

	// BlockedTimes, if not work should be done before begin this work, it will be 0
	// When not 0, it should not be pushed into ReadyQueue
	BlockedTimes int
}

// WorkFilter for filter valid work for current condition
// There is no need to check if work.status is StatusTODO
// The [WorkProvider.Provide] function will skip invalid status
type WorkFilter func(*Work) bool

func EnergyFilter(e Energy) WorkFilter {
	return func(w *Work) bool {
		return w.BlockedTimes == 0 && w.EnergyRequirement == e
	}
}

type WorkProvider struct {
	AllWorks []*Work
}

func (wp *WorkProvider) Provide(filter WorkFilter) *Work {
	for _, work := range wp.AllWorks {
		// Check if status is todo
		// if not skip
		// Check on this provide loop because any work provided should be status todo
		if work.Status != StatusTODO {
			continue
		}

		if filter(work) {
			return work
		}
	}
	return nil
}

func CombinedFiler(filters ...WorkFilter) WorkFilter {
	return func(w *Work) bool {
		for _, f := range filters {
			if !f(w) {
				return false
			}
		}

		return true
	}
}
