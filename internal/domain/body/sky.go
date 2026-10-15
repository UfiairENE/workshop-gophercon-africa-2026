package body

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"
)

// constellation is one of the 88 IAU constellations: the name it goes
// by, and the genitive used in a star's name ("51 Pegasi").
type constellation struct {
	Name, Genitive string
}

var constellations = map[string]constellation{
	"And": {"Andromeda", "Andromedae"}, "Ant": {"Antlia", "Antliae"}, "Aps": {"Apus", "Apodis"},
	"Aqr": {"Aquarius", "Aquarii"}, "Aql": {"Aquila", "Aquilae"}, "Ara": {"Ara", "Arae"},
	"Ari": {"Aries", "Arietis"}, "Aur": {"Auriga", "Aurigae"}, "Boo": {"Boötes", "Boötis"},
	"Cae": {"Caelum", "Caeli"}, "Cam": {"Camelopardalis", "Camelopardalis"}, "Cnc": {"Cancer", "Cancri"},
	"CVn": {"Canes Venatici", "Canum Venaticorum"}, "CMa": {"Canis Major", "Canis Majoris"},
	"CMi": {"Canis Minor", "Canis Minoris"}, "Cap": {"Capricornus", "Capricorni"},
	"Car": {"Carina", "Carinae"}, "Cas": {"Cassiopeia", "Cassiopeiae"}, "Cen": {"Centaurus", "Centauri"},
	"Cep": {"Cepheus", "Cephei"}, "Cet": {"Cetus", "Ceti"}, "Cha": {"Chamaeleon", "Chamaeleontis"},
	"Cir": {"Circinus", "Circini"}, "Col": {"Columba", "Columbae"}, "Com": {"Coma Berenices", "Comae Berenices"},
	"CrA": {"Corona Australis", "Coronae Australis"}, "CrB": {"Corona Borealis", "Coronae Borealis"},
	"Crv": {"Corvus", "Corvi"}, "Crt": {"Crater", "Crateris"}, "Cru": {"Crux", "Crucis"},
	"Cyg": {"Cygnus", "Cygni"}, "Del": {"Delphinus", "Delphini"}, "Dor": {"Dorado", "Doradus"},
	"Dra": {"Draco", "Draconis"}, "Equ": {"Equuleus", "Equulei"}, "Eri": {"Eridanus", "Eridani"},
	"For": {"Fornax", "Fornacis"}, "Gem": {"Gemini", "Geminorum"}, "Gru": {"Grus", "Gruis"},
	"Her": {"Hercules", "Herculis"}, "Hor": {"Horologium", "Horologii"}, "Hya": {"Hydra", "Hydrae"},
	"Hyi": {"Hydrus", "Hydri"}, "Ind": {"Indus", "Indi"}, "Lac": {"Lacerta", "Lacertae"},
	"Leo": {"Leo", "Leonis"}, "LMi": {"Leo Minor", "Leonis Minoris"}, "Lep": {"Lepus", "Leporis"},
	"Lib": {"Libra", "Librae"}, "Lup": {"Lupus", "Lupi"}, "Lyn": {"Lynx", "Lyncis"},
	"Lyr": {"Lyra", "Lyrae"}, "Men": {"Mensa", "Mensae"}, "Mic": {"Microscopium", "Microscopii"},
	"Mon": {"Monoceros", "Monocerotis"}, "Mus": {"Musca", "Muscae"}, "Nor": {"Norma", "Normae"},
	"Oct": {"Octans", "Octantis"}, "Oph": {"Ophiuchus", "Ophiuchi"}, "Ori": {"Orion", "Orionis"},
	"Pav": {"Pavo", "Pavonis"}, "Peg": {"Pegasus", "Pegasi"}, "Per": {"Perseus", "Persei"},
	"Phe": {"Phoenix", "Phoenicis"}, "Pic": {"Pictor", "Pictoris"}, "Psc": {"Pisces", "Piscium"},
	"PsA": {"Piscis Austrinus", "Piscis Austrini"}, "Pup": {"Puppis", "Puppis"}, "Pyx": {"Pyxis", "Pyxidis"},
	"Ret": {"Reticulum", "Reticuli"}, "Sge": {"Sagitta", "Sagittae"}, "Sgr": {"Sagittarius", "Sagittarii"},
	"Sco": {"Scorpius", "Scorpii"}, "Scl": {"Sculptor", "Sculptoris"}, "Sct": {"Scutum", "Scuti"},
	"Ser": {"Serpens", "Serpentis"}, "Sex": {"Sextans", "Sextantis"}, "Tau": {"Taurus", "Tauri"},
	"Tel": {"Telescopium", "Telescopii"}, "Tri": {"Triangulum", "Trianguli"},
	"TrA": {"Triangulum Australe", "Trianguli Australis"}, "Tuc": {"Tucana", "Tucanae"},
	"UMa": {"Ursa Major", "Ursae Majoris"}, "UMi": {"Ursa Minor", "Ursae Minoris"}, "Vel": {"Vela", "Velorum"},
	"Vir": {"Virgo", "Virginis"}, "Vol": {"Volans", "Volantis"}, "Vul": {"Vulpecula", "Vulpeculae"},
}

