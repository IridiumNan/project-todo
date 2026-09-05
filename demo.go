package main

import (
	"time"
)

type Energy int

const (
	EnergyLow Energy = iota
	EnergyMedium
	EnergyHigh
)

// WorkStatus mark the status of a work
// it can be TODO, DOING, DONE
type WorkStatus int

const (
	StatusTODO WorkStatus = iota
	StatusDOING
	StatusDONE
)

type WorkContext interface {
	// Provide offer information about specific work with suitable way
	// For school work, it can just print the link which will be used
	// For program development, it can use $EDITOR to open document
	Provide()
}

// WorkDoc for manage the document about contains details about this work
// usually for personal project development
type WorkDoc struct {
	Path string
}

const WorkDocSuffix = "doc"

// Provide for WorkDoc use $EDITOR to open related document
func (d *WorkDoc) Provide()

// WorkLink usually for school homework
type WorkLink struct {
	Link string

	Comment string
}

const WorkLinkSuffix = "link"

func (l *WorkLink) Provide()

// You can add more struct for context

// NewContext Try all possible context type for this file
func NewContext(contextPath string) WorkContext

type Work struct {
	// Use timestamp hash as id
	ID string

	// EnergyRequirement mark the suitable status for handling this work
	EnergyRequirement Energy

	Status WorkStatus

	// Context provide the useful information about this work
	// It will be store on the file Path
	// system will load it into specific struct when it's used
	// See [WorkContext.Provide]
	// The ContextPath should use meaningful suffix for file name
	// for [WorkDoc], it should be [WorkDocSuffix]
	// for [WorkLink], it should be [WorkLinkSuffix]
	// So that system can identify different types easier
	ContextPath string

	// Time record for logging and work analysis
	CreateTime time.Time

	StartTime time.Time

	EndTime time.Time

	// BlockedWorksID contains the ID of works which is blocked by this work
	// When this work is done, system will update these works with BlockedTimes -= 1
	BlockedWorksID []string

	// BlockedTimes, if not work should be done before begin this work, it will be 0
	// When not 0, it should not be pushed into ReadyQueue
	BlockedTimes int
}

// ReadyQueue which store all task with BlockedCount == 0 & Status == TODO
// The key is the ID which is the hash value of create timestamp
type ReadyQueue map[string]*Work

type WorkEngine struct {
	// LowQueue store the task which contains tasks whose Energy is EnergyLow and BlockedCount == 0
	LowQueue ReadyQueue

	// MediumQueue tasks whose Energy is EnergyMedium and BlockedCount == 0
	MediumQueue ReadyQueue

	// HighQueue tasks whose Energy is EnergyHigh and BlockedCount == 0
	HighQueue ReadyQueue

	AllWorks []*Work
}
