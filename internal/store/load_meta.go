package store

import (
	"io"
	"os"

	"github.com/IridiumNan/project-todo/internal/models"
)

// Store a interface for change meta data
type Store interface {
	// LoadMeta function load all works from
	LoadMetadata(dataSrc io.ReadCloser) (works []*models.Work, err error)

	// AppendMetadata function append single work struct into the database
	AppendMetadata(dataDst io.WriteCloser, newWorks []*models.Work) (err error)

	// DumpMetadataToFile function trunc the old data then write all new metadata
	DumpMetadataToFile(dstFile *os.File, allWorks []*models.Work) (err error)
}
