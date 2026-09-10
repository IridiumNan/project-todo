package utils

import (
	"slices"
	"testing"
)

func TestMapValues(t *testing.T) {
	type testPerson struct {
		ID string

		Name string

		Age int
	}
	persons := []testPerson{
		{
			"1",
			"John",
			18,
		},
		{
			"2",
			"Alice",
			19,
		},
		{
			"3",
			"Mike",
			20,
		},
	}

	personMap := make(map[string]testPerson, 3)
	for _, p := range persons {
		personMap[p.ID] = p
	}

	values := MapValues(personMap)

	slices.SortFunc(values, func(p1, p2 testPerson) int {
		if p1.ID == p2.ID {
			return 0
		}
		if p1.ID > p2.ID {
			return 1
		}

		return -1
	})
	if !slices.Equal(persons, values) {
		t.Errorf("fail to convert map values into slices, want: %v, got: %v", persons, values)
	}
}
