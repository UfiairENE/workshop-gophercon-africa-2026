package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/web/views/pages"
)

type SightingsLister interface {
	List(ctx context.Context, filter sighting.ListFilter) ([]sighting.Sighting, error)
}

// seaLife is every animal the home page charts: the ones that live in or by the sea.
var seaLife = []species.Animal{
	species.AnimalBlueWhale,
	species.AnimalHumpbackWhale,
	species.AnimalSpermWhale,
	species.AnimalFinWhale,
	species.AnimalCommonMinkeWhale,
	species.AnimalGreyWhale,
	species.AnimalNorthAtlanticRightWhale,
	species.AnimalBowheadWhale,
	species.AnimalBeluga,
	species.AnimalNarwhal,
	species.AnimalOrca,
	species.AnimalCommonBottlenoseDolphin,
	species.AnimalCommonDolphin,
	species.AnimalSpinnerDolphin,
	species.AnimalDuskyDolphin,
	species.AnimalRissosDolphin,
	species.AnimalAmazonRiverDolphin,
	species.AnimalHarbourSeal,
	species.AnimalLeopardSeal,
	species.AnimalWalrus,
	species.AnimalSeaOtter,
	species.AnimalGreatWhiteShark,
	species.AnimalWhaleShark,
	species.AnimalGreenTurtle,
	species.AnimalLeatherbackTurtle,
	species.AnimalEmperorPenguin,
	species.AnimalKingPenguin,
	species.AnimalGentooPenguin,
	species.AnimalAfricanPenguin,
	species.AnimalAtlanticPuffin,
}

func Sea(logger *slog.Logger, sightings SightingsLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		items, err := sightings.List(ctx, sighting.ListFilter{
			Animals: seaLife,
		})
		if err != nil {
			logger.ErrorContext(ctx, "handlers: listing sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to load sightings.")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err = pages.Sea(items).Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering sightings: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
