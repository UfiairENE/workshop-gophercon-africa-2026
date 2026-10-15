package components

import (
	"slices"

	"workshop/internal/domain/sighting"
)

// MapGroup is one colour on the map, with its legend entry.
type MapGroup struct {
	Key, Label, Icon, Colour, Blurb string
	Count                           int
}

// Href is the group's own page.
func (g MapGroup) Href() string { return "/groups/" + g.Key }

// mapGroups are in legend order.
var mapGroups = []MapGroup{
	{Key: "ocean", Blurb: "Whales, dolphins, seals, sharks and sea turtles: animals that spend most or all of their lives in the sea.", Label: "Ocean", Icon: "🐋", Colour: "#0b2f6b"},
	{Key: "cats", Blurb: "Hunters built for ambush, with forward-facing eyes for judging distance and, in most species, claws that retract to stay sharp.", Label: "Big cats", Icon: "🐆", Colour: "#e97200"},
	{Key: "hunters", Blurb: "Wild dogs, hyenas and bears. Three families of meat-eaters, though some bears now live mostly on plants.", Label: "Dogs, hyenas & bears", Icon: "🐺", Colour: "#892554"},
	{Key: "grazers", Blurb: "Hoofed animals, elephants and hippos. Plant-eaters that shape whole landscapes by what they eat and where they walk.", Label: "Grazers & giants", Icon: "🦓", Colour: "#e5a823"},
	{Key: "birds", Blurb: "From penguins that fly underwater to eagles and owls. Every bird alive today descends from dinosaurs.", Label: "Birds", Icon: "🦅", Colour: "#329a94"},
	{Key: "reptiles", Blurb: "Crocodiles, alligators and lizards. They can't make their own body heat, so they move between sun and shade to stay at the right temperature.", Label: "Reptiles", Icon: "🐊", Colour: "#4d7c0f"},
	{Key: "others", Blurb: "Primates, marsupials, rodents, sloths and the egg-laying platypus: the mammals that don't fit anywhere else.", Label: "Primates & more", Icon: "🐒", Colour: "#7c5c4a"},
}

var oceanFamilies = []string{"Phocidae", "Odobenidae", "Mustelidae", "Cheloniidae", "Dermochelyidae"}

// GroupOf puts a sighting's animal into one of the map's colour groups.
func GroupOf(s sighting.Sighting) string {
	sp := s.Species
	switch {
	case sp.Order == "Cetacea", sp.Class == "Elasmobranchii", slices.Contains(oceanFamilies, sp.Family):
		return "ocean"
	case sp.Family == "Felidae":
		return "cats"
	case sp.Family == "Canidae", sp.Family == "Hyaenidae", sp.Family == "Ursidae":
		return "hunters"
	case sp.Order == "Artiodactyla", sp.Order == "Perissodactyla", sp.Order == "Proboscidea":
		return "grazers"
	case sp.Class == "Aves":
		return "birds"
	case sp.Class == "Crocodylia", sp.Class == "Squamata", sp.Class == "Testudines", sp.Class == "Reptilia":
		return "reptiles"
	}
	return "others"
}

func groupColour(key string) string { return groupByKey(key).Colour }

// MapGroups is the legend for these sightings: each group that appears, with its count.
func MapGroups(sightings []sighting.Sighting) []MapGroup {
	counts := map[string]int{}
	for _, s := range sightings {
		counts[GroupOf(s)]++
	}
	var out []MapGroup
	for _, g := range mapGroups {
		if counts[g.Key] > 0 {
			g.Count = counts[g.Key]
			out = append(out, g)
		}
	}
	return out
}