// greek spells out the letter abbreviations the archive uses in star names.
var greek = map[string]string{
	"alf": "Alpha", "bet": "Beta", "gam": "Gamma", "del": "Delta", "eps": "Epsilon", "zet": "Zeta",
	"eta": "Eta", "tet": "Theta", "the": "Theta", "iot": "Iota", "kap": "Kappa", "lam": "Lambda",
	"mu": "Mu", "nu": "Nu", "ksi": "Xi", "xi": "Xi", "omi": "Omicron", "pi": "Pi", "rho": "Rho",
	"sig": "Sigma", "tau": "Tau", "ups": "Upsilon", "phi": "Phi", "chi": "Chi", "psi": "Psi", "ome": "Omega",
}

// properNames are the names the IAU has given some planet-hosting stars.
var properNames = map[string]string{
	"51 Peg": "Helvetios", "eps Eri": "Ran", "ups And": "Titawin", "47 UMa": "Chalawan",
	"55 Cnc": "Copernicus", "gam Cep": "Errai", "bet Gem": "Pollux", "eps Tau": "Ain",
	"iot Dra": "Edasich", "14 And": "Veritate", "mu Ara": "Cervantes", "alf Tau": "Aldebaran",
	"gam1 Leo": "Algieba", "alf Ari": "Hamal", "bet UMi": "Kochab", "bet And": "Mirach",
	"omi UMa": "Muscida", "del Vir": "Minelauva", "bet Cnc": "Tarf", "mu Leo": "Rasalas",
	"gam Lib": "Zubenelhakrabi", "Sun": "Sol",
}

// bayerNames are the letter-and-constellation names of bright stars the
// archive files under a catalogue number instead.
var bayerNames = map[string]string{
	"HD 62509": "bet Gem", "HD 6860": "bet And", "HD 112300": "del Vir", "HD 20794": "82 Eri",
}

// Names is how a star is known: the name the catalogue files it under,
// that name spelled out, a proper name if it has one, and the
// constellation it lies in, when its name says.
type Names struct {
	Catalogue     string
	Full          string
	Proper        string
	Constellation string
}

// StarNames works out a star's names from the one the archive uses, so
// "eps Eri" is Epsilon Eridani, called Ran, in Eridanus.
func StarNames(name string) Names {
	if bayer, ok := bayerNames[name]; ok {
		n := StarNames(bayer)
		n.Catalogue = name
		return n
	}
	n := Names{Catalogue: name, Full: name, Proper: properNames[name]}

	tokens := strings.Fields(name)
	if len(tokens) < 2 {
		return n
	}
	c, ok := constellations[tokens[1]]
	if !ok {
		return n
	}

	first := tokens[0]
	letters := strings.TrimRightFunc(first, unicode.IsDigit)
	if g, ok := greek[letters]; ok {
		first = g + first[len(letters):]
	}
	n.Full = strings.Join(append([]string{first, c.Genitive}, tokens[2:]...), " ")
	n.Constellation = c.Name
	return n
}

// StarType is what sort of star something is, and what that means it does.
type StarType struct {
	Class       string
	Description string
}

