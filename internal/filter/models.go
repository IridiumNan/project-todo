package filter

import "github.com/IridiumNan/project-todo/internal/models"

// WorkFilter for filter valid work for current condition
// There is no need to check if work.status is StatusTODO
// The [WorkProvider.Provide] function will skip invalid status
type WorkFilter func(*models.Work) bool

func EnergyFilter(e models.Energy) WorkFilter {
	return func(w *models.Work) bool {
		return w.BlockedTimes == 0 && w.EnergyRequirement == e
	}
}

type WorkProvider struct {
	AllWorks []*models.Work
}

func (wp *WorkProvider) Provide(filter WorkFilter) *models.Work {
	for _, work := range wp.AllWorks {
		// Check if status is todo
		// if not skip
		// Check on this provide loop because any work provided should be status todo
		if work.Status != models.StatusTODO {
			continue
		}

		if filter(work) {
			return work
		}
	}
	return nil
}

func MultiFiler(filters ...WorkFilter) WorkFilter {
	return func(w *models.Work) bool {
		for _, f := range filters {
			if !f(w) {
				return false
			}
		}

		return true
	}
}
