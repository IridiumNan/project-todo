// Package runner defines the struct [WorkRunner] which is used to display and change status of work then write log into file
// When a runner start a work, it will write metadata of this work into [models.DataDOINGTomlName]
package runner

type WorkRunner interface {
	Run() error
}
