package store

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
)

// TomlViewDB is an implement of [ViewDB]
type TomlViewDB struct {
	// store all works on the memory
	allWorks []*models.Work
}

func (td *TomlViewDB) Load(dataFilePath string) error {
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

func (td *TomlViewDB) Clear() {
	td.allWorks = []*models.Work{}
}

func (td *TomlViewDB) Reload(dataFilePath string) error {
	td.Clear()

	return td.Load(dataFilePath)
}

func (td *TomlViewDB) All(f filter.WorkFilter) ([]*models.Work, error) {
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
