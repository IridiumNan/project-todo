package models

import "time"

type WorkStatus int

const (
	StatusTODO WorkStatus = iota
	StatusDOING
	StatusDONE
)

type Energy int

const (
	EnergyLow Energy = iota
	EnergyMedium
	EnergyHigh
)

// ViewerType for sore type into file
type ViewerType int

const (
	ViewerEditor ViewerType = iota
	ViewerPlainPrint
	ViewerBatPrint
)

type Work struct {
	// Use timestamp hash as id
	ID string

	// EnergyRequirement mark the suitable status for handling this work
	EnergyRequirement Energy

	Status WorkStatus

	// Context provide the useful information about this work
	// It will be store on the file Path
	ContextPath string

	// Viewer for displaying the Context file content
	// Current available types
	// [EditorViewer] [PlainPrintViewer] [BatPrintViewer]
	Viewer ViewerType

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
