package store

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/BurntSushi/toml"
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
	byteData, err := os.ReadFile(dataFilePath)
	if err != nil {
		return fmt.Errorf("error when load works, err: %s", err.Error())
	}

	workParts := bytes.Split(byteData, []byte(defaultTomlSep))

	for idx := range workParts {
		var work models.Work

		if workByte := bytes.Trim(workParts[idx], "\n\t "); string(workByte) == models.EmptyStr {
			// Skip invalid part
			continue
		}

		err := toml.Unmarshal(workParts[idx], &work)
		if err != nil {
			slog.Error("while unmarshal toml metadata", "err", err, "raw_toml_str", string(workParts[idx]))
		}

		td.allWorks = append(td.allWorks, &work)
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
	if f == nil {
		return td.allWorks, nil
	}

	out := []*models.Work{}
	for _, work := range td.allWorks {
		if f(work) {
			out = append(out, work)
		}
	}

	return out, nil
}

func (td *TomlDoneDB) allWithMap() map[string]*models.Work {
	m := make(map[string]*models.Work, len(td.allWorks))

	for _, w := range td.allWorks {
		m[w.ID] = w
	}

	return m
}

func (td *TomlDoneDB) AllWithMap(f filter.WorkFilter) (map[string]*models.Work, error) {
	if f == nil {
		return td.allWithMap(), nil
	}

	m := make(map[string]*models.Work, len(td.allWorks))
	for _, w := range td.allWorks {
		if f(w) {
			m[w.ID] = w
		}
	}

	return m, nil
}
