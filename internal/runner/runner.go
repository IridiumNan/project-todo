// Package runner defines the struct [WorkRunner] which is used to display and change status of work then write log into file
// When a runner start a work, it will write metadata of this work into [models.DataDOINGTomlName]
package runner

import (
	"log/slog"

	"github.com/IridiumNan/project-todo/internal/models"
)

// TODO: If this runner should manage multi works ?
// NO, just a single work

type WorkRunner interface {
	Run() error
}

// WorkTomlRunner is used to exec the work display command and monitor work status
// When this work' status changed, it write log into file
type WorkTomlRunner struct {
	Work *models.Work

	Logger *slog.Logger
}

func NewWorkRunner() (*WorkRunner, error)

// Run starts the specified command and waits for it to complete.
func (wr *WorkTomlRunner) Run() error {
	return nil
}

// TODO: decide if this function should start work then not block terminal ?
// func (wt *WorkTomlRunner) Start() error {
// 	return nil
// }
