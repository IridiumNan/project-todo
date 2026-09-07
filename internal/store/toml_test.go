package store

import (
	"fmt"
	"log/slog"
	"math/rand"
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

func getTomlTestData(size int) (works []*models.Work) {
	for range size {
		works = append(works, createTestWork())
	}

	return
}

// TODO: Test with load the file
func TestTomlDump(t *testing.T) {
	allWorks := getTomlTestData(candidateSize)

	st := StoreToml{}

	defaultDumpTestPath := "/tmp/go-test-metadata.toml"
	err := st.DumpMetadataToFile(defaultDumpTestPath, allWorks)
	if err != nil {
		t.Error(err)
	}

	slog.Info("please check the data file", "file_path", defaultDumpTestPath)
}

func TestTomlLoad(t *testing.T) {
}

func TestTomlAppend(t *testing.T) {
}
