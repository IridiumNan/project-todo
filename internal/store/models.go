// Package store define the interface [TodoDB]
// which expose the functions
// [TodoDB.Push] [TodoDB.Pop] [TodoDB.Done] [TodoDB.Sync]
// which is used for work database management
// Function details see the define of interface
package store

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID(inputTime time.Time) string {
	return utils.HashByTimestamp(inputTime.UnixNano(), 8)
}

type DB interface {
	// All return all works for filter return true
	// if f == nil, use [filter.DefaultFilter]
	//
	// For [TodoDB], it just store todo works
	// For [DoneDB], it just store done works
	//
	All(f filter.WorkFilter) ([]*models.Work, error)

	// AllWithMap function returns all works match the filter
	// The map key is [models.Work.ID]
	AllWithMap(f filter.WorkFilter) (map[string]*models.Work, error)
}

// TodoDB provide todo work with Pop function and support Push new todo work
// Use Done function to change the status on database
type TodoDB interface {
	// Push create a new work then build metadata from user input
	// It generate an ID for this work then store the context file path and it's content on the memory until [TodoDB.Sync] is called
	Push(conf *models.MDTomlConfig, contextByte []byte) (string, error)

	// Pop next Work
	// If doing data file has work which is doing, pop it first
	// else load todo data file then check if there is an available work
	//
	// this function will not manage the status of work, the status update will be handled by runner
	Pop(filter filter.WorkFilter) (*models.Work, error)

	// then write this work metadata into data file
	// data file name and path depends on the database format
	Done(work *models.Work) error

	DB

	// Sync the function makes changes on memory saved to disk
	// It contains the work metadata, context file for new work
	Sync() error
}

// DoneDB Provide read-only functions to visit works
type DoneDB interface {
	// Load all works from the data dir
	// ViewDB will search toml data file from this data dir
	// This will maintain all works has been loaded
	Load(dataDirPath string) error

	// Reload Equals to Clear then Load
	Reload(dataDirPath string) error

	// Clear remove all works on the database (memory), it will not change the database file, just remove loaded content
	Clear()

	DB
}

// buildTomlByteData convert all works into the toml file format with [defaultTomlSep]
// It will exclude all work in excludeID whose context file is broken or fail to write
func buildTomlByteData(works []*models.Work, excludeID []string) (byteData []byte) {
	for _, work := range works {
		if slices.Contains(excludeID, work.ID) {
			slog.Warn("skip work with broken context", "id", work.ID)
			continue
		}
		data, currErr := toml.Marshal(work)
		if currErr != nil {
			slog.Error("while marshal work into byte data", "err", currErr, "work", work)
			continue
		}

		byteData = append(byteData, data...)
		// add the [defaultTomlSep] for sep different todo work variant
		// WARN: DON NOT TOUCH THIS
		// or system will fail to parse the metadata file
		byteData = append(byteData, []byte(defaultTomlSep)...)
	}

	return
}

func loadFromTomlFileToSlice(tomlFilePath string) ([]*models.Work, error) {
	if _, err := os.Stat(tomlFilePath); os.IsNotExist(err) {
		slog.Warn("toml file not found, creating a new one", "path", tomlFilePath)
		// create new file then return a empty slice
		_, err := os.Create(tomlFilePath)
		if err != nil {
			return nil, fmt.Errorf("error while creating a new toml file, err: %s", err.Error())
		}

		return make([]*models.Work, 0), nil
	}
	byteData, err := os.ReadFile(tomlFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when load works, err: %s", err.Error())
	}

	workParts := bytes.Split(byteData, []byte(defaultTomlSep))

	allWorks := make([]*models.Work, 0, 20)

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

		allWorks = append(allWorks, &work)
	}

	return allWorks, nil
}

func loadFromTomlFileToMap(tomlFilePath string) (map[string]*models.Work, error) {
	if _, err := os.Stat(tomlFilePath); os.IsNotExist(err) {
		slog.Warn("toml file not found, creating a new one", "path", tomlFilePath)
		_, err := os.Create(tomlFilePath)
		if err != nil {
			return nil, fmt.Errorf("error while creating a new toml file, err: %s", err.Error())
		}

		return make(map[string]*models.Work), nil
	}

	byteData, err := os.ReadFile(tomlFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when load works, err: %s", err.Error())
	}
	workParts := bytes.Split(byteData, []byte(defaultTomlSep))

	allWorks := make(map[string]*models.Work, 20)

	for idx := range workParts {
		var work models.Work

		if workByte := bytes.Trim(workParts[idx], "\n\t "); string(workByte) == models.EmptyStr {
			// skip invalid part
			continue
		}

		err := toml.Unmarshal(workParts[idx], &work)
		if err != nil {
			slog.Error("while unmarshal toml metadata", "err", err, "raw_toml_str", string(workParts[idx]))
		}

		allWorks[work.ID] = &work
	}
	return allWorks, nil
}
