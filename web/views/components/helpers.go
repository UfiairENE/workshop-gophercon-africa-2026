package components

import (
	"fmt"

	"workshop/internal/domain/sighting"
)

// Href is the page for one sighting.
func Href(s sighting.Sighting) string { return "/sightings/" + s.ID.String() }

// Name is the sighting's animal, as a reader would call it.
func Name(s sighting.Sighting) string { return displayName(s.Species) }

// Place is the sighting's location in words.
func Place(s sighting.Sighting) string { return place(s.Location) }

// Count is how many animals were seen, or empty when unknown.
func Count(s sighting.Sighting) string {
	if s.IndividualCount == nil {
		return ""
	}
	return individuals(*s.IndividualCount)
}

// Cover is the first still image of a sighting, if it has one.
func Cover(s sighting.Sighting) (string, bool) {
	for _, m := range s.Media {
		if m.Type == "StillImage" && m.URL != "" {
			return m.URL, true
		}
	}
	return "", false
}

// Temp renders a temperature in whole degrees.
func Temp(c float64) string { return fmt.Sprintf("%.0f°C", c) }

// WeatherIcon picks an emoji for a WMO weather code.
func WeatherIcon(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code <= 2:
		return "🌤️"
	case code == 3:
		return "☁️"
	case code <= 48:
		return "🌫️"
	case code <= 57:
		return "🌦️"
	case code <= 67, code >= 80 && code <= 82:
		return "🌧️"
	case code <= 77, code == 85, code == 86:
		return "❄️"
	case code >= 95:
		return "⛈️"
	}
	return "🌡️"
}

// Compass turns a wind bearing into a compass point.
func Compass(deg int) string {
	points := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	return points[((deg%360+360)%360+22)/45%8]
}

// MapPoint is one marker, handed to map.js as JSON.
type MapPoint struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Title  string  `json:"title"`
	Date   string  `json:"date"`
	Place  string  `json:"place"`
	Count  string  `json:"count,omitzero"`
	Href   string  `json:"href"`
	Group  string  `json:"group"`
	Colour string  `json:"colour"`
	Photo  string  `json:"photo,omitzero"`
}

// MapPoints turns sightings into markers, leaving out any with no position.
func MapPoints(sightings []sighting.Sighting) []MapPoint {
	points := make([]MapPoint, 0, len(sightings))
	for _, s := range sightings {
		if s.Location.Latitude == 0 && s.Location.Longitude == 0 {
			continue
		}
		photo, _ := Cover(s)
		group := GroupOf(s)
		points = append(points, MapPoint{
			Lat: s.Location.Latitude, Lng: s.Location.Longitude,
			Title: Name(s), Date: s.HappenedAt.Format("2 Jan 2006"),
			Place: Place(s), Count: Count(s), Href: Href(s),
			Group: group, Colour: groupColour(group), Photo: photo,
		})
	}
	return points
}
