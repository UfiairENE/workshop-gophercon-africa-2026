package components

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"workshop/internal/domain/sighting"
)

// Profile is everything the site knows about one species, gathered from its sightings.
type Profile struct {
	Key                               int
	Name, Scientific                  string
	Class, Order, Family, Genus       string
	Group                             MapGroup
	Photo, Fact, Sound                string
	Sightings, Countries, Individuals int
	First, Last                       time.Time
}

// Href is the species' own page.
func (p Profile) Href() string { return fmt.Sprintf("/animals/%d", p.Key) }

// Taxonomy is the species' place in the tree of life, broadest first.
func (p Profile) Taxonomy() []string {
	return slices.DeleteFunc([]string{p.Class, p.Order, p.Family, p.Genus}, func(s string) bool { return s == "" })
}

// Profiles builds one profile per species, most-seen first.
func Profiles(items []sighting.Sighting) []Profile {
	byKey := map[int]*Profile{}
	countries := map[int]map[string]bool{}
	var order []int
	for _, s := range items {
		k := s.Species.GBIFKey
		p, ok := byKey[k]
		if !ok {
			sp := s.Species
			p = &Profile{
				Key: k, Name: Name(s), Scientific: sp.CanonicalName,
				Class: sp.Class, Order: sp.Order, Family: sp.Family, Genus: sp.Genus,
				Group: groupByKey(GroupOf(s)), Fact: Fact(s), Sound: SoundFor(s),
				First: s.HappenedAt, Last: s.HappenedAt,
			}
			byKey[k] = p
			countries[k] = map[string]bool{}
			order = append(order, k)
		}
		p.Sightings++
		if s.IndividualCount != nil {
			p.Individuals += *s.IndividualCount
		}
		if s.Location.Country != "" {
			countries[k][s.Location.Country] = true
		}
		if s.HappenedAt.Before(p.First) {
			p.First = s.HappenedAt
		}
		if s.HappenedAt.After(p.Last) {
			p.Last = s.HappenedAt
		}
		if p.Photo == "" {
			p.Photo, _ = Cover(s)
		}
	}
	out := make([]Profile, 0, len(order))
	for _, k := range order {
		p := byKey[k]
		p.Countries = len(countries[k])
		out = append(out, *p)
	}
	slices.SortStableFunc(out, func(a, b Profile) int { return cmp.Compare(b.Sightings, a.Sightings) })
	return out
}

func groupByKey(key string) MapGroup {
	for _, g := range mapGroups {
		if g.Key == key {
			return g
		}
	}
	return mapGroups[len(mapGroups)-1]
}

// GroupByKey finds a group by its key.
func GroupByKey(key string) (MapGroup, bool) {
	for _, g := range mapGroups {
		if g.Key == key {
			return g, true
		}
	}
	return MapGroup{}, false
}

// Related picks up to n other species to explore next: the same family first,
// then the same group.
func Related(p Profile, all []Profile, n int) []Profile {
	var out []Profile
	add := func(match func(Profile) bool) {
		for _, o := range all {
			if len(out) < n && o.Key != p.Key && o.Photo != "" && match(o) && !slices.ContainsFunc(out, func(x Profile) bool { return x.Key == o.Key }) {
				out = append(out, o)
			}
		}
	}
	add(func(o Profile) bool { return o.Family == p.Family })
	add(func(o Profile) bool { return o.Group.Key == p.Group.Key })
	return out
}
