package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"uuid"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/web/views/components"
	"workshop/web/views/pages"
)

// SightingFinder is the part of the sighting service the detail page needs.
type SightingFinder interface {
	Find(ctx context.Context, id uuid.UUID) (sighting.Sighting, error)
}

func render(w http.ResponseWriter, r *http.Request, logger *slog.Logger, c interface {
	Render(ctx context.Context, w io.Writer) error
}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(r.Context(), w); err != nil {
		logger.ErrorContext(r.Context(), "handlers: rendering: "+err.Error())
	}
}

func all(logger *slog.Logger, w http.ResponseWriter, r *http.Request, s SightingsLister, f sighting.ListFilter) ([]sighting.Sighting, bool) {
	items, err := s.List(r.Context(), f)
	if err != nil {
		logger.ErrorContext(r.Context(), "handlers: listing sightings: "+err.Error())
		renderError(r.Context(), w, logger, http.StatusInternalServerError, "Unable to load sightings.")
		return nil, false
	}
	return items, true
}

func Sighting(logger *slog.Logger, finder SightingFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			NotFound(logger)(w, r)
			return
		}
		s, err := finder.Find(r.Context(), id)
		switch {
		case errors.Is(err, sighting.ErrNotFound):
			NotFound(logger)(w, r)
			return
		case err != nil:
			logger.ErrorContext(r.Context(), "handlers: finding sighting: "+err.Error())
			renderError(r.Context(), w, logger, http.StatusInternalServerError, "Unable to load this sighting.")
			return
		}
		render(w, r, logger, pages.Sighting(s))
	}
}

func Animal(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := strconv.Atoi(r.PathValue("key"))
		if err != nil {
			NotFound(logger)(w, r)
			return
		}
		everything, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		items := slices.DeleteFunc(slices.Clone(everything), func(s sighting.Sighting) bool { return s.Species.GBIFKey != key })
		if len(items) == 0 {
			NotFound(logger)(w, r)
			return
		}
		profiles := components.Profiles(everything)
		i := slices.IndexFunc(profiles, func(p components.Profile) bool { return p.Key == key })
		render(w, r, logger, pages.Animal(profiles[i], items, components.Related(profiles[i], profiles, 3)))
	}
}

func Country(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := country.Parse(r.PathValue("code"))
		if err != nil {
			NotFound(logger)(w, r)
			return
		}
		items, ok := all(logger, w, r, lister, sighting.ListFilter{Country: code})
		if !ok {
			return
		}
		if len(items) == 0 {
			NotFound(logger)(w, r)
			return
		}
		render(w, r, logger, pages.Listing("Countries", items[0].Location.Country, sighting.Summarise(items), items))
	}
}

func Stats(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}

		byYear := sighting.CountBy(items, func(s sighting.Sighting) string { return strconv.Itoa(s.HappenedAt.Year()) })
		years := make([]sighting.Tally, 0, len(byYear))
		for y, n := range byYear {
			years = append(years, sighting.Tally{Name: y, Count: n})
		}
		slices.SortFunc(years, func(a, b sighting.Tally) int { return strings.Compare(a.Name, b.Name) })

		animalKey := map[string]int{}
		animals := sighting.CountBy(items, func(s sighting.Sighting) string {
			animalKey[components.Name(s)] = s.Species.GBIFKey
			return components.Name(s)
		})
		countryCode := map[string]string{}
		countries := sighting.CountBy(items, func(s sighting.Sighting) string {
			countryCode[s.Location.Country] = s.Location.CountryCode.String()
			return s.Location.Country
		})

		var temp float64
		for _, s := range items {
			temp += s.Weather.Temperature.Actual
		}
		avg := 0.0
		if len(items) > 0 {
			avg = temp / float64(len(items))
		}

		data := pages.StatsData{
			Summary: sighting.Summarise(items),
			Years:   components.Bars(years, nil),
			Animals: components.Bars(sighting.Top(animals, 10), func(n string) string {
				return fmt.Sprintf("/animals/%d", animalKey[n])
			}),
			Countries: components.Bars(sighting.Top(countries, 10), func(n string) string {
				if c := countryCode[n]; c != "" {
					return "/countries/" + c
				}
				return ""
			}),
			Average: avg,
		}
		render(w, r, logger, pages.Stats(data))
	}
}

