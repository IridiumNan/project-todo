package store

import (
	"fmt"
	"path"

	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
)

// TomlDoneDB is an implement of [DoneDB]
type TomlDoneDB struct {
	// store all works on the memory
	allWorks []*models.Work
}

// NewTomlDoneDB create a new toml done database without any metadata
// You should call function [TomlDoneDB.Load] before get data
func NewTomlDoneDB() *TomlDoneDB {
	return &TomlDoneDB{}
}

// Load the works which is read-only for now
func (td *TomlDoneDB) Load(dataDirPath string) error {
	dataFilePath := path.Join(dataDirPath, models.DataDONETomlName)
	var err error
	td.allWorks, err = loadFromTomlFileToSlice(dataFilePath)
	if err != nil {
		td.allWorks = make([]*models.Work, 0)
		return fmt.Errorf("error while loading works from toml file, err: %s", err.Error())
	}

	return nil
}

// Clear all works
func (td *TomlDoneDB) Clear() {
	td.allWorks = []*models.Work{}
}

// Reload clear and load from dataDirPath
func (td *TomlDoneDB) Reload(dataDirPath string) error {
	td.Clear()

	return td.Load(dataDirPath)
}

// All return all works that make f(work) == true
func (td *TomlDoneDB) All(f filter.WorkFilter) ([]*models.Work, error) {
	return filter.NewFilterAssistant().SliceToSlice(td.allWorks, f)
}

func (td *TomlDoneDB) AllWithMap(f filter.WorkFilter) (map[string]*models.Work, error) {
	return filter.NewFilterAssistant().SliceToMap(td.allWorks, f)
}
