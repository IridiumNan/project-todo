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
	defaultTomlSep = "@@@@@GAP@@@@@\n"
)

type pathType string

type metaWriteTaskFunc func() error

type TomlDB struct {
	// todoWorks store id as key
	// key is autogenerate by timestamp
	// this map just store works with status [models.StatusTODO]
	todoWorks map[string]*models.Work

	// the data dir named .project-todo where [tomlDataFileName] and context dir is stored
	DataDirPath string

	// Write task will be write to disk only when [TomlDB.Sync] is called

	// ctxWriteTask cache the context file
	ctxWriteTask map[pathType]([]byte)

	// metaWriteTask exec when [TomlDB.Sync] is called
	// NOTE: just use this map when update for statusDOING and statusDONE
	// todo works just dump [TomlDB.todoWorks] into todo metadata file
	metaWriteTask map[pathType]metaWriteTaskFunc
}

func (td *TomlDB) addCtxWriteTask(CtxFilePath string, data []byte) {
	td.ctxWriteTask[pathType(CtxFilePath)] = data
}

func (td *TomlDB) delCtxWriteTask(CtxFilePath string, data []byte) {
	delete(td.ctxWriteTask, pathType(CtxFilePath))
}

func (td *TomlDB) addMetaWriteFunc(metaFilePath string, f metaWriteTaskFunc) {
	td.metaWriteTask[pathType(metaFilePath)] = f
}

// resolvePushDependency when a new work is pushed, resolve the dependencies
// That's push new id into it's dependency [models.Work.BlockedWorksID]
// This support prefix match
func (td *TomlDB) resolvePushDependency(newID string, dependencies []string) {
	for _, targetPrefix := range dependencies {
		for id := range td.todoWorks {
			// use prefix match
			if strings.HasPrefix(id, targetPrefix) {
				td.todoWorks[id].BlockedWorksID = append(td.todoWorks[id].BlockedWorksID, newID)
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

	td.todoWorks[newWorkID] = &newWork

	// write this when [TomlDB.Sync] is called
	td.addCtxWriteTask(contextFilePath, contextByte)

	return nil
}

// loadDoingWork load if doing file exist (which means that there is doing work now)
// if file not exist, return [os.ErrNotExist]
// else encounter error when load metadata from file, you should check the logic
// The metadata file just contains one task without sep
func (td *TomlDB) loadDoingWork() (*models.Work, error) {
	doingFilePath := path.Join(td.DataDirPath, models.DataDOINGTomlName)

	if _, err := os.Stat(doingFilePath); os.IsNotExist(err) {
		// doing file not exist, return now
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(doingFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when read doingFilePath, err: %s", err.Error())
	}
	var work models.Work
	err = toml.Unmarshal(data, &work)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshal status doing work metadata, err: %s", err.Error())
	}

	return &work, nil
}

func (td *TomlDB) fetchWorkByFilter(filter filter.WorkFilter) (*models.Work, error) {
	var work *models.Work

	for _, todoWork := range td.todoWorks {
		if todoWork.BlockedTimes == 0 && filter(work) {
			work = todoWork
		}
	}

	if work == nil {
		return nil, fmt.Errorf("error when fetch work by filter, there is no valid work")
	}

	return work, nil
}

func (td *TomlDB) Pop(filter filter.WorkFilter) (*models.Work, error) {
	var poppedWork *models.Work = nil
	var err error

	// check if there is work whose status is [models.StatusDOING]
	// if load success, return now
	poppedWork, err = td.loadDoingWork()
	if err == nil {
		return poppedWork, nil
	}

	if os.IsNotExist(err) {
		slog.Info("not doing work found, try todo work")
	}

	slog.Warn("fail to load doing work from metadata", "err", err)

	work, err := td.fetchWorkByFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("error when pop todo work, err: %s", err.Error())
	}

	doingFilePath := path.Join(td.DataDirPath, models.DataDOINGTomlName)
	td.addMetaWriteFunc(doingFilePath, func() error {
		data, err := toml.Marshal(work)
		if err != nil {
			return fmt.Errorf("error when marshal work data, err: %s", err.Error())
		}

		err = os.WriteFile(doingFilePath, data, 0o644)
		if err != nil {
			return fmt.Errorf("error when write metadata to file, file path: %s, err: %s", doingFilePath, err.Error())
		}

		return nil
	})

	return work, nil
}

// resolveDoneDependency call this function when a work status changed from [models.StatusDOING] to [models.StatusDONE]
func (td *TomlDB) resolveDoneDependency(blockedIDs []string) {
	for _, id := range blockedIDs {
		td.todoWorks[id].BlockedTimes -= 1
	}
}

// NOTE: Remove the metadata file for doing work
func (td *TomlDB) Done(work *models.Work) error {
	td.resolveDoneDependency(work.BlockedWorksID)
	// TODO: append this into [models.DataDONETomlName]
}

// Sync -> See interface [TodoDB] Sync
func (td *TomlDB) Sync() error {
	// TODO: rewrite this function for write all metadata
	// And context file name
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
	for ctxPath, byteData := range td.ctxWriteTask {
		currErr := os.WriteFile(string(ctxPath), byteData, 0o644)
		if currErr != nil {
			newErr := fmt.Errorf("file path, file path: %s, err: %s", ctxPath, currErr)
			errs = append(errs, newErr)

			id, _ := utils.ParseCtxPath(string(ctxPath))
			excludeID = append(excludeID, id)
		}
	}
	return
}

func (td *TomlDB) buildTomlByTodoWorks(excludeID []string) (byteData []byte) {
	works := utils.MapValues(td.todoWorks)

	return td.buildTomlByteData(works, excludeID)
}

// buildTomlByteData convert all works into the toml file format with [defaultTomlSep]
// It will exclude all work in excludeID whose context file is broken or fail to write
func (td *TomlDB) buildTomlByteData(works []*models.Work, excludeID []string) (byteData []byte) {
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

// NewTomlDB return the TomlDB which is the struct of interface [TodoDB]
// You can use Add, Pop, and Sync functions
func NewTomlDB(dataDirPath string) (*TomlDB, error) {
	allTodoWorks := map[string]*models.Work{}

	// NOTE:
	// just load works with status todo
	dataFilePath := path.Join(dataDirPath, models.DataTODOTomlName)

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

		allTodoWorks[work.ID] = &work
	}

	return &TomlDB{
		todoWorks:   allTodoWorks,
		DataDirPath: dataDirPath,

		ctxWriteTask:  map[pathType][]byte{},
		metaWriteTask: map[pathType]metaWriteTaskFunc{},
	}, nil
}
