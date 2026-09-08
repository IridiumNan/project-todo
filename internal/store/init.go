package store

import (
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/models"
)

func InitDataDir() (err error) {
	currPath, err := os.Getwd()
	if err != nil {
		slog.Error("getting current dir, use . as current path", "err", err)
		currPath = "."
	}

	dataDir := path.Join(currPath, models.DataDirName)

	err = os.Mkdir(dataDir, 0o755)
	if err != nil {
		slog.Error("creating date dir", "err", err, "data_dir_path", dataDir)
		return err
	}

	slog.Info("init new project-todo data dir", "data_dir_path", dataDir)

	contextPath := path.Join(dataDir, "context")

	err = os.Mkdir(contextPath, 0o755)
	if err != nil {
		slog.Error("creating context dir", "context_path", contextPath, "err", err)
		return err
	}

	slog.Info("creating context dir", "context_path", contextPath)

	return
}
