package body

import (
	"time"
	"uuid"
)

// Kind is what sort of body something is.
type Kind string

// The kinds of body the catalogue holds.
const (
	KindStar        Kind = "star"
	KindPlanet      Kind = "planet"
	KindDwarfPlanet Kind = "dwarf_planet"
	KindMoon        Kind = "moon"
	KindExoplanet   Kind = "exoplanet"
)

// Kinds lists every kind, nearest to home first.
var Kinds = []Kind{KindStar, KindPlanet, KindDwarfPlanet, KindMoon, KindExoplanet}

// SolarSystem is the System every body orbiting the Sun belongs to.
const SolarSystem = "Solar System"

// String implements fmt.Stringer, in words a reader would use.
func (k Kind) String() string {
	switch k {
	case KindStar:
		return "Star"
	case KindPlanet:
		return "Planet"
	case KindDwarfPlanet:
		return "Dwarf planet"
	case KindMoon:
		return "Moon"
	case KindExoplanet:
		return "Exoplanet"
	default:
		return string(k)
	}
}

// Valid reports whether k is one of Kinds.
func (k Kind) Valid() bool {
	for _, known := range Kinds {
		if k == known {
			return true
		}
	}
	return false
}

type (
	// Body is a star, planet, dwarf planet or moon. A nil measurement
	// is one nobody has made yet, which for exoplanets is most of them.
	Body struct {
		ID uuid.UUID `json:"id"`
		// Key identifies the body in its source catalogue, such as
		// "exo:51 Peg b" or "sol:Earth", so seeding twice adds nothing.
		Key    string `json:"key"`
		Name   string `json:"name"`
		Kind   Kind   `json:"kind"`
		System string `json:"system"`
		// Parent is what the body orbits: a star for planets, a planet for moons.
		Parent          string    `json:"parent,omitzero"`
		PlanetsInSystem *int      `json:"planetsInSystem,omitzero"`
		Discovery       Discovery `json:"discovery"`
		Physical        Physical  `json:"physical"`
		Orbit           Orbit     `json:"orbit"`
		Star            Star      `json:"star"`
		Sky             Sky       `json:"sky"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
	}
	// Discovery is who found a body, when and how. A zero Year means it
	// has been known since antiquity.
	Discovery struct {
		Year   *int   `json:"year,omitzero"`
		By     string `json:"by,omitzero"`
		Method string `json:"method,omitzero"`
	}
	// Physical is a body's size and mass, measured against Earth's.
	Physical struct {
		RadiusEarth      *float64 `json:"radiusEarth,omitzero"`
		MassEarth        *float64 `json:"massEarth,omitzero"`
		EquilibriumTempK *float64 `json:"equilibriumTempK,omitzero"`
	}
	// Orbit is how a body goes round its parent.
	Orbit struct {
		PeriodDays      *float64 `json:"periodDays,omitzero"`
		SemiMajorAxisAU *float64 `json:"semiMajorAxisAU,omitzero"`
	}
	// Star describes the star a body belongs to, or the body itself if
	// it's a star.
	Star struct {
		SpectralType string   `json:"spectralType,omitzero"`
		TempK        *float64 `json:"tempK,omitzero"`
		RadiusSun    *float64 `json:"radiusSun,omitzero"`
	}
	// Sky is where a body's system appears from Earth: right ascension
	// and declination in degrees, and its distance in parsecs.
	Sky struct {
		RightAscension  *float64 `json:"rightAscension,omitzero"`
		Declination     *float64 `json:"declination,omitzero"`
		DistanceParsecs *float64 `json:"distanceParsecs,omitzero"`
	}
	// CreateParams is a body without its ID or timestamps.
	CreateParams struct {
		Key             string    `json:"key"`
		Name            string    `json:"name"`
		Kind            Kind      `json:"kind"`
		System          string    `json:"system"`
		Parent          string    `json:"parent,omitzero"`
		PlanetsInSystem *int      `json:"planetsInSystem,omitzero"`
		Discovery       Discovery `json:"discovery"`
		Physical        Physical  `json:"physical"`
		Orbit           Orbit     `json:"orbit"`
		Star            Star      `json:"star"`
		Sky             Sky       `json:"sky"`
	}
)

// InSolarSystem reports whether the body orbits our Sun, or is it.
func (b Body) InSolarSystem() bool {
	return b.System == SolarSystem
}

// LightYears is how far the body's system is from Earth, if known.
func (s Sky) LightYears() (float64, bool) {
	if s.DistanceParsecs == nil {
		return 0, false
	}
	return *s.DistanceParsecs * 3.26156, true
}
