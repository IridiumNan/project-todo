package models

import (
	"time"
)

type Work struct {
	// Use timestamp hash as id
	// Auto generated
	ID string `toml:"id"`

	// Title provide overview information for work dependency reference
	// NOTE: User set
	Title string `toml:"title"`

	// EnergyRequirement mark the suitable status for handling this work
	// NOTE: User set
	EnergyRequirement Energy `toml:"energy"`

	// Auto generated
	// default StatusTODO
	Status WorkStatus `toml:"work_status"`

	// Context provide the useful information about this work
	// It will be store on the file Path
	// Auto generated
	ContextPath string `toml:"context_path"`

	// Viewer for displaying the Context file content
	// Current available types
	// [viewer.EditorViewer] [viewer.PlainPrintViewer] [viewer.BatPrintViewer]
	// NOTE: User Set (global default configuration is supported)
	Viewer ViewerType `toml:"viewer"`

	// Time record for logging and work analysis
	// Auto generated
	CreateTime time.Time `toml:"create_time"`

	StartTime time.Time `toml:"start_time"`

	EndTime time.Time `toml:"end_time"`

	// BlockedWorksID contains the ID of works which is blocked by this work
	// When this work is done, system will update these works with BlockedTimes -= 1
	// NOTE: User set
	BlockedWorksID []string `toml:"blocked_works_id"`

	// BlockedTimes, if not work should be done before begin this work, it will be 0
	// When not 0, it should not be pushed into ReadyQueue
	// Auto generated
	BlockedTimes int `toml:"blocked_times"`
}
