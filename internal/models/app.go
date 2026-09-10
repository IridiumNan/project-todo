package models

import (
	"fmt"
	"os"
)

const (
	AppName        = "project-todo"
	EmptyStr       = ""
	DataDirName    = ".project-todo"
	ContextDirName = "context"

	// Toml data file name
	DataTODOTomlName  = "todo-data.toml"
	DataDOINGTomlName = "doing-data.toml"
	DataDONETomlName  = "done-data.toml"
)

var homeDir string = EmptyStr

func GetHomeDir() (string, error) {
	if homeDir != EmptyStr {
		return homeDir, nil
	}

	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return EmptyStr, fmt.Errorf("while getting home dir, err: %s", err)
	}

	return homeDir, nil
}
