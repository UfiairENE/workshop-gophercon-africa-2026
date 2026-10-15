package body

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"uuid"
	"workshop/internal/clients/exoplanets"
	"workshop/internal/clients/solarsystem"
)

// Store defines the operations used to persist and query bodies.
type Store interface {
	Find(ctx context.Context, id uuid.UUID) (Body, error)
	FindByName(ctx context.Context, name string) (Body, error)
	List(ctx context.Context, filter ListFilter) ([]Body, error)
	// SaveMany saves bodies in one go, refreshing any whose Key is
	// already stored, and reports how many it added.
	SaveMany(ctx context.Context, in []CreateParams) (int, error)
}

// ExoplanetSource fetches the live archive, for seeding with --live.
type ExoplanetSource interface {
	Fetch(ctx context.Context) ([]exoplanets.Planet, error)
}

// ListFilter narrows a listing. The zero value matches everything.
type ListFilter struct {
	Kinds  []Kind
	System string
	Limit  int
}

// Service is the entry point for working with bodies. Handlers and
// commands use it rather than the store.
type Service struct {
	logger *slog.Logger
	store  Store
	live   ExoplanetSource
}

// NewService returns a Service backed by store, fetching live planets
// from live when asked to.
func NewService(logger *slog.Logger, store Store, live ExoplanetSource) *Service {
	return &Service{logger: logger, store: store, live: live}
}

var (
	// ErrNotFound is returned when no body matches.
	ErrNotFound = errors.New("body not found")

	// ErrInvalid is returned for a filter or body that doesn't make sense.
	ErrInvalid = errors.New("invalid body")
)

// Find returns the body with id, or ErrNotFound.
func (s *Service) Find(ctx context.Context, id uuid.UUID) (Body, error) {
	return s.store.Find(ctx, id)
}

// FindByName returns the body called name, ignoring case, or ErrNotFound.
func (s *Service) FindByName(ctx context.Context, name string) (Body, error) {
	return s.store.FindByName(ctx, name)
}

// List returns the bodies matching filter, the Solar System first and
// then exoplanets in the order they were found.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]Body, error) {
	for _, k := range filter.Kinds {
		if !k.Valid() {
			return nil, fmt.Errorf("%w: unknown kind %q", ErrInvalid, k)
		}
	}
	return s.store.List(ctx, filter)
}

// SeedRequest says where to seed exoplanets from.
type SeedRequest struct {
	// Live fetches the NASA Exoplanet Archive rather than reading the
	// snapshot bundled with the app.
	Live bool
}

// SeedResult counts what a seed found and what it added.
type SeedResult struct {
	Found int
	Added int
}

// Seed saves the Solar System, every confirmed exoplanet and the stars
// they orbit. It's safe to run again: bodies already stored are refreshed
// rather than duplicated.
func (s *Service) Seed(ctx context.Context, req SeedRequest) (SeedResult, error) {
	sol, err := solarsystem.Bodies()
	if err != nil {
		return SeedResult{}, err
	}

	var planets []exoplanets.Planet
	if req.Live {
		planets, err = s.live.Fetch(ctx)
	} else {
		planets, err = exoplanets.Snapshot()
	}
	if err != nil {
		return SeedResult{}, err
	}

	params := make([]CreateParams, 0, len(sol)+len(planets)*2)
	for _, b := range sol {
		params = append(params, fromSolarSystem(b))
	}
	hosts := map[string]bool{}
	for _, p := range planets {
		// Each system's star is saved once, from its first planet's row.
		if !hosts[p.Host] {
			hosts[p.Host] = true
			params = append(params, starOf(p))
		}
		params = append(params, fromExoplanet(p))
	}
	for _, p := range params {
		if err := p.Validate(); err != nil {
			return SeedResult{}, err
		}
	}

	added, err := s.store.SaveMany(ctx, params)
	if err != nil {
		return SeedResult{}, err
	}

	s.logger.InfoContext(ctx, "body: seeded", slog.Int("found", len(params)), slog.Int("added", added))
	return SeedResult{Found: len(params), Added: added}, nil
}

