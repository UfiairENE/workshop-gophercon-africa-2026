package body

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStarNames(t *testing.T) {
	t.Parallel()

	tt := map[string]Names{
		"eps Eri":     {Catalogue: "eps Eri", Full: "Epsilon Eridani", Proper: "Ran", Constellation: "Eridanus"},
		"51 Peg":      {Catalogue: "51 Peg", Full: "51 Pegasi", Proper: "Helvetios", Constellation: "Pegasus"},
		"Proxima Cen": {Catalogue: "Proxima Cen", Full: "Proxima Centauri", Constellation: "Centaurus"},
		"psi1 Dra B":  {Catalogue: "psi1 Dra B", Full: "Psi1 Draconis B", Constellation: "Draco"},
		"TRAPPIST-1":  {Catalogue: "TRAPPIST-1", Full: "TRAPPIST-1"},
		"GJ 887":      {Catalogue: "GJ 887", Full: "GJ 887"},
		"HD 62509":    {Catalogue: "HD 62509", Full: "Beta Geminorum", Proper: "Pollux", Constellation: "Gemini"},
	}
	for in, want := range tt {
		assert.Equal(t, want, StarNames(in), in)
	}
}

func TestTypeOf(t *testing.T) {
	t.Parallel()

	temp := func(k float64) *float64 { return &k }
	assert.Equal(t, "Yellow dwarf", TypeOf("51 Peg", Star{SpectralType: "G5V"}).Class)
	assert.Equal(t, "Red dwarf", TypeOf("TRAPPIST-1", Star{SpectralType: "M8.0 V"}).Class)
	assert.Equal(t, "Giant", TypeOf("bet Gem", Star{SpectralType: "K0 III"}).Class)
	assert.Equal(t, "Subgiant", TypeOf("x", Star{SpectralType: "G2 IV"}).Class)
	assert.Equal(t, "Orange dwarf", TypeOf("x", Star{TempK: temp(4500)}).Class)
	assert.Equal(t, "Pulsar", TypeOf("PSR B1257+12", Star{}).Class)
	assert.Equal(t, "Star", TypeOf("x", Star{}).Class)
}

func TestBestMonth(t *testing.T) {
	t.Parallel()

	// Orion's belt, around 84°, is a December–January sight.
	assert.Contains(t, []time.Month{time.December, time.January}, BestMonth(84))
	// 51 Peg, at 344°, is an autumn star in the north.
	assert.Contains(t, []time.Month{time.August, time.September}, BestMonth(344))
}

func TestStarlight(t *testing.T) {
	t.Parallel()

	earth := Body{Kind: KindPlanet, System: SolarSystem, Orbit: Orbit{SemiMajorAxisAU: ptr(1.0)}}
	got, ok := earth.Starlight()
	assert.True(t, ok)
	assert.InDelta(t, 1, got, 1e-9)
	assert.True(t, Temperate(got))

	hot := Body{Kind: KindExoplanet, System: "x", Orbit: Orbit{SemiMajorAxisAU: ptr(0.05)}, Star: Star{RadiusSun: ptr(1.0), TempK: ptr(5772.0)}}
	got, ok = hot.Starlight()
	assert.True(t, ok)
	assert.InDelta(t, 400, got, 1e-6)
	assert.False(t, Temperate(got))

	_, ok = Body{Kind: KindExoplanet, System: "x"}.Starlight()
	assert.False(t, ok)
}

func TestGravity(t *testing.T) {
	t.Parallel()

	g, ok := Body{Kind: KindPlanet, Physical: Physical{MassEarth: ptr(317.8), RadiusEarth: ptr(11.21)}}.Gravity()
	assert.True(t, ok)
	assert.InDelta(t, 2.53, g, 0.01)
}
