# Design: Earthbound Observatory

Based on the BlueYard Capital system from Refero
(https://styles.refero.design/style/ea11696a-17c0-41f8-9b08-c820851e0ea9): a printed editorial
object, a white page with one warm peach wash and big tightly set headlines, with celestial
artwork hanging in the white space.

## Pages

- `/`, **Worlds**: planets and stars from the `body` domain (NASA Exoplanet Archive plus the
  bundled Solar System). Orrery banner, Portraits, Stars you can see, Star map, Almanac,
  Nearest worlds. Portraits, star cards, table rows, orrery planets and star-map stars all link
  to a world's page.
- `/worlds/{name}`, **one world**: a hero with its name and a large 3D body, then fact grids.
  A planet or moon gets At a glance, How we know, Its star, and its moons or neighbours. A star
  gets its names (proper, full and catalogue), what kind of star it is and what that means,
  In our sky (brightness, best month, where on Earth, constellation, distance), The star, and
  its planets. The plain-English rules live in `internal/domain/body/sky.go`.
- Anything unmatched gets the error page with a 404 ("Nothing orbits here.").

## Colours (`:root` in `assets/css/app.css`)

| Token | Hex | Use |
|---|---|---|
| `--canvas` | `#ffffff` | Page background, cards |
| `--ink` | `#3a3a3e` | Nearly all text, hairlines, orbit lines |
| `--carbon` | `#090b11` | Headlines, the star map panel |
| `--ash` | `#b5b0b0` | Soft borders only, never text |
| `--muted` | `#66666b` | Secondary text (5.7:1 on white) |
| `--apricot` | `#ffcf9e` | Banner wash, warm plate, Sun, Sun-like stars |
| `--iris` | `#babfff` | Cool plate, cold worlds |
| `--fuchsia` | `#e3a2ef` | Bloom plate, scorching worlds, cool red dwarfs |
| `--polar` | `#bfe0f7` | Temperate worlds, hot blue-white stars |

The tints are fills and artwork only: never button backgrounds or body text.

## Lettering

Instrument Sans (Google Fonts), weights 400 and 500 only. Headlines are 400 with tight negative
tracking (-0.03em to -0.035em) and a line height of 1. Weight 500 is for nav, labels and the
wordmark. Never use positive letter-spacing.

## Shape and spacing

- Sharp corners everywhere (radius 0). No drop shadows.
- Sections are 90px apart, each opened by a 1px `--ink` rule above its title.
- Max width 1200px. Side gutter 16px on phones, 32px from 720px.
- Cards are hairline-bordered and transparent, so the portraits canvas shows through.

## 3D and motion (three.js r186 from jsDelivr, `assets/js/scenes/`)

- `orrery.js`: the Solar System in the banner. Scrolling swings the view from overhead to
  edge-on and draws the camera in. One Earth year takes 14 seconds.
- `portraits.js`: one canvas behind the Portraits grid, drawing each planet into its card with
  scissor viewports (one WebGL context for all of them). Hovering a card spins its planet up.
- `starmap.js`: every planet-hosting star at its real direction, with distance log-scaled.
  Scrolling flies the camera in. Drag to turn on desktop. On touch screens dragging is off, so
  the page still scrolls.
- `world.js`: the large body at the top of a world's page. Planets turn and lean towards the
  pointer; stars glow, tinted by their temperature.
- `planet-material.js`: the shared banded pastel planet shader.
- Hover: cards lift 4–6px, table rows tint `--polar`, and links in lists slide 6px right.
- Fact grids: hairline `--ash` cells, a 12px label, a large figure, and a muted note saying
  what the number means.
- Reduced motion: scenes draw one still frame, and lifts are off.

## Avoid

Rounded corners, shadows or glows on the page itself, coloured buttons, and more than two or three
WebGL canvases on a page.
