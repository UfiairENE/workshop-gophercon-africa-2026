// Package exoplanets reads confirmed planets from NASA's Exoplanet
// Archive, either from the snapshot bundled beside this file or live
// from the archive's TAP service.
package exoplanets

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BaseURL is the archive's synchronous TAP endpoint.
const BaseURL = "https://exoplanetarchive.ipac.caltech.edu/TAP/sync"

// query is the ADQL behind both the snapshot and a live fetch: one row per
// planet from the archive's composite parameters table.
const query = "select pl_name,hostname,sy_pnum,discoverymethod,disc_year,disc_facility," +
	"pl_orbper,pl_orbsmax,pl_rade,pl_bmasse,pl_eqt,st_spectype,st_teff,st_rad,ra,dec,sy_dist," +
	"sy_snum,sy_vmag,st_lum,st_mass,st_age,hd_name,hip_name " +
	"from pscomppars order by disc_year,pl_name"

// snapshot is the archive's answer to query, saved so seeding needs no network.
//
//go:embed snapshot.csv
var snapshot []byte

// Planet is one confirmed exoplanet and what the archive knows of its star.
// A nil field is one the archive left blank.
type Planet struct {
	Name              string
	Host              string
	PlanetsInSystem   *int
	DiscoveryMethod   string
	DiscoveryYear     *int
	DiscoveryFacility string
	OrbitalPeriodDays *float64
	SemiMajorAxisAU   *float64
	RadiusEarth       *float64
	MassEarth         *float64
	EquilibriumTempK  *float64
	StarSpectralType  string
	StarTempK         *float64
	StarRadiusSun     *float64
	RightAscension    *float64
	Declination       *float64
	DistanceParsecs   *float64
	StarsInSystem     *int
	// ApparentMagnitude is how bright the system looks from Earth in
	// visible light: lower is brighter, and 6 is about the faintest an
	// eye can see.
	ApparentMagnitude *float64
	// StarLuminosityLog is the star's output as log10 of the Sun's.
	StarLuminosityLog *float64
	StarMassSun       *float64
	StarAgeGyr        *float64
	HDName            string
	HIPName           string
}

// Client fetches planets from the archive.
type Client struct {
	baseURL string
	http    *http.Client
}

// New returns a Client for the archive at baseURL.
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 2 * time.Minute}}
}

// Snapshot returns the planets bundled with the app.
func Snapshot() ([]Planet, error) {
	return parse(bytes.NewReader(snapshot))
}

// Fetch asks the live archive for every confirmed planet.
func (c *Client) Fetch(ctx context.Context) ([]Planet, error) {
	u := c.baseURL + "?" + url.Values{"query": {query}, "format": {"csv"}}.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exoplanets: fetching: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exoplanets: archive answered %s", res.Status)
	}
	return parse(res.Body)
}

// columns is the order query selects in, which parse checks the header against.
var columns = []string{
	"pl_name", "hostname", "sy_pnum", "discoverymethod", "disc_year", "disc_facility",
	"pl_orbper", "pl_orbsmax", "pl_rade", "pl_bmasse", "pl_eqt", "st_spectype", "st_teff",
	"st_rad", "ra", "dec", "sy_dist", "sy_snum", "sy_vmag", "st_lum", "st_mass", "st_age", "hd_name",
	"hip_name",
}

func parse(r io.Reader) ([]Planet, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = len(columns)

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("exoplanets: reading header: %w", err)
	}
	for i, name := range columns {
		if header[i] != name {
			return nil, fmt.Errorf("exoplanets: column %d is %q, want %q", i+1, header[i], name)
		}
	}

	var planets []Planet
	for line := 2; ; line++ {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return planets, nil
		} else if err != nil {
			return nil, fmt.Errorf("exoplanets: line %d: %w", line, err)
		}

		p := Planet{
			Name:              rec[0],
			Host:              rec[1],
			DiscoveryMethod:   rec[3],
			DiscoveryFacility: rec[5],
			StarSpectralType:  rec[11],
			HDName:            rec[22],
			HIPName:           rec[23],
		}
		var errs []error
		p.PlanetsInSystem = integer(rec[2], &errs)
		p.DiscoveryYear = integer(rec[4], &errs)
		p.OrbitalPeriodDays = number(rec[6], &errs)
		p.SemiMajorAxisAU = number(rec[7], &errs)
		p.RadiusEarth = number(rec[8], &errs)
		p.MassEarth = number(rec[9], &errs)
		p.EquilibriumTempK = number(rec[10], &errs)
		p.StarTempK = number(rec[12], &errs)
		p.StarRadiusSun = number(rec[13], &errs)
		p.RightAscension = number(rec[14], &errs)
		p.Declination = number(rec[15], &errs)
		p.DistanceParsecs = number(rec[16], &errs)
		p.StarsInSystem = integer(rec[17], &errs)
		p.ApparentMagnitude = number(rec[18], &errs)
		p.StarLuminosityLog = number(rec[19], &errs)
		p.StarMassSun = number(rec[20], &errs)
		p.StarAgeGyr = number(rec[21], &errs)
		if err := errors.Join(errs...); err != nil {
			return nil, fmt.Errorf("exoplanets: line %d: %w", line, err)
		}

		planets = append(planets, p)
	}
}

// number parses an optional decimal, recording a malformed one in errs.
func number(s string, errs *[]error) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	return &v
}

// integer parses an optional whole number, recording a malformed one in errs.
func integer(s string, errs *[]error) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	return &v
}
