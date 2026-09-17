package viewer

import (
	"github.com/IridiumNan/project-todo/internal/config"
	"github.com/IridiumNan/project-todo/internal/models"
)

// GetWorkViewer get the Viewer struct by the [WorkDefaultRunner.Work] viewerType
func GetWorkViewer(t models.ViewerType) (ctxViewer ContextViewer) {
	switch t {
	case models.ViewerBatPrint:
		return NewBatPrintViewer()
	case models.ViewerPlainPrint:
		return NewPlainPrintViewer()
	case models.ViewerEditor:
		return NewEditorViewer(config.GlobalConf.EditorViewCommands)
	}

	return nil
}
