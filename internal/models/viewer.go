package models

// ViewerType for sore type into file
type ViewerType int

const (
	ViewerEditor ViewerType = iota
	ViewerPlainPrint
	ViewerBatPrint
)
