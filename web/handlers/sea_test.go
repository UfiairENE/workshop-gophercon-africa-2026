package handlers_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"workshop/internal/common/logs"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/web/handlers"

	"github.com/stretchr/testify/assert"
)

func TestSea(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		got := serveSea(t, fakeLister{items: []sighting.Sighting{
			{GBIFKey: 1, Species: species.Species{VernacularName: "lion"}},
		}})
		assert.Equal(t, http.StatusOK, got.Code)
		assert.Equal(t, "text/html; charset=utf-8", got.Header().Get("Content-Type"))
	})

	t.Run("No sightings", func(t *testing.T) {
		t.Parallel()

		got := serveSea(t, fakeLister{})
		assert.Equal(t, http.StatusOK, got.Code)
	})

	t.Run("List error", func(t *testing.T) {
		t.Parallel()

		got := serveSea(t, fakeLister{err: errors.New("store unavailable")})
		assert.Equal(t, http.StatusInternalServerError, got.Code)
		assert.Equal(t, "text/html; charset=utf-8", got.Header().Get("Content-Type"))
	})
}

type fakeLister struct {
	items []sighting.Sighting
	err   error
}

func (f fakeLister) List(context.Context, sighting.ListFilter) ([]sighting.Sighting, error) {
	return f.items, f.err
}

func serveSea(t *testing.T, lister handlers.SightingsLister) *httptest.ResponseRecorder {
	t.Helper()

	logger := logs.NewProduction(io.Discard, slog.LevelError)
	rec := httptest.NewRecorder()
	handlers.Sea(logger, lister)(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/sea", nil))

	return rec
}
