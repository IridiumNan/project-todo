package store

import (
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
)

// TomlArchiveDB load all archive works then update their status as [models.StatusArchive]
// Then dump data into
type TomlArchiveDB struct {
	dataDir string

	archiveWorks []*models.Work
}

// NewTomlArchiveDB return a new toml archive database based on the datadir
func NewTomlArchiveDB(dataDir string) *TomlArchiveDB {
	tomlFilePath := path.Join(dataDir, models.DataArchiveTomlName)

	archivedWorks, err := loadFromTomlFileToSlice(tomlFilePath)
	if err != nil {
		slog.Error("while loading archive works data", "err", err)
		archivedWorks = make([]*models.Work, 0)
	}

	return &TomlArchiveDB{
		dataDir:      dataDir,
		archiveWorks: archivedWorks,
	}
}

func (td *TomlArchiveDB) All(f filter.WorkFilter) ([]*models.Work, error) {
	return filter.NewFilterAssistant().SliceToSlice(td.archiveWorks, f)
}

func (td *TomlArchiveDB) AllWithMap(f filter.WorkFilter) (map[string]*models.Work, error) {
	return filter.NewFilterAssistant().SliceToMap(td.archiveWorks, f)
}

// updateStatus update [TomlArchiveDB.archiveWorks] status as [models.StatusArchive]
func (td *TomlArchiveDB) updateStatus() {
	for _, work := range td.archiveWorks {
		work.Status = models.StatusArchive
	}
}

// isDuplicate check if this work is duplicate (exists on current archiveWorks)
// using work id as the identity
func (td *TomlArchiveDB) isDuplicate(work *models.Work) bool {
	for _, w := range td.archiveWorks {
		if work.ID == w.ID {
			return true
		}
	}

	return false
}

// Append a new work, if this work id has exist ont archiveWorks, it return [os.ErrExist]
func (td *TomlArchiveDB) Append(work *models.Work) error {
	if td.isDuplicate(work) {
		return os.ErrExist
	}

	td.archiveWorks = append(td.archiveWorks, work)
	return nil
}

// Sync build the toml byte data then write the byte data into file
// return the [os.WriteFile] error
func (td *TomlArchiveDB) Sync() error {
	td.updateStatus()

	byteData := buildTomlByteData(td.archiveWorks, []string{})

	archiveFile := path.Join(td.dataDir, models.DataArchiveTomlName)

	return os.WriteFile(archiveFile, byteData, 0o644)
}

// AllWorkFilter return a filter
// If a work has archived, it will return false
func (td *TomlArchiveDB) AllWorkFilter() filter.WorkFilter {
	return func(w *models.Work) bool {
		for _, word := range td.archiveWorks {
			if w.ID == word.ID {
				return false
			}
		}
		return true
	}
}
