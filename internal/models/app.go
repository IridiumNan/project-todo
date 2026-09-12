package models

import (
	"fmt"
	"os"
	"strings"
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

const ProjectAppend = `This project is powered by IridiumNan
See <https://github.com/IridiumNan/project-todo>`

func ProjectAppendWithMDQuote() string {
	parts := strings.Split(ProjectAppend, "\n")

	out := ""

	for _, p := range parts {
		out += fmt.Sprintf("> %s\n", p)
	}

	return "\n\n---\n\n" + out
}

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
