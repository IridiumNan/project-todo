package models

import (
	"fmt"
	"os"
)

const (
	AppName     = "project-todo"
	EmptyStr    = ""
	DataDirName = ".project-todo"
)

var HomeDir string = EmptyStr

func GetHomeDir() (string, error) {
	if HomeDir != EmptyStr {
		return HomeDir, nil
	}

	var err error
	HomeDir, err = os.UserHomeDir()
	if err != nil {
		return EmptyStr, fmt.Errorf("while getting home dir", "err", err)
	}

	return HomeDir, nil
}
