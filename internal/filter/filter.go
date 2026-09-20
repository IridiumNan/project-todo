package filter

import (
	"github.com/IridiumNan/project-todo/internal/models"
)

// WorkFilter for filter valid work for current condition
// There is no need to check if work.status is StatusTODO
// If filter return true, it means this work passed filter, it should appear on result
type WorkFilter func(*models.Work) bool

func EnergyFilter(e models.Energy) WorkFilter {
	if e == models.EnergyAll {
		return nil
	}
	return func(w *models.Work) bool {
		return w.BlockedTimes == 0 && w.EnergyRequirement == e
	}
}

func MultiFilter(filters ...WorkFilter) WorkFilter {
	return func(w *models.Work) bool {
		for _, f := range filters {
			if !f(w) {
				return false
			}
		}

		return true
	}
}

func DefaultFilter() WorkFilter {
	return func(w *models.Work) bool {
		return w.BlockedTimes == 0
	}
}
