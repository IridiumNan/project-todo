package models

type WorkStatus int

const (
	StatusTODO WorkStatus = iota
	StatusDOING
	StatusDONE
)
