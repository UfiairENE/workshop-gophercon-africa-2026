package components

import (
	"cmp"
	"fmt"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"workshop/internal/domain/body"
)

// Fact is one labelled figure on a world's page, with an optional note
// saying what it means.
type Fact struct {
	Label, Value, Note string
}

// WorldURL is the page for a body.
func WorldURL(name string) string {
	return "/worlds/" + url.PathEscape(name)
}

// LookOf is how a body is painted in 3D: Solar System bodies by hand,
// stars by their heat and exoplanets by how hot they run.
func LookOf(b body.Body) Portrait {
	if l, ok := solarLooks[b.Name]; ok && b.InSolarSystem() {
		return Portrait{Body: b, Colour: l.colour, Accent: l.accent, Bands: l.bands, Ring: l.ring}
	}
	if b.Kind == body.KindStar {
		colour, accent := starColour(b.Star.TempK)
		return Portrait{Body: b, Colour: colour, Accent: accent, Bands: 0.15}
	}
	if b.Kind == body.KindMoon {
		return Portrait{Body: b, Colour: "#d4cfcb", Accent: "#b5b0b0", Bands: 0.08}
	}
	colour, accent := heat(b.Physical.EquilibriumTempK)
	return Portrait{Body: b, Colour: colour, Accent: accent, Bands: bandsFor(b.Physical.RadiusEarth)}
}

// starColour tints a star by its surface temperature, in the site's pastels.
func starColour(k *float64) (colour, accent string) {
	switch t := deref(k); {
	case t == 0:
		return "#ffe3c7", "#ffcf9e"
	case t >= 7500:
		return "#dcefff", "#bfe0f7"
	case t >= 6000:
		return "#fff3e2", "#ffe3c7"
	case t >= 5200:
		return "#ffdcb5", "#ffcf9e"
	case t >= 3700:
		return "#ffc28f", "#f4b49c"
	default:
		return "#f6a99a", "#e3a2ef"
	}
}

// DisplayName is what a reader calls a body: a star's proper name if it
// has one, else its name spelled out.
func DisplayName(b body.Body) string {
	if b.Kind != body.KindStar || b.InSolarSystem() {
		return b.Name
	}
	n := body.StarNames(b.Name)
	return cmp.Or(n.Proper, n.Full)
}

// PlanetFacts are a planet's or moon's figures, measured against Earth.
func PlanetFacts(b body.Body) []Fact {
	var facts []Fact
	if r := b.Physical.RadiusEarth; r != nil && b.Kind == body.KindExoplanet {
		facts = append(facts, Fact{"What it is", body.PlanetClass(*r), ""})
	}
	if r := b.Physical.RadiusEarth; r != nil {
		facts = append(facts, Fact{"Size", times(r) + " Earth", fmt.Sprintf("%s km across", thousands(*r*12742))})
	}
	if m := b.Physical.MassEarth; m != nil {
		facts = append(facts, Fact{"Mass", times(m) + " Earth", ""})
	}
	if g, ok := b.Gravity(); ok {
		facts = append(facts, Fact{"Gravity", times(&g) + " Earth", fmt.Sprintf("A 70 kg person would feel like %s kg", strconv.FormatFloat(math.Round(70*g), 'f', 0, 64))})
	}
	if p := b.Orbit.PeriodDays; p != nil {
		label := "A year lasts"
		if b.Kind == body.KindMoon {
			label = "One orbit takes"
		}
		facts = append(facts, Fact{label, yearLength(p), ""})
	}
	if a := b.Orbit.SemiMajorAxisAU; a != nil {
		facts = append(facts, Fact{"Distance from its star", strconv.FormatFloat(*a, 'g', 3, 64) + " AU", "1 AU is Earth's distance from the Sun, 150 million km"})
	}
	if t := b.Physical.EquilibriumTempK; t != nil {
		facts = append(facts, Fact{"Temperature", fmt.Sprintf("%.0f K", *t), fmt.Sprintf("About %.0f °C, before any atmosphere warms it", *t-273.15)})
	}
	if s, ok := b.Starlight(); ok {
		note := "Too much for liquid water, by this measure"
		if body.Temperate(s) {
			note = "In the habitable zone: about the starlight a world needs for liquid water"
		} else if s < 0.36 {
			note = "Too little for liquid water on the surface"
		}
		facts = append(facts, Fact{"Starlight", times(&s) + " Earth's", note})
	}
	return facts
}

