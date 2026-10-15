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
	// CreateMany saves bodies in one go, skipping any whose Key is
	// already stored, and reports how many it added.
	CreateMany(ctx context.Context, in []CreateParams) (int, error)
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

// Seed saves the Solar System and every confirmed exoplanet. It's safe
// to run again: bodies already stored are skipped.
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

	params := make([]CreateParams, 0, len(sol)+len(planets))
	for _, b := range sol {
		params = append(params, fromSolarSystem(b))
	}
	for _, p := range planets {
		params = append(params, fromExoplanet(p))
	}
	for _, p := range params {
		if err := p.Validate(); err != nil {
			return SeedResult{}, err
		}
	}

	added, err := s.store.CreateMany(ctx, params)
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
			SpectralType: b.SpectralType,
			TempK:        b.StarTempK,
			RadiusSun:    b.StarRadiusSun,
		},
	}
}

func fromExoplanet(p exoplanets.Planet) CreateParams {
	return CreateParams{
		Key:             "exo:" + p.Name,
		Name:            p.Name,
		Kind:            KindExoplanet,
		System:          p.Host,
		Parent:          p.Host,
		PlanetsInSystem: p.PlanetsInSystem,
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
		Star: Star{
			SpectralType: p.StarSpectralType,
			TempK:        p.StarTempK,
			RadiusSun:    p.StarRadiusSun,
		},
		Sky: Sky{
			RightAscension:  p.RightAscension,
			Declination:     p.Declination,
			DistanceParsecs: p.DistanceParsecs,
		},
	}
}
