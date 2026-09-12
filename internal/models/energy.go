package models

type Energy int

const (
	EnergyLow Energy = iota
	EnergyMedium
	EnergyHigh

	EnergyInvalid
)
