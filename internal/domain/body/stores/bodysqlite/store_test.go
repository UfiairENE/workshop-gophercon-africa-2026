package bodysqlite_test

import (
	"testing"
	"workshop/internal/domain/body"
	"workshop/internal/domain/body/stores/bodysqlite"
	"workshop/internal/infra/db/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

var (
	earth = body.CreateParams{
		Key: "sol:Earth", Name: "Earth", Kind: body.KindPlanet, System: body.SolarSystem, Parent: "Sun",
		Physical: body.Physical{RadiusEarth: ptr(1.0)},
	}
	pegasi = body.CreateParams{
		Key: "exo:51 Peg b", Name: "51 Peg b", Kind: body.KindExoplanet, System: "51 Peg", Parent: "51 Peg",
		Discovery: body.Discovery{Year: ptr(1995), By: "Haute-Provence Observatory", Method: "Radial Velocity"},
		Sky:       body.Sky{RightAscension: ptr(344.37), Declination: ptr(20.77), DistanceParsecs: ptr(15.46)},
	}
)

func TestStore(t *testing.T) {
	t.Parallel()

	ctx, conn, cleanup := dbtest.Setup(t)
	t.Cleanup(cleanup)
	store := bodysqlite.New(conn)

	added, err := store.CreateMany(ctx, []body.CreateParams{pegasi, earth})
	require.NoError(t, err)
	assert.Equal(t, 2, added)

	t.Run("Seeding again adds nothing", func(t *testing.T) {
		added, err := store.CreateMany(ctx, []body.CreateParams{earth})
		require.NoError(t, err)
		assert.Equal(t, 0, added)
	})

	t.Run("List puts the Solar System first", func(t *testing.T) {
		got, err := store.List(ctx, body.ListFilter{})
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "Earth", got[0].Name)
		assert.Equal(t, "51 Peg b", got[1].Name)
		assert.Equal(t, 1995, *got[1].Discovery.Year)
		assert.InDelta(t, 15.46, *got[1].Sky.DistanceParsecs, 1e-9)
	})

	t.Run("List by kind", func(t *testing.T) {
		got, err := store.List(ctx, body.ListFilter{Kinds: []body.Kind{body.KindExoplanet}})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, body.KindExoplanet, got[0].Kind)
	})

	t.Run("Find by name ignores case", func(t *testing.T) {
		got, err := store.FindByName(ctx, "51 PEG B")
		require.NoError(t, err)
		assert.Equal(t, "Haute-Provence Observatory", got.Discovery.By)

		_, err = store.Find(ctx, got.ID)
		require.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		_, err := store.FindByName(ctx, "Vulcan")
		assert.ErrorIs(t, err, body.ErrNotFound)
	})
}
