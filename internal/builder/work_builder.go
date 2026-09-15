// Package builder provide the work builder
// work builder is used for build struct [models.Work]
package builder

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/store"
)

type WorkBuilder interface {
	Build() error
}

type WorkTomlBuilder struct {
	// DataDir is /path/to/.project-todo
	DataDir string
}

func (wb *WorkTomlBuilder) buildMDTemplate(idTitleMap map[string]string) (*os.File, error) {
	mi := NewMDInjector(idTitleMap)
	mdContent := mi.Inject(newWorkMDTemplate)

	tmpDir := path.Join(wb.DataDir, models.DataBuildDirName)

	mdFile, err := os.CreateTemp(tmpDir, "project-todo-markdown-builder-*.md")
	if err != nil {
		return nil, fmt.Errorf("error when create a tmp file for editing configuration for build a new work, err: %s", err)
	}

	_, err = mdFile.WriteString(mdContent)
	if err != nil {
		return nil, fmt.Errorf("error when write markdown content into tmp file, file_path: %s, err: %s", mdFile.Name(), err.Error())
	}

	err = mdFile.Sync()
	if err != nil {
		return nil, fmt.Errorf("error when sync injected md content into file, err: %s", err.Error())
	}

	return mdFile, nil
}

// Build create a new markdown template file for new command
// use edit this file then work_parser Parse it
func (wb *WorkTomlBuilder) Build(db store.DB) (outFilePath string, err error) {
	works, err := db.All(nil)

	idTitleMap := map[string]string{}

	for idx := range works {
		idTitleMap[works[idx].ID] = works[idx].Title
	}
	mdFile, err := wb.buildMDTemplate(idTitleMap)
	if err != nil {
		return models.EmptyStr, fmt.Errorf("error when build markdown template, err: %s", err)
	}

	defer func() {
		mdFile.Close()
		// os.Remove(mdFile.Name())
	}()

	return mdFile.Name(), nil
}

var ErrParse = errors.New("parse markdown file failed")

func (wb *WorkTomlBuilder) ParseTodoWork(filePath string) (conf *models.MDTomlConfig, mdCtx []byte, err error) {
	var byteData []byte

	byteData, err = os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error when read byte data from md file before parsing it, err: %s", err.Error())
	}

	var rawMDCtx []byte
	conf, rawMDCtx, err = NewMDWorkParser().Parse(byteData)
	if err != nil {
		slog.Error("error when parsing your config file, please edit it again", "err", err)

		return nil, nil, ErrParse
	}
	mdCtx = wb.buildMDCtx(conf, rawMDCtx)

	return
}

func (wb *WorkTomlBuilder) buildMDCtx(conf *models.MDTomlConfig, rawMDCtx []byte) (mdCtx []byte) {
	prefix := wb.mdEnergyPrefix(conf)

	mdCtx = append(mdCtx, prefix...)
	mdCtx = append(mdCtx, rawMDCtx...)
	mdCtx = append(mdCtx, []byte(models.ProjectAppendWithMDQuote())...)

	return mdCtx
}

func (wb *WorkTomlBuilder) mdEnergyPrefix(conf *models.MDTomlConfig) []byte {
	noteEnergyStr := `> [!NOTE]
> energy requirement: %s`
	title := fmt.Sprintf("# %s\n", conf.Title)

	var energyStr string
	switch conf.Energy {
	case models.EnergyLow:
		energyStr = "Low Energy"
	case models.EnergyMedium:
		energyStr = "Medium Energy"
	case models.EnergyHigh:
		energyStr = "High Energy"
	}

	energyInfo := fmt.Sprintf(noteEnergyStr, energyStr)

	prefix := title + "\n" + energyInfo + "\n\n---\n\n"

	return []byte(prefix)
}

func NewWorkTomlBuilder(dataDir string) *WorkTomlBuilder {
	return &WorkTomlBuilder{
		DataDir: dataDir,
	}
}
