package store

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

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

var defaultTomlTestpath = "/tmp/go-test-metadata.toml"

var tomlTestWorks = getTomlTestData(candidateSize)

var tomlAppendWorks = getTomlTestData(2)

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

func removeTomlTestFile(t *testing.T) {
	if err := os.Remove(defaultTomlTestpath); err != nil {
		t.Errorf("error when remove the defaultTomlTestpath: %s, err: %s", defaultTomlTestpath, err.Error())
	}
}

func dumpToml(t *testing.T) {
	st := StoreToml{}

	err := st.DumpMetadataToFile(defaultTomlTestpath, tomlTestWorks)
	if err != nil {
		t.Error(err)
	}
}

func loadToml() (works []*models.Work, err error) {
	st := StoreToml{}

	tomlSrc, err := os.OpenFile(defaultTomlTestpath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error when opening default toml test file, file path: %s, err: %s", defaultTomlTestpath, err.Error())
	}

	works, err = st.LoadMetadata(tomlSrc)
	if err != nil {
		slog.Error("when load metadata", "err", err)
	}

	return
}

func TestTomlAppend(t *testing.T) {
}

func TestTomlDumpLoad(t *testing.T) {
	removeTomlTestFile(t)
	dumpToml(t)

	works, err := loadToml()
	if err != nil {
		t.Errorf("error when load work from defaultTomlTestpath: %s, err: %s", defaultTomlTestpath, err.Error())
	}
	for idx := range works {
		if !isSameWork(tomlTestWorks[idx], works[idx]) {
			t.Errorf("want %v, got %v", tomlTestWorks[idx], works[idx])
		}
	}
}
