package components

import (
	"fmt"
	"slices"

	"github.com/a-h/templ"

	"workshop/internal/domain/sighting"
)

var cetaceans = []string{"Balaenidae", "Balaenopteridae", "Delphinidae", "Physeteridae", "Phocoenidae", "Eschrichtiidae", "Ziphiidae", "Kogiidae", "Monodontidae"}

// KindOf decides which drawing and which voice a sighting's animal gets.
func KindOf(s sighting.Sighting) string {
	sp := s.Species
	switch {
	case sp.Family == "Elephantidae":
		return "elephant"
	case sp.Genus == "Acinonyx":
		return "cheetah"
	case sp.Family == "Felidae":
		return "lion"
	case sp.Family == "Canidae":
		return "wolf"
	case sp.Family == "Ursidae":
		return "bear"
	case sp.Family == "Spheniscidae":
		return "penguin"
	case sp.Class == "Aves":
		return "bird"
	case slices.Contains(cetaceans, sp.Family) || sp.Order == "Cetacea":
		return "whale"
	case sp.Order == "Testudines" || sp.Family == "Cheloniidae" || sp.Family == "Dermochelyidae":
		return "turtle"
	}
	return "wild"
}

// SlideStyle staggers a hero slide's crossfade.
func SlideStyle(i, n int) templ.SafeCSS {
	if n == 0 {
		n = 1
	}
	return templ.SafeCSS(fmt.Sprintf("--i:%d;--n:%d;--delay:%ds;--cycle:%ds", i, n, i*6, n*6))
}
