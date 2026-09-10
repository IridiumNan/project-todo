package store

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/builder"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

const (
	tomlDataFileName = "data.toml"
	defaultTomlSep   = "@@@@@GAP@@@@@\n"
)

type TomlDB struct {
	// AllWorks store id as key and [models.Work] as value
	AllWorks map[string]*models.Work

	// the data dir named .project-todo where [tomlDataFileName] and context dir is stored
	DataDirPath string

	// ContextToWrite store context on memory
	// The key is context path (which contains the id and title of a work)
	// value is context to be write into disk file
	// Write then Empty this when [TomlDB.Sync] is called
	ContextToWrite map[string]([]byte)
}

// NewTomlDB return the TomlDB which is the struct of interface [TodoDB]
// You can use Add, Pop, and Sync functions
func NewTomlDB(dataDirPath string) (*TomlDB, error) {
	allWorks := map[string]*models.Work{}

	dataFilePath := path.Join(dataDirPath, tomlDataFileName)

	byteData, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when read data file, file path: %s, err: %s", dataFilePath, err.Error())
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

		allWorks[work.ID] = &work
	}

	return &TomlDB{
		AllWorks:    allWorks,
		DataDirPath: dataDirPath,
	}, nil
}

// resolvePushDependency when a new work is pushed, resolve the dependencies
// That's push new id into it's dependency [models.Work.BlockedWorksID]
// This support prefix match
func (td *TomlDB) resolvePushDependency(newID string, dependencies []string) {
	for _, targetPrefix := range dependencies {
		for id := range td.AllWorks {
			// use prefix match
			if strings.HasPrefix(id, targetPrefix) {
				td.AllWorks[id].BlockedWorksID = append(td.AllWorks[id].BlockedWorksID, newID)
			}
		}
	}
}

// Push -> See interface [TodoDB] Push
func (td *TomlDB) Push(conf *builder.MDTomlConfig, contextByte []byte) error {
	now := time.Now()
	newWorkID := generateTimestampID(now)

	contextFileName := fmt.Sprintf("%s-%s", newWorkID, conf.Title)
	contextFilePath := path.Join(td.DataDirPath, models.ContextDirName, contextFileName)

	td.resolvePushDependency(newWorkID, conf.DependenciesID)

	newWork := models.Work{
		ID:                newWorkID,
		Title:             conf.Title,
		EnergyRequirement: conf.Energy,
		// init with status todo
		Status:      models.StatusTODO,
		ContextPath: contextFilePath,
		Viewer:      conf.Viewer,

		// WARN:
		// Just init with now
		// Then TodoDB should not update these values again
		CreateTime:     now,
		StartTime:      now,
		EndTime:        now,
		BlockedWorksID: []string{},
		BlockedTimes:   len(conf.DependenciesID),
	}

	td.AllWorks[newWorkID] = &newWork

	// write this when [TomlDB.Sync] is called
	td.ContextToWrite[contextFilePath] = contextByte

	return nil
}

// resolveDoneDependency call this function when a work status changed from [models.StatusDOING] to [models.StatusDONE]
func (td *TomlDB) resolveDoneDependency(blockedIDs []string) {
	for _, id := range blockedIDs {
		td.AllWorks[id].BlockedTimes -= 1
	}
}

// FIX: REWRITE THIS FUNCTION
// Pop -> See interface [TodoDB] Pop
func (td *TomlDB) Pop(filter filter.WorkFilter) (*models.Work, error) {
	var poppedWork *models.Work = nil

	for _, work := range td.AllWorks {
		if work.Status != models.StatusTODO {
			continue
		}
		if work.BlockedTimes == 0 && filter(work) {
			poppedWork = work
		}
	}

	if poppedWork == nil {
		return nil, fmt.Errorf("error when pop work, there is no valid work")
	}

	// td.resolvePopDependency(poppedWork.BlockedWorksID)
	// Dependencies resolve is handle by function [TomlDB.Done]

	// The StarTime and EndTime should set by [runner.WorkRunner]

	return poppedWork, nil
}

func (td *TomlDB) Done(work *models.Work) error {
	td.resolveDoneDependency(work.BlockedWorksID)
}

// Sync -> See interface [TodoDB] Sync
func (td *TomlDB) Sync() error {
	tmpFile, err := os.CreateTemp("/tmp", "project-todo-data-*")
	if err != nil {
		return fmt.Errorf("error when create tmp file, data file not changed, err: %s", err.Error())
	}

	excludeID, errs := td.writeContext()
	for _, err := range errs {
		slog.Error("while writing context file", "err", err)
	}

	byteData := td.buildTomlByteData(excludeID)

	if _, err := tmpFile.Write(byteData); err != nil {
		return fmt.Errorf("error when write byte data into tmp file, data file not changed, err: %s", err.Error())
	}

	err = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("error when closing tmp file, data file not changed, err: %s", err.Error())
	}

	dataFilePath := path.Join(td.DataDirPath, tomlDataFileName)
	err = os.Rename(tmpFile.Name(), dataFilePath)
	if err != nil {
		return fmt.Errorf("error when replace data file with tmp file, tmp file path: %s, data file path: %s, err: %s", tmpFile.Name(), dataFilePath, err)
	}

	return nil
}

// writeContext try to write context on [TomlDB.ContextToWrite]
// if failed, push it to excludeID and errs
func (td *TomlDB) writeContext() (excludeID []string, errs []error) {
	for ctxPath, byteData := range td.ContextToWrite {
		currErr := os.WriteFile(ctxPath, byteData, 0o644)
		if currErr != nil {
			newErr := fmt.Errorf("file path, file path: %s, err: %s", ctxPath, currErr)
			errs = append(errs, newErr)

			id, _ := utils.ParseCtxPath(ctxPath)
			excludeID = append(excludeID, id)
		}
	}
	return
}

// buildTomlByteData convert all works into the toml file format with [defaultTomlSep]
// It will exclude all work in excludeID whose context file is broken or fail to write
func (td *TomlDB) buildTomlByteData(excludeID []string) (byteData []byte) {
	for _, work := range td.AllWorks {
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
		byteData = append(byteData, []byte(defaultTomlSep)...)
	}

	return
}
