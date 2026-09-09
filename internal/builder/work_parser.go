package builder

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/models"
)

// mdTomlConf struct is used for UnMarshal raw markdown content
type mdTomlConfig struct {
	Title string `toml:"title"`

	Energy models.Energy `toml:"energy"`

	Viewer models.ViewerType `toml:"viewer"`

	DependenciesID []string `toml:"dependencies_id"`
}

type MDWorkParser struct{}

const (
	mdDefaultSep   = "+++\n"
	mdContextTitle = "## CONTEXT\n"
)

// Parse parse the raw markdown file then return the mdTomlConf and mdContext
// Then WorkBuilder use there data to build a new Work
func (wp *MDWorkParser) Parse(rawMD []byte) (mdTomlConf *mdTomlConfig, mdContext []byte, err error) {
	mdParts := bytes.SplitN(rawMD, []byte(mdDefaultSep), 3)

	if len(mdParts) < 3 {
		slog.Error("invalid length of md Parts loaded from the raw markdown content", "expected_length", 3, "current_length", len(mdParts))
		return nil, nil, fmt.Errorf("error when parse markdown: len of mdParts invalid")
	}

	if string(mdParts[0]) != models.EmptyStr {
		slog.Warn("first markdown part is not empty, ignore it", "content", mdParts[0])
	}

	tomlConfigData := mdParts[1]

	tomlConf := mdTomlConfig{}

	if err := toml.Unmarshal(tomlConfigData, &tomlConf); err != nil {
		return nil, nil, fmt.Errorf("error when Unmarshal toml config byte data from the markdown parts, err: %s", err)
	}

	// the third parts should contains the ## CONTEXT
	// Use the [bytes.Cut] function to get context then build a new context string
	_, mdContext, found := bytes.Cut(rawMD, []byte(mdContextTitle))
	if !found {
		return nil, nil, fmt.Errorf("error when cutting markdown context part, err: %s", err.Error())
	}

	mdContext = bytes.Trim(mdContext, " \n")

	return &tomlConf, mdContext, nil
}
