package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/common/printer"
	"workshop/internal/common/printer/styles"
	"workshop/internal/domain/body"

	"github.com/urfave/cli/v3"
)

func bodiesCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:    "bodies",
		Aliases: []string{"b"},
		Usage:   "Seed, list and inspect stars, planets, moons and exoplanets",
		Commands: []*cli.Command{
			{
				Name:  "seed",
				Usage: "Save the Solar System and every confirmed exoplanet, from the bundled NASA snapshot or live with --live",
				Action: func(ctx context.Context, _ *cli.Command) error {
					return runBodySeed(ctx, app)
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List bodies, the Solar System first",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "kind",
						Aliases: []string{"k"},
						Usage:   "only these kinds: " + kindNames(),
					},
					&cli.StringFlag{
						Name:    "system",
						Aliases: []string{"s"},
						Usage:   "only bodies in this system, such as 'TRAPPIST-1'",
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"n"},
						Usage:   "show at most this many",
						Value:   40,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					kinds := make([]body.Kind, 0, len(cmd.StringSlice("kind")))
					for _, k := range cmd.StringSlice("kind") {
						kinds = append(kinds, body.Kind(k))
					}
					bodies, err := app.Bodies.List(ctx, body.ListFilter{
						Kinds:  kinds,
						System: cmd.String("system"),
						Limit:  cmd.Int("limit"),
					})
					if err != nil {
						return err
					}

					var hint view.Hint
					if len(bodies) > 0 {
						hint = view.Hintf("%d bodies %s workshop bodies find <name>", len(bodies), styles.IconArrow)
					}
					view.Render(app.Printer,
						view.Table[body.Body]{
							Title:   "Bodies",
							Headers: []string{"NAME", "KIND", "SYSTEM", "FOUND", "BY", "RADIUS (EARTHS)"},
							Items:   bodies,
							Empty:   "No bodies yet, run 'workshop bodies seed'",
							Row: func(b body.Body) []string {
								return []string{b.Name, b.Kind.String(), b.System, year(b.Discovery.Year), cmp.Or(b.Discovery.By, "—"), measure(b.Physical.RadiusEarth, 2)}
							},
						},
						hint,
					)
					return nil
				},
			},
			{
				Name:      "find",
				Aliases:   []string{"show"},
				Usage:     "Show everything known about one body",
				ArgsUsage: "<name>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := strings.Join(cmd.Args().Slice(), " ")
					if name == "" {
						return errors.New("give a body's name, such as 'workshop bodies find TRAPPIST-1 e'")
					}
					b, err := app.Bodies.FindByName(ctx, name)
					if errors.Is(err, body.ErrNotFound) {
						return fmt.Errorf("no body called %q, try 'workshop bodies list'", name)
					} else if err != nil {
						return err
					}

					distance := "—"
					if ly, ok := b.Sky.LightYears(); ok {
						distance = fmt.Sprintf("%s light years", strconv.FormatFloat(ly, 'f', 1, 64))
					}
					view.Render(app.Printer,
						view.Card{Title: b.Name, Rows: []printer.KV{
							{Key: "Kind", Value: b.Kind.String()},
							{Key: "System", Value: b.System},
							{Key: "Orbits", Value: cmp.Or(b.Parent, "—")},
							{Key: "Found", Value: year(b.Discovery.Year)},
							{Key: "By", Value: cmp.Or(b.Discovery.By, "—")},
							{Key: "Method", Value: cmp.Or(b.Discovery.Method, "—")},
							{Key: "Radius", Value: measure(b.Physical.RadiusEarth, 3) + " × Earth"},
							{Key: "Mass", Value: measure(b.Physical.MassEarth, 3) + " × Earth"},
							{Key: "Year length", Value: measure(b.Orbit.PeriodDays, 2) + " days"},
							{Key: "Temperature", Value: measure(b.Physical.EquilibriumTempK, 0) + " K"},
							{Key: "Star", Value: strings.TrimSpace(b.Star.SpectralType + " " + measure(b.Star.TempK, 0) + " K")},
							{Key: "Distance", Value: distance},
						}},
					)
					return nil
				},
			},
		},
	}
}

// runBodySeed seeds bodies, from the live archive under --live.
func runBodySeed(ctx context.Context, app *app) error {
	app.Printer.Muted("Seeding the Solar System and the exoplanet catalogue…")
	res, err := app.Bodies.Seed(ctx, body.SeedRequest{Live: app.live})
	if err != nil {
		return fmt.Errorf("seeding bodies: %w", err)
	}
	app.Printer.Successf("Found %d bodies, added %d new\n", res.Found, res.Added)
	return nil
}

func kindNames() string {
	names := make([]string, len(body.Kinds))
	for i, k := range body.Kinds {
		names[i] = string(k)
	}
	return strings.Join(names, ", ")
}

// year shows a discovery year, or that it's been known since antiquity.
func year(y *int) string {
	if y == nil {
		return "antiquity"
	}
	return strconv.Itoa(*y)
}

// measure shows an optional number, or a dash if nobody has measured it.
func measure(v *float64, places int) string {
	if v == nil {
		return "—"
	}
	return strconv.FormatFloat(*v, 'f', places, 64)
}
