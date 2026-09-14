package runner

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/IridiumNan/project-todo/internal/config"
	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/viewer"
)

// WorkDefaultRunner is used to exec the work display command and monitor work status
// It just handle the work with status [models.StatusTODO]
// When this work' status changed, it write log into file
type WorkDefaultRunner struct {
	// dataDir path to [models.DataDirName]
	dataDir string

	// work the work to be run
	work *models.Work

	// viewer load the work viewer [models.ViewerType] to [viewer.ContextViewer]
	viewer viewer.ContextViewer

	logger *slog.Logger

	db store.TodoDB

	filter filter.WorkFilter
}

// NewWorkDefaultRunner create a new work runner based on the data dir and work pointer
func NewWorkDefaultRunner(dir string, log io.Writer, d store.TodoDB, f filter.WorkFilter) *WorkDefaultRunner {
	handler := slog.NewTextHandler(log, nil)
	logger := slog.New(handler)

	return &WorkDefaultRunner{
		dataDir: dir,
		logger:  logger,
		db:      d,
		filter:  f,
	}
}

// prepare pop a new work from the database by filter
func (r *WorkDefaultRunner) prepare() error {
	work, err := r.db.Pop(r.filter)
	if err != nil {
		return err
	}

	r.viewer = getWorkViewer(work.Viewer)
	r.work = work
	return nil
}

// Run starts the specified command and waits for it to complete.
func (r *WorkDefaultRunner) Run() error {
	err := r.prepare()
	if err != nil {
		return fmt.Errorf("while prepare running, err: %s", err.Error())
	}

	r.work.StartTime = time.Now()
	r.work.Status = models.StatusDOING
	r.logger.Info("starting work", "id", r.work.ID, "title", r.work.Title)

	done := r.Wait()
	if !done {
		r.logger.Warn("exit work without done", "id", r.work.ID, "title", r.work.Title)
		return nil
	}

	r.work.EndTime = time.Now()
	r.work.Status = models.StatusDONE
	r.logger.Info("work has done", "id", r.work.ID, "title", r.work.Title)

	return r.db.Done(r.work)
}

// Wait function block the program
// provide a simple shell to exec some command
// utils user enter the done string
func (r *WorkDefaultRunner) Wait() (done bool) {
	reader := bufio.NewReader(os.Stdin)
	prompt := "todo-work ->"
	fmt.Println("type help for help manual")
	fmt.Println(runnerHelp)

	exitWithoutDoneCmd := []string{"quit"}
	for {
		fmt.Print(prompt)
		cmd, err := reader.ReadString('\n')
		cmd = strings.Trim(cmd, " \n\t")
		if err != nil {
			fmt.Println("fail to read from input")
			continue
		}

		if strings.ToLower(cmd) == "done" {
			return true
		} else if slices.Contains(exitWithoutDoneCmd, cmd) {
			return false
		}

		r.execCmd(cmd)
	}
}

func (r *WorkDefaultRunner) execCmd(cmd string) {
	switch cmd {
	case "help":
		r.execHelp()
	case "view":
		r.execView()
	case "edit":
		r.execEdit()
	}
}

const runnerHelp = `==================== Help ====================
	view	open context file with viewer (read-only)
	edit	open context file with editor then edit it
	done	mark this work as done status then exit
	quit	exit without mark this work as done
	help	print help manual

	`

func (r *WorkDefaultRunner) execHelp() {
	fmt.Println(runnerHelp)
}

func (r *WorkDefaultRunner) execView() {
	err := r.viewer.PathView(r.work.ContextPath)
	if err != nil {
		r.logger.Error("while opening context file with viewer", "err", err)
		fmt.Println("error when open context with viewer: ", err)
	}
}

func (r *WorkDefaultRunner) execEdit() {
	err := r.viewer.PathEdit(r.work.ContextPath)
	if err != nil {
		r.logger.Error("while opening context file with viewer", "err", err)
		fmt.Println("error when open context with viewer: ", err)
	}
}

// getWorkViewer get the Viewer struct by the [WorkDefaultRunner.Work] viewerType
func getWorkViewer(t models.ViewerType) (ctxViewer viewer.ContextViewer) {
	switch t {
	case models.ViewerBatPrint:
		return viewer.NewBatPrintViewer()
	case models.ViewerPlainPrint:
		return viewer.NewPlainPrintViewer()
	case models.ViewerEditor:
		return viewer.NewEditorViewer(config.GlobalConf.EditorViewCommands)
	}

	return nil
}
