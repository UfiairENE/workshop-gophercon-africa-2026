package components

import (
	"cmp"
	"slices"

	"workshop/internal/domain/sighting"
)

// Recording is one real call in web/assets/audio/species, with its credit.
type Recording struct {
	File, Animal, Title, Author, Licence, URL string
}

// recordings are keyed by the canonical name of the species they were
// recorded from, so an animal only plays its own call. The giraffes share one.
var recordings = map[string]Recording{
	"Acinonyx jubatus":           {File: "acinonyx-jubatus", Animal: "cheetah", Title: "Acoustic-Structure-and-Contextual-Use-of-Calls-by-Captive-Male-and-Female-Cheetahs-(Acinonyx-pone.0158546.s001.oga", Author: "Smirnova D, Volodin I, Demina T, Volodina E", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File:Acoustic-Structure-and-Contextual-Use-of-Calls-by-Captive-Male-and-Female-Cheetahs-(Acinonyx-pone.0158546.s001.oga"},
	"Aepyceros melampus":         {File: "aepyceros-melampus", Animal: "impala", Title: "Impala (Aepyceros melampus) (W1CDR0000514 BD37).ogg", Author: "Unknown", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AImpala_%28Aepyceros_melampus%29_%28W1CDR0000514_BD37%29.ogg"},
	"Alligator mississippiensis": {File: "alligator-mississippiensis", Animal: "American alligator", Title: "Alligatorbellow1.ogg", Author: "Unknown", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3AAlligatorbellow1.ogg"},
	"Aptenodytes patagonicus":    {File: "aptenodytes-patagonicus", Animal: "King Penguin", Title: "King Penguin Rookery Audio.oga", Author: "Hullwarren", Licence: "CC BY-SA 3.0", URL: "https://commons.wikimedia.org/wiki/File:King_Penguin_Rookery_Audio.oga"},
	"Ara macao":                  {File: "ara-macao", Animal: "Scarlet Macaw", Title: "Scarlet macaw 01.wav", Author: "Ganesh Mohan T", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AScarlet_macaw_01.wav"},
	"Balaenoptera acutorostrata": {File: "balaenoptera-acutorostrata", Animal: "minke whale", Title: "Atlantic minke whale call.ogg", Author: "NOAA", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3AAtlantic_minke_whale_call.ogg"},
	"Balaenoptera musculus":      {File: "balaenoptera-musculus", Animal: "blue whale", Title: "Blue Whale NE Pacific.ogg", Author: "Unknown", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3ABlue_Whale_NE_Pacific.ogg"},
	"Balaenoptera physalus":      {File: "balaenoptera-physalus", Animal: "fin whale", Title: "FinWhaleAtlantic-10x.ogg", Author: "Unknown", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3AFinWhaleAtlantic-10x.ogg"},
	"Bison bison":                {File: "bison-bison", Animal: "American bison", Title: "Bison.ogg", Author: "National Park Service", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3ABison.ogg"},
	"Bubo scandiacus":            {File: "bubo-scandiacus", Animal: "Snowy Owl", Title: "Snowy owl scream (loud).wav", Author: "HEtonis2", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3ASnowy_owl_scream_%28loud%29.wav"},
	"Crocuta crocuta":            {File: "crocuta-crocuta", Animal: "spotted hyaena", Title: "Giggling call of a spotted hyena (Crocuta crocuta) - 1472-6785-10-9-S1.oga", Author: "Mathevon N, Koralek A, Weldele M, Glickman S, Theunissen F", Licence: "CC BY 2.0", URL: "https://commons.wikimedia.org/wiki/File%3AGiggling_call_of_a_spotted_hyena_%28Crocuta_crocuta%29_-_1472-6785-10-9-S1.oga"},
	"Delphinapterus leucas":      {File: "delphinapterus-leucas", Animal: "beluga whale", Title: "Beluga vocalizations.ogg", Author: "Recorded by Fisheries Research Agency, Japan, with contract support from NOAA to Hubbs-SeaWorld Research Institute", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3ABeluga_vocalizations.ogg"},
	"Delphinus delphis":          {File: "delphinus-delphis", Animal: "short-beaked common dolphin", Title: "Dede-whistles-NOAA-PAGroup-01-short-beaked-common-dolphin-clip.mp3", Author: "NOAA Fisheries Passive Acoustics Group", Licence: "Public domain", URL: "https://www.fisheries.noaa.gov/national/science-data/sounds-ocean-mammals"},
	"Dermochelys coriacea":       {File: "dermochelys-coriacea", Animal: "Leatherback", Title: "Dermochelys coriacea 001.ogg", Author: "Philippe Kurlapski", Licence: "CC BY 2.5", URL: "https://commons.wikimedia.org/wiki/File%3ADermochelys_coriacea_001.ogg"},
	"Elephas maximus":            {File: "elephas-maximus", Animal: "Asian elephant", Title: "Elephant voice - trumpeting.ogg", Author: "தகவலுழவன்", Licence: "CC0", URL: "https://commons.wikimedia.org/wiki/File:Elephant_voice_-_trumpeting.ogg"},
	"Equus grevyi":               {File: "equus-grevyi", Animal: "Grevy's zebra", Title: "Grévys zebra (Sound Effects).ogg", Author: "DiegoC472", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AGr%C3%A9vys_zebra_%28Sound_Effects%29.ogg"},
	"Fratercula arctica":         {File: "fratercula-arctica", Animal: "Atlantic Puffin", Title: "Atlantic Puffin (Fratercula arctica) (W1CDR0001416 BD3).ogg", Author: "Unknown", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AAtlantic_Puffin_%28Fratercula_arctica%29_%28W1CDR0001416_BD3%29.ogg"},
	"Giraffa camelopardalis":     {File: "giraffa-camelopardalis", Animal: "giraffe", Title: "Giraffe Hum.oga", Author: "Anton Baotic, Florian Sicks and Angela S. Stoeger", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AGiraffe_Hum.oga"},
	"Grampus griseus":            {File: "grampus-griseus", Animal: "Risso's dolphin", Title: "Grgr-clicks-NOAA-PAGroup-08-rissos-dolphin-clip.mp3", Author: "NOAA Fisheries Passive Acoustics Group", Licence: "Public domain", URL: "https://www.fisheries.noaa.gov/national/science-data/sounds-ocean-mammals"},
	"Haliaeetus leucocephalus":   {File: "haliaeetus-leucocephalus", Animal: "Bald Eagle", Title: "Bald Eagle Yellowstone National Park.ogg", Author: "National Park Service", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3ABald_Eagle_Yellowstone_National_Park.ogg"},
	"Lemur catta":                {File: "lemur-catta", Animal: "Ring-tailed Lemur", Title: "Lemur catta--cackle4.ogg", Author: "Chris Mercer", Licence: "CC BY-SA 3.0", URL: "https://commons.wikimedia.org/wiki/File%3ALemur_catta--cackle4.ogg"},
	"Lycaon pictus":              {File: "lycaon-pictus", Animal: "African wild dog", Title: "Lycaon pictus hoo-sound.wav", Author: "Elsa Sandoval, Stichting Painted Dog Conservation", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3ALycaon_pictus_hoo-sound.wav"},
	"Megaptera novaeangliae":     {File: "megaptera-novaeangliae", Animal: "humpback whale", Title: "Humpbackwhale2.ogg", Author: "Spyrogumas", Licence: "CC0", URL: "https://commons.wikimedia.org/wiki/File%3AHumpbackwhale2.ogg"},
	"Orcinus orca":               {File: "orcinus-orca", Animal: "orca", Title: "Killer whale simple.ogg", Author: "NOAA", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File:Killer_whale_simple.ogg"},
	"Pan troglodytes":            {File: "pan-troglodytes", Animal: "chimpanzee", Title: "Pant-hoot call made by a male chimpanzee.ogg", Author: "Pawel Fedurek et al.", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File%3APant-hoot_call_made_by_a_male_chimpanzee.ogg"},
	"Panthera leo":               {File: "panthera-leo", Animal: "lion", Title: "Lion raring-sound1TamilNadu178.ogg", Author: "த*உழவன்", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File:Lion_raring-sound1TamilNadu178.ogg"},
	"Panthera onca":              {File: "panthera-onca", Animal: "jaguar", Title: "Jaguar saw.flac", Author: "About Zoos", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AJaguar_saw.flac"},
	"Panthera tigris":            {File: "panthera-tigris", Animal: "tiger", Title: "439280 schots angry-tiger.wav", Author: "schots", Licence: "CC0", URL: "https://commons.wikimedia.org/wiki/File%3A439280_schots_angry-tiger.wav"},
	"Phacochoerus africanus":     {File: "phacochoerus-africanus", Animal: "common warthog", Title: "20200322 warthog screeching.ogg", Author: "AGOoiu", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3A20200322_warthog_screeching.ogg"},
	"Phascolarctos cinereus":     {File: "phascolarctos-cinereus", Animal: "koala", Title: "Sonus naturalis - Soundscape XC446868.mp3", Author: "James Ray", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3ASonus_naturalis_-_Soundscape_XC446868.mp3"},
	"Phoca vitulina":             {File: "phoca-vitulina", Animal: "harbor seal", Title: "Phvi-roar-NOAA-PAGroup-01-harbor-seal-clip.wav", Author: "NOAA Fisheries Passive Acoustics Group", Licence: "Public domain", URL: "https://www.fisheries.noaa.gov/national/science-data/sounds-ocean-mammals"},
	"Phoenicopterus roseus":      {File: "phoenicopterus-roseus", Animal: "Greater Flamingo", Title: "Phoenicopterus roseus - Greater Flamingo XC432217.mp3", Author: "Joost van Bruggen", Licence: "CC BY-SA 4.0", URL: "https://commons.wikimedia.org/wiki/File%3APhoenicopterus_roseus_-_Greater_Flamingo_XC432217.mp3"},
	"Physeter macrocephalus":     {File: "physeter-macrocephalus", Animal: "sperm whale", Title: "Sperm Whale Normal Clicks.ogg", Author: "National Oceanic and Atmospheric Administration", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File%3ASperm_Whale_Normal_Clicks.ogg"},
	"Pongo pygmaeus":             {File: "pongo-pygmaeus", Animal: "Bornean orangutan", Title: "Speech-Like-Rhythm-in-a-Voiced-and-Voiceless-Orangutan-Call-pone.0116136.s001.oga", Author: "Lameira A, Hardus M, Bartlett A, Shumaker R, Wich S, Menken S", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File%3ASpeech-Like-Rhythm-in-a-Voiced-and-Voiceless-Orangutan-Call-pone.0116136.s001.oga"},
	"Spheniscus demersus":        {File: "spheniscus-demersus", Animal: "African Penguin", Title: "Vocal-individuality-cues-in-the-African-penguin-(Spheniscus-demersus)-a-source-filter-theory-srep17255-s2.oga", Author: "Favaro L, Gamba M, Alfieri C, Pessani D, McElligott A", Licence: "CC BY 4.0", URL: "https://commons.wikimedia.org/wiki/File%3AVocal-individuality-cues-in-the-African-penguin-%28Spheniscus-demersus%29-a-source-filter-theory-srep17255-s2.oga"},
	"Tursiops truncatus":         {File: "tursiops-truncatus", Animal: "common bottlenose dolphin", Title: "Tutr-multisound-NOAA-PAGroup-03-bottlenose-dolphin-clip.mp3", Author: "NOAA Fisheries Passive Acoustics Group", Licence: "Public domain", URL: "https://www.fisheries.noaa.gov/national/science-data/sounds-ocean-mammals"},
	"Ursus arctos":               {File: "ursus-arctos", Animal: "brown bear", Title: "Yellowstone sound library - Grizzly Bear vocalizations - 001.mp3", Author: "NPS & MSU Acoustic Atlas / Jennifer Jerrett", Licence: "Public domain", URL: "https://commons.wikimedia.org/wiki/File:Yellowstone_sound_library_-_Grizzly_Bear_vocalizations_-_001.mp3"},
}

func init() {
	g := recordings["Giraffa camelopardalis"]
	recordings["Giraffa giraffa"], recordings["Giraffa tippelskirchi"] = g, g
}

// SoundFor is the recording of this sighting's own species, or empty.
func SoundFor(s sighting.Sighting) string { return recordings[s.Species.CanonicalName].File }

// Recordings lists every recording once, by animal, for the credits page.
func Recordings() []Recording {
	seen := map[string]bool{}
	var out []Recording
	for _, r := range recordings {
		if !seen[r.File] {
			seen[r.File] = true
			out = append(out, r)
		}
	}
	slices.SortFunc(out, func(a, b Recording) int { return cmp.Compare(a.Animal, b.Animal) })
	return out
}
