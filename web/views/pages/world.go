package pages

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/domain/body"
	"workshop/web/views/components"
)

// planetLede introduces a planet or moon in a sentence.
func planetLede(b body.Body, host body.Body) string {
	switch {
	case b.Kind == body.KindMoon:
		return fmt.Sprintf("A moon of %s.", b.Parent)
	case b.InSolarSystem():
		return fmt.Sprintf("A %s of our own Solar System.", strings.ToLower(b.Kind.String()))
	}

	what := "A planet"
	if r := b.Physical.RadiusEarth; r != nil {
		class := body.PlanetClass(*r)
		what = "A " + strings.ToLower(class[:1]) + class[1:]
	}
	star := b.Parent
	if host.Name != "" {
		star = components.DisplayName(host)
	}
	s := fmt.Sprintf("%s orbiting %s", what, star)
	if b.Sky.DistanceParsecs != nil {
		s += ", " + components.LightYears(*b.Sky.DistanceParsecs) + " light years from Earth"
	}
	return s + "."
}

// worldDescription is the page's summary for search engines and link previews.
func worldDescription(b body.Body, host body.Body) string {
	if b.Kind == body.KindStar {
		return components.DisplayName(b) + ": " + body.TypeOf(b.Name, b.Star).Description
	}
	return planetLede(b, host)
}

// artLabel describes the 3D art for screen readers.
func artLabel(b body.Body) string {
	return "A slowly turning 3D picture of " + components.DisplayName(b) + "."
}

// starEyebrow is the line above a star's name.
func starEyebrow(b body.Body, names body.Names) string {
	if names.Constellation != "" {
		return "Star in " + names.Constellation
	}
	if b.InSolarSystem() {
		return "Our star"
	}
	return "Star"
}

// familyTitle heads the list of a body's neighbours.
func familyTitle(b body.Body) string {
	switch {
	case b.Kind == body.KindMoon:
		return "Other moons of " + b.Parent
	case b.InSolarSystem():
		return "Its neighbours"
	default:
		return "Other planets around " + b.Parent
	}
}

// by says who found a body, if anyone's recorded.
func by(b body.Body) string {
	if b.Discovery.By == "" {
		return ""
	}
	return " by " + b.Discovery.By
}

func year(y *int) string {
	if y == nil {
		return ""
	}
	return strconv.Itoa(*y)
}

// otherNames are the catalogue names a star goes by besides the ones in its heading.
func otherNames(b body.Body, names body.Names) []string {
	var out []string
	if b.Name != names.Full && b.Name != names.Proper {
		out = append(out, b.Name)
	}
	return append(out, b.CatalogueNames...)
}

func joinNames(names []string) string {
	return strings.Join(names, " and ")
}
