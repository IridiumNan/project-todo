// Package store contains interface [Store] and struct [StoreToml]
// They are used for read and dump metadata into [io.Writer]
package store

import (
	"io"
	"time"

	"github.com/IridiumNan/project-todo/internal/builder"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

// Store a interface for change meta data
type Store interface {
	ReadMetadata(dataSrc io.Reader) (works []*models.Work, err error)

	WriteMetadata(dataDst io.Writer, Works []*models.Work) (err error)
}

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID(inputTime time.Time) string {
	return utils.HashByTimestamp(inputTime.UnixNano(), 8)
}

type TodoDB interface {
	// Push new Work
	// dependencies should be [models.StatusDONE] before this work pop
	// return id of this new work and error
	Push(conf *builder.MDTomlConfig, dependencies []string) (string, error)

	// Pop next Work
	// It assume that this work will be done this time
	// So update blocked_works_id times
	Pop(filter *filter.WorkFilter) (*models.Work, error)

	// Sync data into the database file
	Sync() error
}
