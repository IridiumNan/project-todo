package store

import (
	_ "embed"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/models"
)

const readmeFileName = "README.md"

//go:embed README.md
var dataDirREADME []byte

// InitDataDirOnCurrentDir create data dir and context dir on current path
// return the dataPath and err
func InitDataDirOnCurrentDir() (dataPath string, err error) {
	currPath, err := os.Getwd()
	if err != nil {
		slog.Error("getting current dir, use . as current path", "err", err)
		currPath = "."
	}

	dataDir := path.Join(currPath, models.DataDirName)

	return dataDir, InitDataDir(dataDir)
}

// InitDataDir create the [models.DataDirName] on dir
// if the dir itself contains [models.DataDirName]. this function will treat it as the dataDirPath
func InitDataDir(dir string) (err error) {
	if path.Base(dir) != models.DataDirName {
		dir = path.Join(dir, models.DataDirName)
	}

	err = os.Mkdir(dir, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			slog.Error("creating date dir", "err", err, "data_dir_path", dir)
			return err
		}
		err = nil
		slog.Warn("data dir exist", "path", dir)
	}

	slog.Info("init new project-todo data dir", "data_dir_path", dir)

	slog.Info("add a readme file on data dir")
	err = addDirReadme(dir)
	if err != nil {
		slog.Error("when add readme file on data dir", "err", err)
	}

	contextPath := path.Join(dir, "context")

	err = os.Mkdir(contextPath, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			slog.Error("creating context dir", "context_path", contextPath, "err", err)
			return err
		}
		err = nil
		slog.Warn("context dir exist", "path", dir)
	}

	buildPath := path.Join(dir, models.DataBuildDirName)

	err = os.Mkdir(buildPath, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			slog.Error("creating build dir", "build_path", buildPath, "err", err)
			return err
		}
		err = nil
		slog.Warn("build dir exist", "path", dir)
	}

	slog.Info("creating build dir", "build_path", buildPath)

	summaryPath := path.Join(dir, models.DataSummaryDirName)

	err = os.Mkdir(summaryPath, 0o755)
	if err != nil {
		if !os.IsExist(err) {
			slog.Error("creating summary dir", "summary_path", buildPath, "err", err)
			return err
		}
		err = nil
		slog.Warn("summary dir exist", "path", dir)
	}

	slog.Info("creating summary dir", "summary_path", summaryPath)

	return
}

func addDirReadme(dir string) error {
	readmePath := path.Join(dir, readmeFileName)

	return os.WriteFile(readmePath, dataDirREADME, 0o444)
}
