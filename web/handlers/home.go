package handlers

import (
	"cmp"
	"context"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"workshop/internal/domain/sighting"
	"workshop/web/views/components"
	"workshop/web/views/pages"
)

type SightingsLister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
}

func Home(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, sighting.ListFilter{})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		var photos []sighting.Sighting
		for _, s := range items {
			if _, ok := components.Cover(s); ok {
				photos = append(photos, s)
			}
		}
		// The hero slideshow: one animal per kind, then any other species.
		var heroes []sighting.Sighting
		used := map[int]bool{}
		seenKind := map[string]bool{}
		for _, s := range photos {
			k := components.KindOf(s)
			if len(heroes) < 5 && !seenKind[k] {
				seenKind[k] = true
				used[s.Species.GBIFKey] = true
				heroes = append(heroes, s)
			}
		}
		for _, s := range photos {
			if len(heroes) < 5 && !used[s.Species.GBIFKey] {
				used[s.Species.GBIFKey] = true
				heroes = append(heroes, s)
			}
		}

		// Below the hero, each species appears at most once on the page.
		profiles := components.Profiles(items)
		var groups []pages.GroupTile
		for _, g := range components.MapGroups(items) {
			t := pages.GroupTile{Group: g}
			for _, p := range profiles {
				if p.Group.Key != g.Key {
					continue
				}
				t.Species++
				if t.Photo == "" && p.Photo != "" && !used[p.Key] {
					t.Photo, t.Cover = p.Photo, p.Name
					used[p.Key] = true
				}
			}
			groups = append(groups, t)
		}

		var spotlight *components.Profile
		var candidates []components.Profile
		for _, p := range profiles {
			if p.Photo != "" && p.Fact != "" && !used[p.Key] {
				candidates = append(candidates, p)
			}
		}
		if len(candidates) > 0 {
			p := candidates[time.Now().YearDay()%len(candidates)]
			used[p.Key] = true
			spotlight = &p
		}

		var latest []sighting.Sighting
		for _, s := range photos {
			if len(latest) < 4 && !used[s.Species.GBIFKey] {
				used[s.Species.GBIFKey] = true
				latest = append(latest, s)
			}
		}

		sounds := map[string]bool{}
		for _, p := range profiles {
			if p.Sound != "" && p.Photo != "" {
				sounds[p.Sound] = true
			}
		}

		d := pages.HomeData{
			Summary:   sighting.Summarise(items),
			Heroes:    heroes,
			Groups:    groups,
			Spotlight: spotlight,
			Quiz:      len(sounds),
			Countries: topCountries(items, 8),
			Latest:    latest,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Home(d).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering sightings: "+err.Error())
		}
	}
}

// topCountries ranks countries by how many sightings they have.
func topCountries(items []sighting.Sighting, n int) []pages.CountryTally {
	byCode := map[string]*pages.CountryTally{}
	species := map[string]map[int]bool{}
	for _, s := range items {
		code := s.Location.CountryCode.String()
		if code == "" || s.Location.Country == "" {
			continue
		}
		t, ok := byCode[code]
		if !ok {
			t = &pages.CountryTally{Name: s.Location.Country, Code: code}
			byCode[code] = t
			species[code] = map[int]bool{}
		}
		t.Sightings++
		species[code][s.Species.GBIFKey] = true
	}
	out := make([]pages.CountryTally, 0, len(byCode))
	for code, t := range byCode {
		t.Species = len(species[code])
		out = append(out, *t)
	}
	slices.SortFunc(out, func(a, b pages.CountryTally) int {
		return cmp.Or(cmp.Compare(b.Sightings, a.Sightings), cmp.Compare(a.Name, b.Name))
	})
	return out[:min(n, len(out))]
}
