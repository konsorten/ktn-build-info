package ver

import (
	"sort"
	"strings"
)

type orderedUpdate struct {
	Path   string
	Update string
}

func (s *orderedUpdate) weight() int {
	// no action?
	if !strings.HasPrefix(s.Update, "$") {
		return 200
	}

	// delete actions come first
	if strings.HasPrefix(s.Update, "$delete") {
		return 110
	}

	// create actions
	if strings.HasPrefix(s.Update, "$create") {
		return 120
	}

	// ensure actions
	if strings.HasPrefix(s.Update, "$ensure") {
		return 130
	}

	// ordinary action
	return 199
}

type orderedUpdates []orderedUpdate

func (s orderedUpdates) Len() int {
	return len(s)
}

func (s orderedUpdates) Less(i, j int) bool {
	a := &s[i]
	b := &s[j]

	wa := a.weight()
	wb := b.weight()

	if wa != wb {
		return wa < wb
	}

	return a.Path < b.Path
}

func (s orderedUpdates) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s orderedUpdates) Sort() {
	sort.Sort(s)
}

func orderUpdates(updates UpdateActions) orderedUpdates {
	var ret = make(orderedUpdates, len(updates))
	index := 0

	for path, update := range updates {
		ret[index] = orderedUpdate{path, update}

		index++
	}

	// sort the entries
	ret.Sort()

	return ret
}