// TypeOf describes a star from its spectral type, falling back on its
// temperature when the type isn't recorded.
func TypeOf(name string, s Star) StarType {
	spec := strings.TrimSpace(s.SpectralType)
	switch {
	case strings.HasPrefix(name, "PSR"):
		return StarType{"Pulsar", "The collapsed core of a star that exploded as a supernova. It spins many times a second, sweeping beams of radio waves across space like a lighthouse."}
	case strings.HasPrefix(spec, "D") || strings.Contains(name, "WD "):
		return StarType{"White dwarf", "The hot, Earth-sized core left behind when a Sun-like star died. It no longer makes energy of its own and is slowly cooling for good."}
	}

	letter := spectralLetter(spec, s.TempK)
	colour, main := colours[letter], mainSequence[letter]
	if letter == "" {
		return StarType{"Star", "Too little is known about this star to say what kind it is."}
	}

	switch {
	case strings.Contains(spec, "IV"):
		return StarType{"Subgiant", fmt.Sprintf("A %s star running low on hydrogen in its core. It has begun to swell and brighten on its way to becoming a giant.", colour)}
	case strings.Contains(spec, "III") || strings.Contains(spec, "II"):
		return StarType{"Giant", fmt.Sprintf("A %s giant. It has used up the hydrogen in its core and swollen to many times its old size, now burning hydrogen in a shell or fusing helium.", colour)}
	}
	return main
}

var colours = map[string]string{
	"O": "blue", "B": "blue-white", "A": "white", "F": "yellow-white", "G": "yellow", "K": "orange", "M": "red", "L": "dim red",
}

var mainSequence = map[string]StarType{
	"O": {"Blue giant", "A massive, blazing-hot star burning through its hydrogen fast. It will live only a few million years and may end as a supernova."},
	"B": {"Blue-white star", "A massive, hot star, many times brighter than the Sun. It burns its hydrogen quickly and lives a few hundred million years."},
	"A": {"White star", "A hot white star, bigger and brighter than the Sun, fusing hydrogen into helium for around a billion years."},
	"F": {"Yellow-white star", "A little hotter and bigger than the Sun, steadily fusing hydrogen into helium in its core."},
	"G": {"Yellow dwarf", "A star much like our Sun, steadily fusing hydrogen into helium for about ten billion years."},
	"K": {"Orange dwarf", "Cooler and smaller than the Sun, it burns slowly and steadily for tens of billions of years: a calm home for planets."},
	"M": {"Red dwarf", "The most common kind of star. Small, cool and dim, it sips its hydrogen and could shine for trillions of years, though young ones often flare."},
	"L": {"Brown dwarf", "Too small to fuse hydrogen for long, it glows faintly from the heat it was born with and slowly cools."},
}

// spectralLetter is the star's class letter, from its type or else its temperature.
func spectralLetter(spec string, tempK *float64) string {
	if spec != "" {
		if l := strings.ToUpper(spec[:1]); strings.Contains("OBAFGKML", l) {
			return l
		}
	}
	if tempK == nil {
		return ""
	}
	switch t := *tempK; {
	case t >= 30000:
		return "O"
	case t >= 10000:
		return "B"
	case t >= 7500:
		return "A"
	case t >= 6000:
		return "F"
	case t >= 5200:
		return "G"
	case t >= 3700:
		return "K"
	case t >= 2400:
		return "M"
	default:
		return "L"
	}
}

// NakedEyeLimit is the faintest magnitude an eye can see under dark skies.
const NakedEyeLimit = 6.5

// Brightness says in words what it takes to see something of magnitude m.
func Brightness(m float64) string {
	switch {
	case m < -20:
		return "The Sun. Never look at it directly."
	case m < 1.5:
		return "One of the brightest stars in the sky."
	case m < 4:
		return "Easy to see with the naked eye, even from a city."
	case m < NakedEyeLimit:
		return "Visible to the naked eye under dark, clear skies."
	case m < 10:
		return "Too faint for the eye, but binoculars or a small telescope will show it."
	case m < 15:
		return "Needs a backyard telescope."
	default:
		return "Only large telescopes can see it."
	}
}

