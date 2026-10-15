package bodysqlite

import (
	"uuid"
	"workshop/internal/domain/body"
	"workshop/internal/infra/db/sqlc"
)

func transform(r db.Body) body.Body {
	return body.Body{
		ID:              r.ID,
		Key:             r.CatalogueKey,
		Name:            r.Name,
		Kind:            body.Kind(r.Kind),
		System:          r.System,
		Parent:          r.Parent,
		PlanetsInSystem: r.PlanetsInSystem,
		Discovery: body.Discovery{
			Year:   r.DiscoveredYear,
			By:     r.DiscoveredBy,
			Method: r.DiscoveryMethod,
		},
		Physical: body.Physical{
			RadiusEarth:      r.RadiusEarth,
			MassEarth:        r.MassEarth,
			EquilibriumTempK: r.EquilibriumTempK,
		},
		Orbit: body.Orbit{
			PeriodDays:      r.OrbitalPeriodDays,
			SemiMajorAxisAU: r.SemiMajorAxisAu,
		},
		Star: body.Star{
			SpectralType: r.StarSpectralType,
			TempK:        r.StarTempK,
			RadiusSun:    r.StarRadiusSun,
		},
		Sky: body.Sky{
			RightAscension:  r.RightAscensionDeg,
			Declination:     r.DeclinationDeg,
			DistanceParsecs: r.DistanceParsecs,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func toCreateParams(id uuid.UUID, in body.CreateParams) db.BodyCreateParams {
	return db.BodyCreateParams{
		ID:                id,
		CatalogueKey:      in.Key,
		Name:              in.Name,
		Kind:              string(in.Kind),
		System:            in.System,
		Parent:            in.Parent,
		PlanetsInSystem:   in.PlanetsInSystem,
		DiscoveredYear:    in.Discovery.Year,
		DiscoveredBy:      in.Discovery.By,
		DiscoveryMethod:   in.Discovery.Method,
		RadiusEarth:       in.Physical.RadiusEarth,
		MassEarth:         in.Physical.MassEarth,
		OrbitalPeriodDays: in.Orbit.PeriodDays,
		SemiMajorAxisAu:   in.Orbit.SemiMajorAxisAU,
		EquilibriumTempK:  in.Physical.EquilibriumTempK,
		StarSpectralType:  in.Star.SpectralType,
		StarTempK:         in.Star.TempK,
		StarRadiusSun:     in.Star.RadiusSun,
		RightAscensionDeg: in.Sky.RightAscension,
		DeclinationDeg:    in.Sky.Declination,
		DistanceParsecs:   in.Sky.DistanceParsecs,
	}
}
