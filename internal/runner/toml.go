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
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/viewer"
)

// WorkTomlRunner is used to exec the work display command and monitor work status
// When this work' status changed, it write log into file
type WorkTomlRunner struct {
	DataDir string

	Work *models.Work

	Viewer viewer.ContextViewer

	Logger *slog.Logger
}

// NewWorkTomlRunner create a new work runner based on the data dir and work pointer
func NewWorkTomlRunner(dataDir string, work *models.Work, log io.Writer) *WorkTomlRunner {
	handler := slog.NewTextHandler(log, nil)
	logger := slog.New(handler)

	vi := getWorkViewer(work.Viewer)

	return &WorkTomlRunner{
		DataDir: dataDir,
		Work:    work,
		Viewer:  vi,
		Logger:  logger,
	}
}

// Run starts the specified command and waits for it to complete.
func (r *WorkTomlRunner) Run(db store.TodoDB) error {
	r.Work.StartTime = time.Now()
	r.Work.Status = models.StatusDOING
	r.Logger.Info("starting work", "id", r.Work.ID, "title", r.Work.Title)

	done := r.Wait()
	if !done {
		r.Logger.Warn("exit work without done", "id", r.Work.ID, "title", r.Work.Title)
		return nil
	}

	r.Work.EndTime = time.Now()
	r.Work.Status = models.StatusDONE
	r.Logger.Info("work has done", "id", r.Work.ID, "title", r.Work.Title)

	return db.Done(r.Work)
}

// Wait function block the program
// provide a simple shell to exec some command
// utils user enter the done string
func (r *WorkTomlRunner) Wait() (done bool) {
	reader := bufio.NewReader(os.Stdin)
	prompt := "todo-work ->"
	fmt.Println("type help for help manual")

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

func (r *WorkTomlRunner) execCmd(cmd string) {
	switch cmd {
	case "help":
		r.execHelp()
	case "view":
		r.execView()
	}
}

const runnerHelp = `===== Help =====
	view	open context file with viewer
	done	mark this work as done status then exit
	quit	exit without mark this work as done
	help	print help manual

	`

func (r *WorkTomlRunner) execHelp() {
	fmt.Println(runnerHelp)
}

func (r *WorkTomlRunner) execView() {
	err := r.Viewer.PathView(r.Work.ContextPath)
	if err != nil {
		r.Logger.Error("while opening context file with viewer", "err", err)
		fmt.Println("error when open context with viewer: ", err)
	}
}

// getWorkViewer get the Viewer struct by the [WorkTomlRunner.Work] viewerType
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
