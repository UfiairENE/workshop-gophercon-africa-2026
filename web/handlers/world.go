package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"workshop/internal/domain/body"
	"workshop/web/views/pages"
)

type WorldFinder interface {
	FindByName(ctx context.Context, name string) (body.Body, error)
	List(ctx context.Context, filter body.ListFilter) ([]body.Body, error)
}

func World(logger *slog.Logger, worlds WorldFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := r.PathValue("name")

		b, err := worlds.FindByName(ctx, name)
		if errors.Is(err, body.ErrNotFound) {
			renderError(ctx, w, logger, http.StatusNotFound, "No world called "+name+" is in the catalogue.")
			return
		} else if err != nil {
			logger.ErrorContext(ctx, "handlers: finding world: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load this world.")
			return
		}

		family, err := worlds.List(ctx, body.ListFilter{System: b.System})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing system: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load this world.")
			return
		}
		host, family := split(b, family)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.World(b, host, family).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering world: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}

// split picks out the star b belongs to, and the bodies worth listing
// beside it: a star's planets, a planet's moons or siblings, a moon's
// fellow moons.
func split(b body.Body, system []body.Body) (host body.Body, family []body.Body) {
	for _, other := range system {
		if other.Kind == body.KindStar {
			host = other
		}
	}

	family = slices.DeleteFunc(system, func(o body.Body) bool {
		if o.Name == b.Name || o.Kind == body.KindStar {
			return true
		}
		switch b.Kind {
		case body.KindStar:
			return o.Parent != b.Name
		case body.KindMoon:
			return o.Parent != b.Parent
		case body.KindPlanet, body.KindDwarfPlanet:
			return o.Parent != b.Name && o.Parent != "Sun"
		default:
			return false
		}
	})
	return host, family
}
