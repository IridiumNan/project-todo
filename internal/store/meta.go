// Package store contains interface [Store] and struct [StoreToml]
// They are used for read and dump metadata into [io.Writer]
package store

import (
	"io"

	"github.com/IridiumNan/project-todo/internal/models"
)

// Store a interface for change meta data
type Store interface {
	ReadMetadata(dataSrc io.Reader) (works []*models.Work, err error)

	WriteMetadata(dataDst io.Writer, Works []*models.Work) (err error)
}
