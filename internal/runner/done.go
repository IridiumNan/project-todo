package runner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	_ "embed"

	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/IridiumNan/project-todo/internal/viewer"
)

// WorkDoneRunner provide read-only mode for works whose status is [models.StatusDONE]
// It construct a new tmp file about specific done work with metadata and duration of this work
type WorkDoneRunner struct {
	dataDir string

	logger *slog.Logger

	db store.DoneDB

	f filter.WorkFilter

	work *models.Work

	// summaryPath is the summary for work file path (a markdown path now)
	tmpFilePath string
}

func NewDoneWorkRunner(dir string, log io.Writer, d store.DoneDB, f filter.WorkFilter) *WorkDoneRunner {
	handler := slog.NewTextHandler(log, nil)
	logger := slog.New(handler)

	return &WorkDoneRunner{
		dataDir: dir,
		logger:  logger,
		db:      d,
		f:       f,
	}
}

//go:embed done_work.md
var doneWorkTemplate []byte

const (
	titleReplaceHolder       = "{{Title}}"
	currentTimeReplaceHolder = "{{current_time}}"
	contextReplaceHolder     = "{{Context}}"
	commentHintReplaceHolder = "{{comment_hint}}"
	infoReplaceHolder        = "{{info}}"
)

// buildWorkSummary create a new markdown file for this work then return the path
func (r *WorkDoneRunner) buildWorkSummary() (filePath string, err error) {
	// TODO:
	// Use the r.work
	// inject then create new tmp markdown file on dataDir/build/ then return the file path
	// Build this file on build dir

	// Write markdown build process there now
	// FIX: Rewrite the builder with a markdown builder

	byteData := doneWorkTemplate

	contextByte, err := os.ReadFile(r.work.ContextPath)
	if err != nil {
		err = fmt.Errorf("error when read the context file, file path: %s, err: %s", r.work.ContextPath, err.Error())
		return
	}

	byteData = bytes.Replace(byteData, []byte(titleReplaceHolder), []byte(r.work.Title), 1)
	byteData = bytes.Replace(byteData, []byte(currentTimeReplaceHolder), []byte(time.Now().String()), 1)
	byteData = bytes.Replace(byteData, []byte(contextReplaceHolder), contextByte, 1)
	byteData = bytes.Replace(byteData, []byte(commentHintReplaceHolder), []byte("You can add some comment or append information here or anywhere you want"), 1)
	byteData = bytes.Replace(byteData, []byte(infoReplaceHolder), []byte("Time and id information replacer"), 1)

	return r.dumpWorkSummary(byteData)
}

func (r *WorkDoneRunner) dumpWorkSummary(byteData []byte) (filePath string, err error) {
	summaryDir := path.Join(r.dataDir, models.DataSummaryDirName)

	tmpFile, err := os.CreateTemp(summaryDir, "project-todo-summary-*.md")
	if err != nil {
		err = fmt.Errorf("error while creating a tmp file, err: %s", err.Error())
		return
	}

	_, err = tmpFile.Write(byteData)
	if err != nil {
		err = fmt.Errorf("error when write byteData into tmp file, file_path: %s, err: %s", tmpFile.Name(), err.Error())

		return
	}

	return tmpFile.Name(), nil
}

// selectWork get all matched works map then use query user select one from list (by fzf)
// Then return the work
// if failed, it will return [models.InvalidWork]
func (r *WorkDoneRunner) selectWork() (err error) {
	worksMap, err := r.db.AllWithMap(r.f)
	if err != nil {
		return fmt.Errorf("error while getting works from database, err: %s", err.Error())
	}

	works := utils.MapValues(worksMap)

	opts, err := r.buildFzfOptions(works)
	if err != nil {
		slog.Error("while building options for fzf select", "err", err)
		return fmt.Errorf("error while building options for fzf select, err: %s", err.Error())
	}

	in := bytes.NewBufferString(opts)

	selected, err := utils.SelectByFzf(in)
	if err != nil {
		slog.Error("while select by fzf", "err", err)
		return fmt.Errorf("error while select work by fzf, err: %s", err.Error())
	}

	// NOTE: The format of a single line is <id> <title> <entime>, just use the first part as the key
	strParts := strings.SplitN(selected, " ", 2)

	r.work = worksMap[strParts[0]]
	return nil
}

// buildFzfOptions build a string based on input works
// With format
// id title endtime (a single line)
func (r *WorkDoneRunner) buildFzfOptions(works []*models.Work) (optionsStr string, err error) {
	if len(works) == 0 {
		return models.EmptyStr, errors.New("buildFzfOptions: empty work slices")
	}

	for _, w := range works {
		line := fmt.Sprintf("%s %s %v\n", w.ID, w.Title, w.EndTime)

		optionsStr = optionsStr + line
	}

	return optionsStr, nil
}

func (r *WorkDoneRunner) isBuildTmpFile() bool {
	buildDir := path.Join(r.dataDir, models.DataBuildDirName)

	return path.Dir(r.tmpFilePath) == buildDir
}

// prepare function load the database works into memory
func (r *WorkDoneRunner) prepare() error {
	err := r.db.Load(r.dataDir)
	if err != nil {
		return fmt.Errorf("error while loading works metadata from dataDir: %s, err: %s", r.dataDir, err.Error())
	}

	return nil
}

