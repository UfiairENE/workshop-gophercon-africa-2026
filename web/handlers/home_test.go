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
	"workshop/internal/domain/body"
	"workshop/web/handlers"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		year, dist := 1995, 15.46
		got := serveHome(t, fakeBodies{items: []body.Body{
			{Name: "Earth", Kind: body.KindPlanet, System: body.SolarSystem},
			{Name: "51 Peg b", Kind: body.KindExoplanet, System: "51 Peg", Discovery: body.Discovery{Year: &year, By: "OHP"}, Sky: body.Sky{DistanceParsecs: &dist}},
		}})
		assert.Equal(t, http.StatusOK, got.Code)
		assert.Equal(t, "text/html; charset=utf-8", got.Header().Get("Content-Type"))
		assert.Contains(t, got.Body.String(), "51 Peg b")
	})

	t.Run("No bodies", func(t *testing.T) {
		t.Parallel()

		got := serveHome(t, fakeBodies{})
		assert.Equal(t, http.StatusOK, got.Code)
		assert.Contains(t, got.Body.String(), "bodies seed")
	})

	t.Run("List error", func(t *testing.T) {
		t.Parallel()

		got := serveHome(t, fakeBodies{err: errors.New("store unavailable")})
		assert.Equal(t, http.StatusInternalServerError, got.Code)
	})
}

type fakeBodies struct {
	items []body.Body
	err   error
}

func (f fakeBodies) List(context.Context, body.ListFilter) ([]body.Body, error) {
	return f.items, f.err
}

func serveHome(t *testing.T, lister handlers.BodiesLister) *httptest.ResponseRecorder {
	t.Helper()

	logger := logs.NewProduction(io.Discard, slog.LevelError)
	rec := httptest.NewRecorder()
	handlers.Home(logger, lister)(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))

	return rec
}
