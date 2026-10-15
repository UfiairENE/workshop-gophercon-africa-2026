-- +goose Up
-- +goose StatementBegin
ALTER TABLE bodies ADD COLUMN stars_in_system INTEGER;
ALTER TABLE bodies ADD COLUMN apparent_magnitude REAL;
ALTER TABLE bodies ADD COLUMN star_luminosity_log REAL;
ALTER TABLE bodies ADD COLUMN star_mass_sun REAL CHECK (star_mass_sun > 0);
ALTER TABLE bodies ADD COLUMN star_age_gyr REAL CHECK (star_age_gyr >= 0);
ALTER TABLE bodies ADD COLUMN catalogue_names TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bodies DROP COLUMN catalogue_names;
ALTER TABLE bodies DROP COLUMN star_age_gyr;
ALTER TABLE bodies DROP COLUMN star_mass_sun;
ALTER TABLE bodies DROP COLUMN star_luminosity_log;
ALTER TABLE bodies DROP COLUMN apparent_magnitude;
ALTER TABLE bodies DROP COLUMN stars_in_system;
-- +goose StatementEnd