// Run function query user for selecting a work then build a summary tmp file
// Then call [WorkRunner.Wait] function where the user can choose view edit summary, save it etc...
func (r *WorkDoneRunner) Run() error {
	err := r.prepare()
	if err != nil {
		return r.errWrap(err.Error())
	}

	r.logger.Info("starting done work runner")

	err = r.selectWork()
	if err != nil {
		return fmt.Errorf("while selecting work, err: %s", err.Error())
	}

	r.tmpFilePath, err = r.buildWorkSummary()
	if err != nil {
		slog.Error("while building the summary for work", "workId", r.work.ID, "workTitle", r.work.Title, "err", err)
		return fmt.Errorf("while building the summary for work, id: %s, title: %s, err: %w", r.work.ID, r.work.Title, err)
	}

	// start the wait shell
	done := r.Wait()

	if !done {
		return fmt.Errorf("not saved on default Path")
	}

	defer func() {
		// if tmp file not point into the build dir, skip clean
		// else remove the tmp file
		if !r.isBuildTmpFile() {
			return
		}

		os.Remove(r.tmpFilePath)
	}()

	return r.saveOnDefaultSummaryPath()
}

// saveOnDefaultSummaryPath save the summary file on default path
func (r *WorkDoneRunner) saveOnDefaultSummaryPath() error {
	fileName := fmt.Sprintf("%s-%s.md", r.work.ID, r.work.Title)

	defaultPath := path.Join(r.dataDir, models.DataSummaryDirName, fileName)

	r.logger.Info("save summary file on default path", "file_path", defaultPath)
	fmt.Printf("save summary file on default path, file_path: %s", defaultPath)

	return r.saveSummaryWithPath(defaultPath)
}

// Wait start a simple shell and user run command for view, edit or save summary file
// if return true, you save work on default summary path
// else just exit
func (r *WorkDoneRunner) Wait() (done bool) {
	reader := bufio.NewReader(os.Stdin)
	prompt := "done-work ->"
	fmt.Println("type help for help manual")
	fmt.Println(doneRunnerHelp)

	exitWithoutDoneCmd := []string{"quit"}

	for {
		fmt.Print(prompt)

		cmd, err := reader.ReadString('\n')
		cmd = strings.Trim(cmd, " \n\t")
		if err != nil {
			fmt.Println("fail to read from standard input")
			continue
		}
		cmd = strings.ToLower(cmd)
		if cmd == "exit" {
			return true
		}
		if slices.Contains(exitWithoutDoneCmd, cmd) {
			return false
		}

		r.execCmd(cmd)
	}
}

// saveSummaryWithPath just move the summary file from default path to dstPath
func (r *WorkDoneRunner) saveSummaryWithPath(dstPath string) error {
	err := os.Rename(r.tmpFilePath, dstPath)
	fmt.Printf("move summary file from %s to %s", r.tmpFilePath, dstPath)
	slog.Info("move summary file", "src_path", r.tmpFilePath, "dst_path", dstPath)
	if err != nil {
		return fmt.Errorf("error while save summary file to dst, dstPath: %s, err: %s", dstPath, err)
	}

	return nil
}

func (r *WorkDoneRunner) execCmd(cmd string) {
	// exec functions below just handle the error inner
	switch cmd {
	case "view":
		r.execView()
	case "edit":
		r.execEdit()
	case "save":
		r.execSave()
	case "help":
		r.execHelp()
	}
}

func (r *WorkDoneRunner) execView() {
	vi := viewer.GetWorkViewer(r.work.Viewer)

	err := vi.PathView(r.tmpFilePath)
	if err != nil {
		slog.Error("while viewing summary file", "path", r.tmpFilePath)
		fmt.Printf("error when viewing summary file with editor, path: %s, err: %s", r.tmpFilePath, err)
	}
}

func (r *WorkDoneRunner) execEdit() {
	vi := viewer.GetWorkViewer(r.work.Viewer)

	err := vi.PathEdit(r.tmpFilePath)
	if err != nil {
		slog.Error("while open summary file with editor", "path", r.tmpFilePath)
		fmt.Printf("error when opening summary file with editor, path: %s, err: %s", r.tmpFilePath, err)
	}
}

func (r *WorkDoneRunner) execSave() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("save file name (it will be save on current dir) ->")

	fileName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("error when read from standard input, err: %s", err.Error())
		return
	}
	fileName = strings.Trim(fileName, "\t\n ")

	// dstPath := path.Join(r.dataDir, models.DataSummaryDirName, fileName)
	dstPath := fileName

	// read the path then save
	err = r.saveSummaryWithPath(dstPath)
	if err != nil {
		slog.Error("while save summary file with new name", "savedPath", dstPath)
	}

	r.tmpFilePath = dstPath
}

const doneRunnerHelp = `==================== Help ====================
	view	open summary file with viewer (read-only)
	edit	open context file with editor then edit it
	save	save this summary file (or it will be saved as default path)
	exit	exit this program (with save on default summary dir
	quit	exit this program and not save the summary file
	help	print help manual`

func (r *WorkDoneRunner) execHelp() {
	fmt.Println(doneRunnerHelp)
}

func (r *WorkDoneRunner) errWrapf(format string, a ...any) error {
	errStr := fmt.Sprintf(format, a...)
	return fmt.Errorf("WorkDoneRunner error: %s", errStr)
}

func (r *WorkDoneRunner) errWrap(errStr string) error {
	return fmt.Errorf("WorkDoneRunner error: %s", errStr)
}