// StarFacts are a star's figures, measured against the Sun.
func StarFacts(b body.Body) []Fact {
	var facts []Fact
	if t := b.Star.TempK; t != nil {
		facts = append(facts, Fact{"Surface", fmt.Sprintf("%s K", thousands(*t)), "The Sun's is 5,772 K"})
	}
	if l, ok := b.Star.Luminosity(); ok {
		facts = append(facts, Fact{"Shines", strconv.FormatFloat(l, 'g', 3, 64) + "× the Sun", "Its total output of light and heat"})
	}
	if r := b.Star.RadiusSun; r != nil {
		facts = append(facts, Fact{"Size", times(r) + " the Sun", ""})
	}
	if m := b.Star.MassSun; m != nil {
		facts = append(facts, Fact{"Mass", times(m) + " the Sun", ""})
	}
	if a := b.Star.AgeGyr; a != nil && *a > 0 {
		facts = append(facts, Fact{"Age", strconv.FormatFloat(*a, 'g', 3, 64) + " billion years", "The Sun is 4.6 billion"})
	}
	if n := b.StarsInSystem; n != nil && *n > 1 {
		facts = append(facts, Fact{"Stars in the system", strconv.Itoa(*n), "It has companion stars"})
	}
	return facts
}

// SkyFacts say how bright a star looks from Earth, and when and where to look.
func SkyFacts(b body.Body) []Fact {
	var facts []Fact
	if m := b.Sky.ApparentMagnitude; m != nil {
		facts = append(facts, Fact{"How bright it looks", fmt.Sprintf("Magnitude %.1f", *m), body.Brightness(*m) + " Lower is brighter; the faintest the eye can see is about 6.5."})
	}
	if b.InSolarSystem() {
		return facts
	}
	if ra := b.Sky.RightAscension; ra != nil {
		facts = append(facts, Fact{"Best seen", body.BestMonth(*ra).String(), "When it's highest in the sky around midnight"})
	}
	if dec := b.Sky.Declination; dec != nil {
		facts = append(facts, Fact{"Where on Earth", hemisphere(*dec), body.Visibility(*dec)})
	}
	if c := body.StarNames(b.Name).Constellation; c != "" {
		facts = append(facts, Fact{"Constellation", c, ""})
	}
	if ly, ok := b.Sky.LightYears(); ok {
		facts = append(facts, Fact{"Distance", LightYears(*b.Sky.DistanceParsecs) + " light years", fmt.Sprintf("The light you see left it %s years ago", thousands(math.Round(ly)))})
	}
	return facts
}

// hemisphere is a one-word answer to where a star is best seen from.
func hemisphere(dec float64) string {
	switch {
	case dec >= 20:
		return "Northern skies"
	case dec <= -20:
		return "Southern skies"
	default:
		return "Everywhere"
	}
}

// Moons are the bodies in family that orbit b.
func Moons(b body.Body, family []body.Body) []body.Body {
	return slices.DeleteFunc(slices.Clone(family), func(o body.Body) bool { return o.Parent != b.Name })
}

// Neighbours are the bodies in family that share b's parent.
func Neighbours(b body.Body, family []body.Body) []body.Body {
	return slices.DeleteFunc(slices.Clone(family), func(o body.Body) bool { return o.Parent != b.Parent })
}

// thousands writes a whole number with commas.
func thousands(v float64) string {
	s := strconv.FormatFloat(math.Round(v), 'f', 0, 64)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// SeeableStars lists the planet-hosting stars bright enough for the naked eye, brightest first.
func SeeableStars(bodies []body.Body, n int) []body.Body {
	var out []body.Body
	for _, b := range bodies {
		if b.Kind == body.KindStar && !b.InSolarSystem() && b.Sky.ApparentMagnitude != nil && *b.Sky.ApparentMagnitude < body.NakedEyeLimit {
			out = append(out, b)
		}
	}
	slices.SortStableFunc(out, func(a, b body.Body) int {
		return cmp.Compare(*a.Sky.ApparentMagnitude, *b.Sky.ApparentMagnitude)
	})
	return out[:min(n, len(out))]
}

// CountSeeable is how many planet-hosting stars the naked eye can see.
func CountSeeable(bodies []body.Body) int {
	return len(SeeableStars(bodies, len(bodies)))
}

// magnitude shows how bright a star looks, or a dash.
func magnitude(b body.Body) string {
	if m := b.Sky.ApparentMagnitude; m != nil {
		return strconv.FormatFloat(*m, 'f', 1, 64)
	}
	return "—"
}

// bestMonth is when a star is highest at midnight, or a dash.
func bestMonth(b body.Body) string {
	if ra := b.Sky.RightAscension; ra != nil {
		return body.BestMonth(*ra).String()
	}
	return "—"
}

// planetCount is how many planets a star is known to have.
func planetCount(b body.Body) string {
	if n := b.PlanetsInSystem; n != nil {
		return strconv.Itoa(*n)
	}
	return "—"
}
