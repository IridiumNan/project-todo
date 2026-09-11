// Package builder provide the work builder
// work builder is used for build struct [models.Work]
package builder

import (
	"bufio"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/utils"
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

	mdFile, err := os.CreateTemp(wb.DataDir, "project-todo-markdown-builder-*.md")
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

func (wb *WorkTomlBuilder) Build() error {
	td, err := store.NewTomlDB(wb.DataDir)
	if err != nil {
		return fmt.Errorf("error when create a new toml database connection, err: %s", err)
	}

	idTitleMap := map[string]string{}
	for _, work := range td.TodoWorks {
		idTitleMap[work.ID] = work.Title
	}

	mdFile, err := wb.buildMDTemplate(idTitleMap)
	if err != nil {
		return fmt.Errorf("error when build markdown template, err: %s", err)
	}

	err = utils.OpenWithEnvEditor(mdFile.Name())
	if err != nil {
		return fmt.Errorf("error when open with environment editor, err: %s", err)
	}

	defer func() {
		mdFile.Close()
		os.Remove(mdFile.Name())
	}()

	blockReader := bufio.NewReader(os.Stdin)
	for {

		byteData, err := os.ReadFile(mdFile.Name())
		if err != nil {
			return fmt.Errorf("error when read byte data from md file before parsing it, err: %s", err.Error())
		}
		conf, mdCtx, err := NewMDWorkParser().Parse(byteData)
		if err != nil {
			slog.Error("error when parsing your config file, please edit it again", "err", err)

			fmt.Println("enter for editing again...")
			_, _ = blockReader.ReadString('\n')
			continue
		}

		err = td.Push(conf, mdCtx)
		if err != nil {
			return fmt.Errorf("error when push new work, err: %s", err.Error())
		}

		td.Sync()
		return nil
	}
}

func NewWorkTomlBuilder(dataDir string) *WorkTomlBuilder {
	return &WorkTomlBuilder{
		DataDir: dataDir,
	}
}
