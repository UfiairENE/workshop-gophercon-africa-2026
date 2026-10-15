package components

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/domain/sighting"
)

// The sparkline's drawing box, matching the viewBox in the templ.
const sparkW, sparkH = 100.0, 24.0

// years lists every year from first to last, so a quiet year stays a dip.
func years(first, last int) []int {
	ys := make([]int, 0, max(last-first+1, 0))
	for y := first; y <= last; y++ {
		ys = append(ys, y)
	}
	return ys
}

// peak finds the busiest year and its count.
func peak(byYear map[int]int) (year, count int) {
	for y, n := range byYear {
		if n > count || (n == count && y < year) {
			year, count = y, n
		}
	}
	return year, count
}

// sparkPoints plots a count per year across the box, the busiest year at the top.
func sparkPoints(byYear map[int]int, first, last int) string {
	ys := years(first, last)
	_, high := peak(byYear)
	var b strings.Builder
	for i, y := range ys {
		x := 0.0
		if len(ys) > 1 {
			x = float64(i) / float64(len(ys)-1) * sparkW
		}
		top := sparkH
		if high > 0 {
			top = sparkH - float64(byYear[y])/float64(high)*(sparkH-1)
		}
		fmt.Fprintf(&b, "%.1f,%.1f ", x, top)
	}
	return strings.TrimSpace(b.String())
}

// sparkArea closes the line down to the baseline, for the fill under it.
func sparkArea(byYear map[int]int, first, last int) string {
	return fmt.Sprintf("0,%v %s %v,%v", sparkH, sparkPoints(byYear, first, last), sparkW, sparkH)
}

// sparkLabel says in words what the line shows.
func sparkLabel(noun string, byYear map[int]int, first, last int) string {
	y, n := peak(byYear)
	return fmt.Sprintf("%s per year, %d to %d, most in %d with %d", noun, first, last, y, n)
}

// bar sets --bar to a row's share of the longest, for charts.css to draw.
func bar(count, longest int) map[string]string {
	if longest <= 0 {
		return map[string]string{"--bar": "0%"}
	}
	return map[string]string{"--bar": strconv.Itoa(count*100/longest) + "%"}
}

// TopBodies ranks the animals seen most often.
func TopBodies(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, func(s sighting.Sighting) string {
		return displayName(s.Species)
	}), n)
}

// BusiestYear is the year with the most sightings.
func BusiestYear(byYear map[int]int) int {
	y, _ := peak(byYear)
	return y
}
