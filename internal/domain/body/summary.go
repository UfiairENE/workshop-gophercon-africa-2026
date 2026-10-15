package body

import (
	"cmp"
	"slices"
)

type (
	// Summary is the headline figures over a set of bodies, worked out
	// once so anything reporting on them can lead with it.
	Summary struct {
		Bodies     int
		Exoplanets int
		Systems    int
		Facilities int
		Methods    int
		FirstYear  int
		LatestYear int
		// NearestParsecs is the distance to the closest system outside
		// our own, or zero if none has a distance.
		NearestParsecs float64
	}
	// Tally is one row of a leaderboard: a name and how often it appeared.
	Tally struct {
		Name  string
		Count int
	}
)

// Summarise reduces bodies to the headline figures. It reads every body
// once, so callers should hold on to the result.
func Summarise(bodies []Body) Summary {
	s := Summary{Bodies: len(bodies)}
	systems := map[string]struct{}{}
	facilities := map[string]struct{}{}
	methods := map[string]struct{}{}

	for _, b := range bodies {
		systems[b.System] = struct{}{}
		if b.Kind != KindExoplanet {
			continue
		}
		s.Exoplanets++
		if b.Discovery.By != "" {
			facilities[b.Discovery.By] = struct{}{}
		}
		if b.Discovery.Method != "" {
			methods[b.Discovery.Method] = struct{}{}
		}
		if y := b.Discovery.Year; y != nil {
			if s.FirstYear == 0 || *y < s.FirstYear {
				s.FirstYear = *y
			}
			s.LatestYear = max(s.LatestYear, *y)
		}
		if d := b.Sky.DistanceParsecs; d != nil && *d > 0 && (s.NearestParsecs == 0 || *d < s.NearestParsecs) {
			s.NearestParsecs = *d
		}
	}

	s.Systems = len(systems)
	s.Facilities = len(facilities)
	s.Methods = len(methods)
	return s
}

// CountBy tallies bodies by key, skipping those it returns "" for.
func CountBy(bodies []Body, key func(Body) string) map[string]int {
	counts := map[string]int{}
	for _, b := range bodies {
		if k := key(b); k != "" {
			counts[k]++
		}
	}
	return counts
}

// CountByYear tallies exoplanets by the year they were found.
func CountByYear(bodies []Body) map[int]int {
	counts := map[int]int{}
	for _, b := range bodies {
		if b.Kind == KindExoplanet && b.Discovery.Year != nil {
			counts[*b.Discovery.Year]++
		}
	}
	return counts
}

// Top ranks counts highest first, ties by name, keeping at most n.
func Top(counts map[string]int, n int) []Tally {
	rows := make([]Tally, 0, len(counts))
	for name, count := range counts {
		rows = append(rows, Tally{Name: name, Count: count})
	}
	slices.SortFunc(rows, func(a, b Tally) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), cmp.Compare(a.Name, b.Name))
	})
	return rows[:min(n, len(rows))]
}
