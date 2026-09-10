package store

import (
	"fmt"
	"math/rand"
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

// func initTestDir(t *testing.T) {
// 	dataDirPath, err := InitDataDirOnCurrentDir()
// 	if err != nil {
// 		t.Errorf("error when init data dir, err: %s", err)
// 	}
//
// }

// TODO: Write test for tomlDB
