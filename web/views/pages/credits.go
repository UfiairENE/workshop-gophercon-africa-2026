package pages

// Credit is one openly licensed file used on the site.
type Credit struct {
	Used, Title, Author, Licence, URL string
}

// credits lists the footage the site shows; recordings come from components.Recordings.
var credits = []Credit{
	{Used: "Hero video", Title: "Lions (Panthera leo) on the move ....webm", Author: "Bernard DUPONT", Licence: "CC BY-SA 2.0", URL: "https://commons.wikimedia.org/wiki/File:Lions_(Panthera_leo)_on_the_move_....webm"},
	{Used: "Hero video", Title: "Savanna Elephant (Loxodonta africana) bull ....webm", Author: "Bernard DUPONT", Licence: "CC BY-SA 2.0", URL: "https://commons.wikimedia.org/wiki/File:Savanna_Elephant_(Loxodonta_africana)_bull_....webm"},
	{Used: "Hero video", Title: "Giraffe (Giraffa camelopardalis) male ....webm", Author: "Bernard DUPONT", Licence: "CC BY-SA 2.0", URL: "https://commons.wikimedia.org/wiki/File:Giraffe_(Giraffa_camelopardalis)_male_....webm"},
}
