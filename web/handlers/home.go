package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"workshop/internal/domain/body"
	"workshop/web/views/pages"
)

type BodiesLister interface {
	List(ctx context.Context, filter body.ListFilter) ([]body.Body, error)
}

func Home(logger *slog.Logger, bodies BodiesLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := bodies.List(ctx, body.ListFilter{})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing bodies: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load the catalogue.")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Home(items).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering bodies: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
