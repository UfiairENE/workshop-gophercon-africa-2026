package exoplanets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sample = `pl_name,hostname,sy_pnum,discoverymethod,disc_year,disc_facility,pl_orbper,pl_orbsmax,pl_rade,pl_bmasse,pl_eqt,st_spectype,st_teff,st_rad,ra,dec,sy_dist,sy_snum,sy_vmag,st_lum,st_mass,st_age,hd_name,hip_name
"51 Peg b","51 Peg",1,"Radial Velocity",1995,"Haute-Provence Observatory",4.2307869,0.0527,13.7,146,,"G2 IV",5758.00,1.15,344.3665852,20.7690844,15.4614,1,5.45,0.136,1.07,10.0,"HD 217014","HIP 113357"
`

func TestSnapshot(t *testing.T) {
	t.Parallel()

	planets, err := Snapshot()
	require.NoError(t, err)
	assert.Greater(t, len(planets), 5000)
	for _, p := range planets {
		assert.NotEmpty(t, p.Name)
		assert.NotEmpty(t, p.Host)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		got, err := parse(strings.NewReader(sample))
		require.NoError(t, err)
		require.Len(t, got, 1)

		p := got[0]
		assert.Equal(t, "51 Peg b", p.Name)
		assert.Equal(t, "51 Peg", p.Host)
		assert.Equal(t, 1995, *p.DiscoveryYear)
		assert.Equal(t, "Haute-Provence Observatory", p.DiscoveryFacility)
		assert.InDelta(t, 15.4614, *p.DistanceParsecs, 1e-9)
		assert.Nil(t, p.EquilibriumTempK)
		assert.InDelta(t, 5.45, *p.ApparentMagnitude, 1e-9)
		assert.Equal(t, "HD 217014", p.HDName)
		assert.Equal(t, 1, *p.StarsInSystem)
	})

	t.Run("Wrong columns", func(t *testing.T) {
		t.Parallel()

		_, err := parse(strings.NewReader(strings.Replace(sample, "pl_name", "name", 1)))
		assert.ErrorContains(t, err, `want "pl_name"`)
	})

	t.Run("Bad number", func(t *testing.T) {
		t.Parallel()

		_, err := parse(strings.NewReader(strings.Replace(sample, "1995", "soon", 1)))
		assert.ErrorContains(t, err, "line 2")
	})
}

func TestFetch(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "csv", r.URL.Query().Get("format"))
			assert.Contains(t, r.URL.Query().Get("query"), "from pscomppars")
			_, _ = w.Write([]byte(sample))
		}))
		t.Cleanup(srv.Close)

		got, err := New(srv.URL).Fetch(context.Background())
		require.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("Server error", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		t.Cleanup(srv.Close)

		_, err := New(srv.URL).Fetch(context.Background())
		assert.ErrorContains(t, err, "503")
	})
}
