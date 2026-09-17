// Package store define the interface [TodoDB]
// which expose the functions
// [TodoDB.Push] [TodoDB.Pop] [TodoDB.Done] [TodoDB.Sync]
// which is used for work database management
// Function details see the define of interface
package store

import (
	"time"

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
	// if f == nil, return all works loaded
	//
	// For [TodoDB], it just store todo works
	// For [DoneDB], it just store done works
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
