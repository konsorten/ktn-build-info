package ver

import (
	"sort"
)

type ItemDocumentationEntry struct {
	Name        string
	Description string
}

type ItemDocumentation []ItemDocumentationEntry

func (s ItemDocumentation) Len() int {
	return len(s)
}

func (s ItemDocumentation) Less(i, j int) bool {
	a := &s[i]
	b := &s[j]

	return a.Name < b.Name
}

func (s ItemDocumentation) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ItemDocumentation) Sort() {
	sort.Sort(s)
}