// Validate checks a body has what every listing relies on.
func (p CreateParams) Validate() error {
	switch {
	case p.Key == "":
		return fmt.Errorf("%w: %q has no key", ErrInvalid, p.Name)
	case p.Name == "":
		return fmt.Errorf("%w: %q has no name", ErrInvalid, p.Key)
	case !p.Kind.Valid():
		return fmt.Errorf("%w: %q has unknown kind %q", ErrInvalid, p.Name, p.Kind)
	case p.System == "":
		return fmt.Errorf("%w: %q has no system", ErrInvalid, p.Name)
	}
	return nil
}

func fromSolarSystem(b solarsystem.Body) CreateParams {
	return CreateParams{
		Key:    "sol:" + b.Name,
		Name:   b.Name,
		Kind:   Kind(b.Kind),
		System: SolarSystem,
		Parent: b.Parent,
		Discovery: Discovery{
			Year: b.DiscoveredYear,
			By:   b.DiscoveredBy,
		},
		Physical: Physical{
			RadiusEarth:      b.RadiusEarth,
			MassEarth:        b.MassEarth,
			EquilibriumTempK: b.TempK,
		},
		Orbit: Orbit{
			PeriodDays:      b.OrbitalPeriodDays,
			SemiMajorAxisAU: b.SemiMajorAxisAU,
		},
		Star: Star{
			SpectralType:  b.SpectralType,
			TempK:         b.StarTempK,
			RadiusSun:     b.StarRadiusSun,
			MassSun:       b.StarMassSun,
			AgeGyr:        b.StarAgeGyr,
			LuminosityLog: sunLuminosity(b),
		},
		Sky: Sky{ApparentMagnitude: b.ApparentMagnitude},
	}
}

// sunLuminosity is 0, the Sun's own, for the Sun and nothing for the rest.
func sunLuminosity(b solarsystem.Body) *float64 {
	if b.Kind != string(KindStar) {
		return nil
	}
	zero := 0.0
	return &zero
}

func fromExoplanet(p exoplanets.Planet) CreateParams {
	return CreateParams{
		Key:             "exo:" + p.Name,
		Name:            p.Name,
		Kind:            KindExoplanet,
		System:          p.Host,
		Parent:          p.Host,
		PlanetsInSystem: p.PlanetsInSystem,
		StarsInSystem:   p.StarsInSystem,
		Discovery: Discovery{
			Year:   p.DiscoveryYear,
			By:     p.DiscoveryFacility,
			Method: p.DiscoveryMethod,
		},
		Physical: Physical{
			RadiusEarth:      p.RadiusEarth,
			MassEarth:        p.MassEarth,
			EquilibriumTempK: p.EquilibriumTempK,
		},
		Orbit: Orbit{
			PeriodDays:      p.OrbitalPeriodDays,
			SemiMajorAxisAU: p.SemiMajorAxisAU,
		},
		Star: starFields(p),
		Sky:  skyFields(p),
	}
}

// starOf is the star an exoplanet orbits, as a body of its own.
func starOf(p exoplanets.Planet) CreateParams {
	var names []string
	for _, n := range []string{p.HDName, p.HIPName} {
		if n != "" && n != p.Host {
			names = append(names, n)
		}
	}
	return CreateParams{
		Key:             "star:" + p.Host,
		Name:            p.Host,
		Kind:            KindStar,
		System:          p.Host,
		PlanetsInSystem: p.PlanetsInSystem,
		StarsInSystem:   p.StarsInSystem,
		CatalogueNames:  names,
		Physical: Physical{
			RadiusEarth: sunRadiiToEarth(p.StarRadiusSun),
		},
		Star: starFields(p),
		Sky:  skyFields(p),
	}
}

// sunRadiiToEarth converts a star's radius to Earth radii, the unit every
// body's size is kept in.
func sunRadiiToEarth(r *float64) *float64 {
	if r == nil {
		return nil
	}
	v := *r * 109.1
	return &v
}

func starFields(p exoplanets.Planet) Star {
	return Star{
		SpectralType:  p.StarSpectralType,
		TempK:         p.StarTempK,
		RadiusSun:     p.StarRadiusSun,
		MassSun:       p.StarMassSun,
		AgeGyr:        p.StarAgeGyr,
		LuminosityLog: p.StarLuminosityLog,
	}
}

func skyFields(p exoplanets.Planet) Sky {
	return Sky{
		RightAscension:    p.RightAscension,
		Declination:       p.Declination,
		DistanceParsecs:   p.DistanceParsecs,
		ApparentMagnitude: p.ApparentMagnitude,
	}
}
