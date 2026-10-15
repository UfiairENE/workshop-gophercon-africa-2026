package components

import (
	"fmt"

	"workshop/internal/domain/sighting"
)

// catalogueNo numbers a row like an entry in a printed catalogue.
func catalogueNo(i int) string {
	return fmt.Sprintf("%03d", i+1)
}

// firstPhoto returns the first still image, if there is one.
func firstPhoto(media []sighting.Media) (sighting.Media, bool) {
	for _, m := range media {
		if m.Type == "StillImage" && m.URL != "" {
			return m, true
		}
	}
	return sighting.Media{}, false
}
