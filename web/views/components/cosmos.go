package components

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"

	"workshop/internal/domain/body"
)

// orbiter is one body circling the Sun in the orrery.
type orbiter struct {
	Name   string  `json:"name"`
	Kind   string  `json:"kind"`
	A      float64 `json:"a"`
	Period float64 `json:"period"`
	Radius float64 `json:"radius"`
	Colour string  `json:"colour"`
	Ring   bool    `json:"ring"`
}

// Portrait is a body drawn up close, with the look its scene uses.
type Portrait struct {
	Body   body.Body
	Colour string
	Accent string
	Bands  float64
	Ring   bool
}

// look is how a Solar System body is painted: its colour, a second
// colour for its bands, how banded it is and whether it has rings.
type look struct {
	colour, accent string
	bands          float64
	ring           bool
}

var solarLooks = map[string]look{
	"Sun":      {"#ffcf9e", "#ffb877", 0.15, false},
	"Mercury":  {"#c9c4c0", "#b5b0b0", 0.1, false},
	"Venus":    {"#ffdcb5", "#ffcf9e", 0.35, false},
	"Earth":    {"#bfe0f7", "#babfff", 0.25, false},
	"Mars":     {"#f4b49c", "#e8977f", 0.15, false},
	"Jupiter":  {"#ffd9b0", "#e3a2ef", 0.8, false},
	"Saturn":   {"#f6e2b8", "#ffcf9e", 0.6, true},
	"Uranus":   {"#cfeaf7", "#bfe0f7", 0.2, true},
	"Neptune":  {"#babfff", "#9fa6ff", 0.35, false},
	"Ceres":    {"#d4cfcb", "#b5b0b0", 0.05, false},
	"Pluto":    {"#efd9c8", "#d9bfae", 0.2, false},
	"Haumea":   {"#e4e1df", "#cdc9c6", 0.05, false},
	"Makemake": {"#f0d0bb", "#e3b9a3", 0.1, false},
	"Eris":     {"#ecebea", "#d6d3d1", 0.05, false},
}

// featured are the strangers given a portrait beside our own planets.
var featured = []string{"51 Peg b", "Proxima Cen b", "TRAPPIST-1 e", "Kepler-452 b", "K2-18 b", "WASP-12 b"}

// Orrery lists the planets and dwarf planets with an orbit to draw.
func Orrery(bodies []body.Body) []orbiter {
	var out []orbiter
	for _, b := range bodies {
		if !b.InSolarSystem() || b.Parent != "Sun" || b.Orbit.SemiMajorAxisAU == nil || b.Orbit.PeriodDays == nil {
			continue
		}
		l := solarLooks[b.Name]
		out = append(out, orbiter{
			Name:   b.Name,
			Kind:   string(b.Kind),
			A:      *b.Orbit.SemiMajorAxisAU,
			Period: *b.Orbit.PeriodDays,
			Radius: deref(b.Physical.RadiusEarth),
			Colour: cmp.Or(l.colour, "#d9d4d0"),
			Ring:   l.ring,
		})
	}
	return out
}

// starMap is what the star map draws: each planet-hosting star once, as
// right ascension, declination, distance in parsecs and temperature (0 if
// unknown), with names alongside for hovering and clicking.
type starMap struct {
	Names  []string     `json:"names"`
	Points [][4]float64 `json:"points"`
}

// StarMap gathers the stars for the star map.
func StarMap(bodies []body.Body) starMap {
	var m starMap
	for _, b := range bodies {
		s := b.Sky
		if b.Kind != body.KindStar || b.InSolarSystem() || s.RightAscension == nil || s.Declination == nil || s.DistanceParsecs == nil {
			continue
		}
		m.Names = append(m.Names, b.Name)
		m.Points = append(m.Points, [4]float64{round(*s.RightAscension, 3), round(*s.Declination, 3), round(*s.DistanceParsecs, 2), round(deref(b.Star.TempK), 0)})
	}
	return m
}

// Portraits picks our planets, Pluto and a few famous exoplanets to draw up close.
func Portraits(bodies []body.Body) []Portrait {
	byName := make(map[string]body.Body, len(bodies))
	for _, b := range bodies {
		byName[b.Name] = b
	}

	var out []Portrait
	for _, name := range []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto"} {
		if b, ok := byName[name]; ok {
			l := solarLooks[name]
			out = append(out, Portrait{Body: b, Colour: l.colour, Accent: l.accent, Bands: l.bands, Ring: l.ring})
		}
	}
	for _, name := range featured {
		if b, ok := byName[name]; ok {
			colour, accent := heat(b.Physical.EquilibriumTempK)
			out = append(out, Portrait{Body: b, Colour: colour, Accent: accent, Bands: bandsFor(b.Physical.RadiusEarth)})
		}
	}
	return out
}

