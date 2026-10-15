package body

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func TestSummarise(t *testing.T) {
	t.Parallel()

	got := Summarise([]Body{
		{Kind: KindPlanet, System: SolarSystem},
		{Kind: KindExoplanet, System: "51 Peg", Discovery: Discovery{Year: ptr(1995), By: "OHP", Method: "Radial Velocity"}, Sky: Sky{DistanceParsecs: ptr(15.5)}},
		{Kind: KindExoplanet, System: "Proxima Cen", Discovery: Discovery{Year: ptr(2016), By: "ESO", Method: "Radial Velocity"}, Sky: Sky{DistanceParsecs: ptr(1.3)}},
	})

	assert.Equal(t, Summary{
		Bodies: 3, Exoplanets: 2, Systems: 3, Facilities: 2, Methods: 1,
		FirstYear: 1995, LatestYear: 2016, NearestParsecs: 1.3,
	}, got)
}

func TestTop(t *testing.T) {
	t.Parallel()

	got := Top(map[string]int{"Transit": 4, "Imaging": 1, "Radial Velocity": 4}, 2)
	assert.Equal(t, []Tally{{"Radial Velocity", 4}, {"Transit", 4}}, got)
}

func TestSeedParams(t *testing.T) {
	t.Parallel()

	assert.ErrorIs(t, CreateParams{Name: "x", Kind: KindMoon, System: SolarSystem}.Validate(), ErrInvalid)
	assert.ErrorIs(t, CreateParams{Key: "k", Name: "x", Kind: "comet", System: SolarSystem}.Validate(), ErrInvalid)
	assert.NoError(t, CreateParams{Key: "k", Name: "x", Kind: KindMoon, System: SolarSystem}.Validate())
}
