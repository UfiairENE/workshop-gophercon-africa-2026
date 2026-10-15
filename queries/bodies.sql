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

-- name: BodyCreate :execrows
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
	distance_parsecs
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (catalogue_key) DO NOTHING;
