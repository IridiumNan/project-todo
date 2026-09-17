package models

type WorkStatus int

// [WorkStatus] mark the [Work] status
const (
	// StatusTODO mark new work which is newly appended and not pop yet
	StatusTODO WorkStatus = iota

	StatusDOING

	StatusDONE

	// Status quit for work which is quit by user
	// That means this work is not reasonable to start or plan
	// TODO: When the shell called, update the quit command with StatusQUIT
	StatusQUIT
)
