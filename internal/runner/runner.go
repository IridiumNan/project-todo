// Package runner provide the struct which update work status then write log
// And monitor the work status
package runner

import (
	"log/slog"

	"github.com/IridiumNan/project-todo/internal/models"
)

// TODO: If this runner should manage multi works ?

// WorkRunner is used to exec the work display command and monitor work status
// When this work' status changed, it write log into file
type WorkRunner struct {
	Work *models.Work

	Logger *slog.Logger
}

func NewWorkRunner() (*WorkRunner, error)

// Run start the work then blocked, waiting for the user end this work
func (wr *WorkRunner) Run() error {
	return nil
}

// TODO: decide if this function should start work then not block terminal ?
func (wt *WorkRunner) Start() error {
	return nil
}
