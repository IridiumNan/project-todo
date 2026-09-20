package filter

import (
	"github.com/IridiumNan/project-todo/internal/models"
)

type FilterAssistant struct{}

func NewFilterAssistant() *FilterAssistant {
	return &FilterAssistant{}
}

// SliceToSlice return the slice contains all works that make f(work) == true
func (fa *FilterAssistant) SliceToSlice(works []*models.Work, f WorkFilter) ([]*models.Work, error) {
	if f == nil {
		f = DefaultFilter()
	}

	passedWorks := make([]*models.Work, 0, len(works))
	for _, w := range works {
		if f(w) {
			passedWorks = append(passedWorks, w)
		}
	}

	return passedWorks, nil
}

// sliceToMapNoFilter convert from slice into map
// key is [models.Work.ID] and value is pointer of [models.Work]
func sliceToMapNoFilter(works []*models.Work) map[string]*models.Work {
	m := make(map[string]*models.Work, len(works))
	for _, w := range works {
		m[w.ID] = w
	}

	return m
}

// SliceToMap return the slice contains all works that make f(work) == true
// key is [models.Work.ID] and value is pointer of [models.Work]
func (fa *FilterAssistant) SliceToMap(works []*models.Work, f WorkFilter) (map[string]*models.Work, error) {
	if f == nil {
		f = DefaultFilter()
	}

	passedWorkMap := make(map[string]*models.Work, len(works))

	for _, w := range works {
		if f(w) {
			passedWorkMap[w.ID] = w
		}
	}

	return passedWorkMap, nil
}

// MapToSlice return the slices contains all works that make f(work) == true
func (fa *FilterAssistant) MapToSlice(workMap map[string]*models.Work, f WorkFilter) ([]*models.Work, error) {
	if f == nil {
		f = DefaultFilter()
	}

	passedWorks := make([]*models.Work, 0, len(workMap))
	for _, w := range workMap {
		if f(w) {
			passedWorks = append(passedWorks, w)
		}
	}

	return passedWorks, nil
}

// MapToMap return the map contains all works that make f(work) == true
// key is [models.Work.ID] and value is pointer of [models.Work]
func (fa *FilterAssistant) MapToMap(workMap map[string]*models.Work, f WorkFilter) (map[string]*models.Work, error) {
	if f == nil {
		f = DefaultFilter()
	}

	passedWorkMap := make(map[string]*models.Work, len(workMap))

	for _, w := range workMap {
		if f(w) {
			passedWorkMap[w.ID] = w
		}
	}

	return passedWorkMap, nil
}
