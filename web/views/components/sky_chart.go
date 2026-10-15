package components

import (
	"fmt"
	"strconv"

	"workshop/internal/domain/sighting"
)

// star is one sighting placed on the chart: longitude across, latitude up.
type star struct {
	X, Y   float64
	Group  string
	Label  string
	Anchor string
}

var (
	gridMeridians = []float64{30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330}
	gridParallels = []float64{30, 60, 120, 150}
	groupOrder    = []string{"mammal", "bird", "reptile", "other"}
	groupNames    = map[string]string{
		"mammal":  "Whales, dolphins and seals",
		"bird":    "Penguins and puffins",
		"reptile": "Sea turtles",
		"other":   "Sharks",
	}
)

// stars plots every sighting with coordinates, skipping those GBIF left at 0, 0.
func stars(sightings []sighting.Sighting) []star {
	out := make([]star, 0, len(sightings))
	for _, s := range sightings {
		c := s.Location.Coordinates
		if c.Latitude == 0 && c.Longitude == 0 {
			continue
		}
		out = append(out, star{
			X:      c.Longitude + 180,
			Y:      90 - c.Latitude,
			Group:  group(s.Species.Class),
			Label:  fmt.Sprintf("%s · %s · %s", displayName(s.Species), s.HappenedAt.Format("2 Jan 2006"), place(s.Location)),
			Anchor: "#" + rowID(s),
		})
	}
	return out
}

// group sorts a taxonomic class into one of the chart's four colours.
func group(class string) string {
	switch class {
	case "Mammalia":
		return "mammal"
	case "Aves":
		return "bird"
	case "Reptilia":
		return "reptile"
	default:
		return "other"
	}
}

// legend counts the stars in each colour, in a fixed order.
func legend(points []star) []sighting.Tally {
	counts := map[string]int{}
	for _, p := range points {
		counts[p.Group]++
	}
	out := make([]sighting.Tally, 0, len(groupOrder))
	for _, g := range groupOrder {
		if counts[g] > 0 {
			out = append(out, sighting.Tally{Name: g, Count: counts[g]})
		}
	}
	return out
}

// skyLabel describes the chart for screen readers.
func skyLabel(points []star) string {
	return fmt.Sprintf("A chart of %d sightings, each placed by where on Earth it was seen. The catalogue below lists them all.", len(points))
}

// coord formats a chart position for an SVG attribute.
func coord(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// rowID is the catalogue row a star links to.
func rowID(s sighting.Sighting) string {
	return "obs-" + s.ID.String()
}