// Luminosity is the star's total output as a multiple of the Sun's.
func (s Star) Luminosity() (float64, bool) {
	if s.LuminosityLog == nil {
		return 0, false
	}
	return math.Pow(10, *s.LuminosityLog), true
}

// BestMonth is when a star at right ascension ra (degrees) is highest
// at midnight: the month the Sun sits on the opposite side of the sky.
func BestMonth(ra float64) time.Month {
	// The Sun is at right ascension 0 on the March equinox, day 79, and
	// moves about 0.9856 degrees a day.
	day := 79 + math.Mod(ra+180, 360)/0.9856
	return time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, int(day)%365).Month()
}

// Visibility says where on Earth a star at declination dec (degrees) rises.
func Visibility(dec float64) string {
	switch {
	case dec >= 60:
		return fmt.Sprintf("A far-northern star: it never rises south of %d°S, and from much of Europe and North America it never sets.", int(90-dec))
	case dec <= -60:
		return fmt.Sprintf("A far-southern star: it never rises north of %d°N, and from southern Africa, Australia and South America it circles the pole all night.", int(90+dec))
	case dec >= 20:
		return fmt.Sprintf("Best seen from the northern hemisphere. It rises everywhere north of %d°S.", int(90-dec))
	case dec <= -20:
		return fmt.Sprintf("Best seen from the southern hemisphere. It rises everywhere south of %d°N.", int(90+dec))
	default:
		return "Close to the celestial equator, so it can be seen from almost everywhere on Earth."
	}
}

// Starlight is how much light a planet gets compared with Earth.
func (b Body) Starlight() (float64, bool) {
	a := b.Orbit.SemiMajorAxisAU
	if a == nil || *a <= 0 || b.Kind == KindStar || b.Kind == KindMoon {
		return 0, false
	}
	if b.InSolarSystem() {
		return 1 / (*a * *a), true
	}
	r, t := b.Star.RadiusSun, b.Star.TempK
	if r == nil || t == nil {
		return 0, false
	}
	return (*r * *r) * math.Pow(*t/5772, 4) / (*a * *a), true
}

// Temperate reports whether a planet gets about the starlight a world
// needs to keep liquid water on its surface: a conservative habitable zone.
func Temperate(starlight float64) bool {
	return starlight >= 0.36 && starlight <= 1.1
}

// Gravity is the pull at a body's surface as a multiple of Earth's.
func (b Body) Gravity() (float64, bool) {
	m, r := b.Physical.MassEarth, b.Physical.RadiusEarth
	if m == nil || r == nil || *r <= 0 || b.Kind == KindStar {
		return 0, false
	}
	return *m / (*r * *r), true
}

// PlanetClass names a planet by its size.
func PlanetClass(radiusEarth float64) string {
	switch {
	case radiusEarth < 1.25:
		return "Rocky, Earth-sized world"
	case radiusEarth < 2:
		return "Super-Earth"
	case radiusEarth < 4:
		return "Mini-Neptune"
	case radiusEarth < 10:
		return "Ice giant"
	default:
		return "Gas giant"
	}
}

// Methods explains each way the archive says planets are found.
var Methods = map[string]string{
	"Transit":                       "The planet passed in front of its star, dimming it slightly on a regular beat.",
	"Radial Velocity":               "The planet's pull made its star wobble towards and away from us, shifting the colour of its light.",
	"Microlensing":                  "The star's gravity bent and magnified light from a star behind it, and the planet added a tell-tale blip.",
	"Imaging":                       "Photographed directly, with the star's glare blocked out.",
	"Transit Timing Variations":     "Another planet's pull made this one's transits arrive early or late.",
	"Eclipse Timing Variations":     "It tugged on a pair of eclipsing stars, shifting when their eclipses happen.",
	"Orbital Brightness Modulation": "The system's light rose and fell as the planet's lit side turned towards us.",
	"Pulsar Timing":                 "It disturbed the clockwork pulses of a pulsar.",
	"Pulsation Timing Variations":   "It shifted the timing of its star's natural pulsations.",
	"Astrometry":                    "Its star was seen to shift back and forth across the sky.",
	"Disk Kinematics":               "It left a wake in the disc of gas and dust around a young star.",
}
