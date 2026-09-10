// Package store define the interface [TodoDB]
// which expose the functions
// [TodoDB.Push] [TodoDB.Pop] [TodoDB.Done] [TodoDB.Sync]
// which is used for work database management
// Function details see the define of interface
package store

import (
	"time"

	"github.com/IridiumNan/project-todo/internal/builder"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

// generateTimestampID return a hashed string which generated from timestamp and length = 8
func generateTimestampID(inputTime time.Time) string {
	return utils.HashByTimestamp(inputTime.UnixNano(), 8)
}

type TodoDB interface {
	// Push create a new work then build metadata from user input
	// It generate an ID for this work then store the context file path and it's content on the memory until [TodoDB.Sync] is called
	Push(conf *builder.MDTomlConfig, contextByte []byte) (string, error)

	// Pop next Work
	// If doing data file has work which is doing, pop it first
	// else load todo data file then check if there is an available work
	//
	// This function will modify the pop work status from [models.StatusTODO] to [models.StatusDOING]
	// Then update the work status on the memory
	//
	// You should call [TodoDB.Sync] function to update the data file status
	Pop(filter filter.WorkFilter) (*models.Work, error)

	// Done Change work status from [models.StatusDOING] to [models.StatusDONE]
	// then write this work metadata into data file
	// data file name and path depends on the database format
	// The work status will not be update on disk until function [TodoDB.Sync] is called
	Done(work *models.Work) error

	// Sync the function makes changes on memory saved to disk
	// It contains the work metadata, context file for new work
	Sync() error
}