// heat colours an exoplanet by how hot it is: lilac when cold, sky blue
// when temperate, peach when hot and pink when scorching.
func heat(k *float64) (colour, accent string) {
	switch t := deref(k); {
	case t == 0:
		return "#d9d4d0", "#b5b0b0"
	case t < 200:
		return "#babfff", "#9fa6ff"
	case t < 330:
		return "#bfe0f7", "#babfff"
	case t < 1000:
		return "#ffcf9e", "#f4b49c"
	default:
		return "#e3a2ef", "#ffcf9e"
	}
}

// bandsFor gives gas giants stripes and leaves rocky worlds smooth.
func bandsFor(radius *float64) float64 {
	if deref(radius) > 4 {
		return 0.7
	}
	return 0.15
}

// Nearest lists the n exoplanets closest to Earth.
func Nearest(bodies []body.Body, n int) []body.Body {
	var out []body.Body
	for _, b := range bodies {
		if b.Kind == body.KindExoplanet && b.Sky.DistanceParsecs != nil {
			out = append(out, b)
		}
	}
	slices.SortStableFunc(out, func(a, b body.Body) int {
		return cmp.Or(cmp.Compare(*a.Sky.DistanceParsecs, *b.Sky.DistanceParsecs), cmp.Compare(a.Name, b.Name))
	})
	return out[:min(n, len(out))]
}

// TopBy ranks exoplanets by key, as rows for Bars.
func TopBy(bodies []body.Body, n int, key func(body.Body) string) []body.Tally {
	exo := slices.DeleteFunc(slices.Clone(bodies), func(b body.Body) bool { return b.Kind != body.KindExoplanet })
	return body.Top(body.CountBy(exo, key), n)
}

// catalogueNo numbers a row like an entry in a printed catalogue.
func catalogueNo(i int) string {
	return fmt.Sprintf("%03d", i+1)
}

// Method is how a body was found, for TopBy.
func Method(b body.Body) string { return b.Discovery.Method }

// Facility is who found a body, for TopBy.
func Facility(b body.Body) string { return b.Discovery.By }

// LightYears shows a distance in parsecs as light years.
func LightYears(pc float64) string {
	ly := pc * 3.26156
	if ly < 100 {
		return strconv.FormatFloat(ly, 'f', 1, 64)
	}
	return strconv.FormatFloat(math.Round(ly), 'f', 0, 64)
}

// bodyDistance shows how far a body's system is, or a dash.
func bodyDistance(b body.Body) string {
	if b.Sky.DistanceParsecs == nil {
		return "—"
	}
	return LightYears(*b.Sky.DistanceParsecs) + " ly"
}

// found says when a body was found.
func found(b body.Body) string {
	if b.Discovery.Year == nil {
		return "Known since antiquity"
	}
	return "Found " + strconv.Itoa(*b.Discovery.Year)
}

// times shows an optional multiple of Earth, or a dash.
func times(v *float64) string {
	if v == nil {
		return "—"
	}
	switch x := *v; {
	case x >= 100:
		return strconv.FormatFloat(x, 'f', 0, 64) + "×"
	case x >= 0.1:
		return strconv.FormatFloat(x, 'f', 2, 64) + "×"
	default:
		return strconv.FormatFloat(x, 'g', 2, 64) + "×"
	}
}

// yearLength shows an orbital period in days or years, whichever reads better.
func yearLength(v *float64) string {
	if v == nil {
		return "—"
	}
	if d := *v; d < 730 {
		return strconv.FormatFloat(d, 'f', 1, 64) + " days"
	}
	return strconv.FormatFloat(*v/365.256, 'f', 1, 64) + " years"
}

// orreryLabel describes the orrery for screen readers.
func orreryLabel(o []orbiter) string {
	return fmt.Sprintf("The Sun and %d planets and dwarf planets turning on their orbits.", len(o))
}

// starMapLabel describes the star map for screen readers.
func starMapLabel(n int) string {
	return fmt.Sprintf("A 3D map of %d star systems with known planets, each at its real direction from Earth.", n)
}

func deref(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func round(v float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(v*p) / p
}

func num(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// year shows a discovery year, or a dash.
func year(y *int) string {
	if y == nil {
		return "—"
	}
	return strconv.Itoa(*y)
}
