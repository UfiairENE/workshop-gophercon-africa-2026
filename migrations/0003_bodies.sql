-- +goose Up
-- +goose StatementBegin
CREATE TABLE bodies (
	id                    UUID PRIMARY KEY,
	created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	catalogue_key         TEXT NOT NULL UNIQUE,
	name                  TEXT NOT NULL,
	kind                  TEXT NOT NULL CHECK (kind IN ('star', 'planet', 'dwarf_planet', 'moon', 'exoplanet')),
	system                TEXT NOT NULL,
	parent                TEXT NOT NULL DEFAULT '',
	planets_in_system     INTEGER,
	discovered_year       INTEGER,
	discovered_by         TEXT NOT NULL DEFAULT '',
	discovery_method      TEXT NOT NULL DEFAULT '',
	radius_earth          REAL CHECK (radius_earth > 0),
	mass_earth            REAL CHECK (mass_earth > 0),
	orbital_period_days   REAL CHECK (orbital_period_days > 0),
	semi_major_axis_au    REAL CHECK (semi_major_axis_au > 0),
	equilibrium_temp_k    REAL,
	star_spectral_type    TEXT NOT NULL DEFAULT '',
	star_temp_k           REAL,
	star_radius_sun       REAL,
	right_ascension_deg   REAL CHECK (right_ascension_deg BETWEEN 0 AND 360),
	declination_deg       REAL CHECK (declination_deg BETWEEN -90 AND 90),
	distance_parsecs      REAL CHECK (distance_parsecs >= 0)
);
CREATE INDEX bodies_kind ON bodies (kind);
CREATE INDEX bodies_name ON bodies (name COLLATE NOCASE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bodies;
-- +goose StatementEnd
