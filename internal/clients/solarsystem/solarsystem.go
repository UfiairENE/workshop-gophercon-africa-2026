// Package solarsystem holds a hand-kept catalogue of our own Solar
// System: the Sun, its planets and dwarf planets, and their best-known
// moons. It's small and changes rarely, so it's bundled rather than fetched.
package solarsystem

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed bodies.json
var catalogue []byte

// Body is one entry in the catalogue. Parent is what it orbits, and a
// nil DiscoveredYear means it has been known since antiquity.
type Body struct {
	Name              string   `json:"name"`
	Kind              string   `json:"kind"`
	Parent            string   `json:"parent"`
	RadiusEarth       *float64 `json:"radiusEarth"`
	MassEarth         *float64 `json:"massEarth"`
	OrbitalPeriodDays *float64 `json:"orbitalPeriodDays"`
	SemiMajorAxisAU   *float64 `json:"semiMajorAxisAU"`
	TempK             *float64 `json:"tempK"`
	SpectralType      string   `json:"spectralType"`
	StarTempK         *float64 `json:"starTempK"`
	StarRadiusSun     *float64 `json:"starRadiusSun"`
	DiscoveredYear    *int     `json:"discoveredYear"`
	DiscoveredBy      string   `json:"discoveredBy"`
}

// Bodies returns the bundled catalogue, Sun first.
func Bodies() ([]Body, error) {
	var bodies []Body
	if err := json.Unmarshal(catalogue, &bodies); err != nil {
		return nil, fmt.Errorf("solarsystem: decoding catalogue: %w", err)
	}
	return bodies, nil
}
