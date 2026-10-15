-- name: BodyFind :one
SELECT * FROM bodies WHERE id = ?;

-- name: BodyFindByKey :one
SELECT * FROM bodies WHERE catalogue_key = ?;

-- name: BodyFindByName :one
SELECT * FROM bodies WHERE name = ? COLLATE NOCASE LIMIT 1;

-- name: BodyList :many
SELECT * FROM bodies
WHERE (CAST(sqlc.arg('any_kind') AS BOOLEAN) OR kind IN (SELECT value FROM json_each(CAST(sqlc.arg('kinds') AS TEXT))))
	AND (CAST(sqlc.arg('system') AS TEXT) = '' OR system = CAST(sqlc.arg('system') AS TEXT) COLLATE NOCASE)
ORDER BY system = 'Solar System' DESC,
	CASE kind WHEN 'star' THEN 0 WHEN 'planet' THEN 1 WHEN 'dwarf_planet' THEN 2 WHEN 'moon' THEN 3 ELSE 4 END,
	CASE WHEN system = 'Solar System' THEN semi_major_axis_au END,
	discovered_year IS NULL DESC, discovered_year, name
LIMIT sqlc.arg('limit');

-- name: BodyCount :one
SELECT COUNT(*) FROM bodies;

-- name: BodyUpsert :exec
-- Seeding again refreshes what the catalogue says, keeping each body's ID.
INSERT INTO bodies (
	id,
	catalogue_key,
	name,
	kind,
	system,
	parent,
	planets_in_system,
	discovered_year,
	discovered_by,
	discovery_method,
	radius_earth,
	mass_earth,
	orbital_period_days,
	semi_major_axis_au,
	equilibrium_temp_k,
	star_spectral_type,
	star_temp_k,
	star_radius_sun,
	right_ascension_deg,
	declination_deg,
	distance_parsecs,
	stars_in_system,
	apparent_magnitude,
	star_luminosity_log,
	star_mass_sun,
	star_age_gyr,
	catalogue_names
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (catalogue_key) DO UPDATE SET
	name = excluded.name,
	kind = excluded.kind,
	system = excluded.system,
	parent = excluded.parent,
	planets_in_system = excluded.planets_in_system,
	discovered_year = excluded.discovered_year,
	discovered_by = excluded.discovered_by,
	discovery_method = excluded.discovery_method,
	radius_earth = excluded.radius_earth,
	mass_earth = excluded.mass_earth,
	orbital_period_days = excluded.orbital_period_days,
	semi_major_axis_au = excluded.semi_major_axis_au,
	equilibrium_temp_k = excluded.equilibrium_temp_k,
	star_spectral_type = excluded.star_spectral_type,
	star_temp_k = excluded.star_temp_k,
	star_radius_sun = excluded.star_radius_sun,
	right_ascension_deg = excluded.right_ascension_deg,
	declination_deg = excluded.declination_deg,
	distance_parsecs = excluded.distance_parsecs,
	stars_in_system = excluded.stars_in_system,
	apparent_magnitude = excluded.apparent_magnitude,
	star_luminosity_log = excluded.star_luminosity_log,
	star_mass_sun = excluded.star_mass_sun,
	star_age_gyr = excluded.star_age_gyr,
	catalogue_names = excluded.catalogue_names,
	updated_at = CURRENT_TIMESTAMP;
