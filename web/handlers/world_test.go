package handlers_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"workshop/internal/common/logs"
	"workshop/internal/domain/body"
	"workshop/web/handlers"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

var (
	epsEri  = body.Body{Name: "eps Eri", Kind: body.KindStar, System: "eps Eri", Star: body.Star{SpectralType: "K2 V"}, Sky: body.Sky{RightAscension: ptr(53.2), Declination: ptr(-9.46), ApparentMagnitude: ptr(3.73), DistanceParsecs: ptr(3.2)}}
	epsEriB = body.Body{Name: "eps Eri b", Kind: body.KindExoplanet, System: "eps Eri", Parent: "eps Eri", Discovery: body.Discovery{Year: ptr(2000), By: "McDonald Observatory", Method: "Radial Velocity"}}
)

type fakeWorlds struct{ bodies []body.Body }

func (f fakeWorlds) FindByName(_ context.Context, name string) (body.Body, error) {
	for _, b := range f.bodies {
		if strings.EqualFold(b.Name, name) {
			return b, nil
		}
	}
	return body.Body{}, body.ErrNotFound
}

func (f fakeWorlds) List(_ context.Context, filter body.ListFilter) ([]body.Body, error) {
	var out []body.Body
	for _, b := range f.bodies {
		if b.System == filter.System {
			out = append(out, b)
		}
	}
	return out, nil
}

func TestWorld(t *testing.T) {
	t.Parallel()

	serve := func(path string) *httptest.ResponseRecorder {
		mux := http.NewServeMux()
		mux.Handle("GET /worlds/{name}", handlers.World(logs.NewProduction(io.Discard, slog.LevelError), fakeWorlds{bodies: []body.Body{epsEri, epsEriB}}))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))
		return rec
	}

	t.Run("Planet", func(t *testing.T) {
		t.Parallel()

		got := serve("/worlds/eps%20Eri%20b")
		assert.Equal(t, http.StatusOK, got.Code)
		assert.Contains(t, got.Body.String(), "McDonald Observatory")
		assert.Contains(t, got.Body.String(), "Its star")
		assert.Contains(t, got.Body.String(), "Ran")
	})

	t.Run("Star", func(t *testing.T) {
		t.Parallel()

		got := serve("/worlds/eps%20Eri")
		assert.Equal(t, http.StatusOK, got.Code)
		html := got.Body.String()
		assert.Contains(t, html, "Epsilon Eridani")
		assert.Contains(t, html, "Orange dwarf")
		assert.Contains(t, html, "November")
		assert.Contains(t, html, "/worlds/eps%20Eri%20b")
	})

	t.Run("Unknown", func(t *testing.T) {
		t.Parallel()

		got := serve("/worlds/Vulcan")
		assert.Equal(t, http.StatusNotFound, got.Code)
		assert.Contains(t, got.Body.String(), "No world called Vulcan")
	})
}
