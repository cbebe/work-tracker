package worktracker

import (
	"sort"
	"time"
)

type Line struct {
	Start    time.Time
	Type     string
	Message  string
	Duration time.Duration
}

func StartDate(l1, l2 *Line) bool {
	return l1.Start.Before(l2.Start)
}

type By func(l1, l2 *Line) bool

func (by By) Sort(lines []Line) {
	ls := &lineSorter{
		lines: lines,
		by:    by,
	}
	sort.Sort(ls)
}

type lineSorter struct {
	lines []Line
	by    func(l1, l2 *Line) bool
}

func (s *lineSorter) Swap(i, j int) {
	s.lines[i], s.lines[j] = s.lines[j], s.lines[i]
}

func (s *lineSorter) Len() int {
	return len(s.lines)
}

func (s *lineSorter) Less(i, j int) bool {
	return s.by(&s.lines[i], &s.lines[j])
}