func MapPage(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		render(w, r, logger, pages.MapPage(items))
	}
}

func GalleryPage(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		withPhotos := items[:0:0]
		for _, s := range items {
			if _, has := components.Cover(s); has {
				withPhotos = append(withPhotos, s)
			}
		}
		if len(withPhotos) > 120 {
			withPhotos = withPhotos[:120]
		}
		render(w, r, logger, pages.GalleryPage(withPhotos))
	}
}

// NotFound renders the custom 404 page.
func NotFound(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		if err := pages.NotFound().Render(r.Context(), w); err != nil {
			logger.ErrorContext(r.Context(), "handlers: rendering 404: "+err.Error())
		}
	}
}

// Search lists sightings whose animal, place or recorder contains the query.
func Search(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	const perPage = 48
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		group := r.URL.Query().Get("group")
		if _, known := components.GroupByKey(group); !known {
			group = ""
		}
		groups := components.MapGroups(items)
		needle := strings.ToLower(q)
		found := slices.DeleteFunc(items, func(s sighting.Sighting) bool {
			if group != "" && components.GroupOf(s) != group {
				return true
			}
			hay := strings.ToLower(components.Name(s) + " " + s.Species.ScientificName + " " + components.Place(s) + " " + s.RecordedBy)
			return !strings.Contains(hay, needle)
		})
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		page = max(page, 1)
		total := len(found)
		from := min((page-1)*perPage, total)
		render(w, r, logger, pages.Sightings(pages.SightingsData{
			Query: q, Group: group, Groups: groups, Total: total,
			Items: found[from:min(from+perPage, total)],
			Page:  page, HasMore: from+perPage < total,
		}))
	}
}

// Group lists every species in one animal group.
func Group(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if _, known := components.GroupByKey(key); !known {
			NotFound(logger)(w, r)
			return
		}
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		var g components.MapGroup
		var others []components.MapGroup
		for _, x := range components.MapGroups(items) {
			if x.Key == key {
				g = x
			} else {
				others = append(others, x)
			}
		}
		if g.Key == "" {
			NotFound(logger)(w, r)
			return
		}
		species := slices.DeleteFunc(components.Profiles(items), func(p components.Profile) bool { return p.Group.Key != key })
		render(w, r, logger, pages.Group(g, species, others))
	}
}

// Quiz serves the quizzes: every species with a photo, with its call and fact when known.
func Quiz(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		var out []pages.QuizAnimal
		sounds, facts := map[string]bool{}, 0
		for _, p := range components.Profiles(items) {
			if p.Photo == "" {
				continue
			}
			a := pages.QuizAnimal{ID: p.Key, Name: p.Name, Photo: p.Photo, Group: p.Group.Key, Genus: p.Genus, Fact: p.Fact, Href: p.Href()}
			// Species that share a recording would make two right answers, so only the first keeps it.
			if p.Sound != "" && !sounds[p.Sound] {
				sounds[p.Sound] = true
				a.Sound = p.Sound
			}
			if p.Fact != "" {
				facts++
			}
			out = append(out, a)
		}
		render(w, r, logger, pages.Quiz(out, len(sounds), facts))
	}
}

// Near lists the sightings' positions for the in-browser "near me" ranking.
func Near(logger *slog.Logger, lister SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, ok := all(logger, w, r, lister, sighting.ListFilter{})
		if !ok {
			return
		}
		render(w, r, logger, pages.Near(components.MapPoints(items)))
	}
}

// Credits lists the openly licensed media used by the site.
func Credits(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, logger, pages.Credits())
	}
}
