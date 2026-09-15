package runner

import "github.com/IridiumNan/project-todo/internal/store"

// WorkDoneRunner provide read-only mode for works whose status is [models.StatusDONE]
// It construct a new tmp file about specific done work with metadata and duration of this work
type WorkDoneRunner struct {
	// ViewDB provide the read-only context file path and metadata
	ViewDB store.DoneDB
}

func Run() error {
	return nil
}

func Wait() (done bool) {
	return true
}
