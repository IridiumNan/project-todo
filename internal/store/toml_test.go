package store

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
)

const candidateSize = 5

var titleCandidates = []string{
	"hello world",
	"hello linux",
	"hello golang",
	"hello kernel",
	"PR 43, writing test",
}

func createTestWork() *models.Work {
	id := utils.HashByTimestamp(time.Now().UnixNano(), 8)
	idx := rand.Intn(candidateSize)
	return &models.Work{
		ID:                id,
		Title:             titleCandidates[idx],
		EnergyRequirement: models.EnergyHigh,
		Status:            models.StatusTODO,
		// TODO: Use real auto generated path
		ContextPath:    fmt.Sprintf("/tmp/fake_path_%s", id),
		Viewer:         models.ViewerEditor,
		CreateTime:     time.Now(),
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		BlockedTimes:   idx,
		BlockedWorksID: []string{},
	}
}

var tomlTestWorks = getTomlTestData(candidateSize)

// isSameWork compare filed except Time
func isSameWork(want *models.Work, got *models.Work) bool {
	res := want.ID == got.ID && want.Title == got.Title && want.EnergyRequirement == got.EnergyRequirement && want.Status == got.Status && want.ContextPath == got.ContextPath && want.BlockedTimes == got.BlockedTimes

	return res
}

func getTomlTestData(size int) (works []*models.Work) {
	for range size {
		works = append(works, createTestWork())
	}

	return
}

func initTestDir(t *testing.T) (dataDirPath string) {
	var err error
	dataDirPath, err = InitDataDirOnCurrentDir()
	if err != nil {
		t.Errorf("error when init data dir, err: %s", err)
	}

	slog.Info("init data dir", "dir_path", dataDirPath)

	return
}

func removeTestDir(dataDirPath string) error {
	if err := os.RemoveAll(dataDirPath); err != nil {
		return err
	}

	return nil
}

func TestTomlPush(t *testing.T) {
	dataDirPath := initTestDir(t)
	// defer removeTestDir(dataDirPath)

	td, err := NewTomlDB(dataDirPath)
	if err != nil {
		t.Errorf("error when create new toml db, err: %s", err)
	}

	conf := models.MDTomlConfig{
		Title:          "Test 0",
		Energy:         models.EnergyLow,
		Viewer:         models.ViewerBatPrint,
		DependenciesID: []string{},
	}

	ctx := []byte("Test 0 context")

	_, err = td.Push(&conf, ctx)
	if err != nil {
		t.Errorf("error when push new work, err: %s", err)
	}

	err = td.Sync()
	if err != nil {
		t.Errorf("error when sync todo works into disk, err: %s", err)
	}

	oldWork, err := td.Pop(filter.EnergyFilter(models.EnergyLow))
	slog.Info("first pop finished")
	if err != nil {
		t.Errorf("error when pop work, err: %s", err)
	}

	newTd, err := NewTomlDB(dataDirPath)
	slog.Info("second td built")
	if err != nil {
		t.Errorf("error when create new toml db, err: %s", err)
	}
	work, err := newTd.Pop(filter.EnergyFilter(models.EnergyLow))
	if err != nil {
		t.Errorf("error when pop work, err: %s", err)
	}

	if !isSameWork(oldWork, work) {
		t.Errorf("load same work failed, want: %v, got: %v", oldWork, work)
	}
}
