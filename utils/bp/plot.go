/*
	pb – Braille Patterns

	Simple way to represent numeric lines like a plot in text.

	TODO: extract to separate repo
*/

package bp

import (
	"fmt"
	"math"
	"strings"

	"github.com/egregors/hk/log"
)

// ⣿⣶⣤⣀ – ok
// ⣾⣷⣴⣦⣠⣄ - ok
// ⣼⣧⣸⣇⣰⣆ - ok
// ⢸⢰⢠⢀⡇⡆⡄⡀ – ok
const (
	m44 = "⣿"
	m33 = "⣶"
	m22 = "⣤"
	m11 = "⣀"

	m34 = "⣾"
	m43 = "⣷"
	m23 = "⣴"
	m32 = "⣦"
	m12 = "⣠"
	m21 = "⣄"

	m24 = "⣼"
	m42 = "⣧"
	m14 = "⣸"
	m41 = "⣇"
	m13 = "⣰"
	m31 = "⣆"

	m40 = "⡇"
	m30 = "⡆"
	m20 = "⡄"
	m10 = "⡀"

	m04 = "⢸"
	m03 = "⢰"
	m02 = "⢠"
	m01 = "⢀"
	m00 = "⠀"
)

var bps = [5][]rune{
	[]rune(m00 + m01 + m02 + m03 + m04),
	[]rune(m10 + m11 + m12 + m13 + m14),
	[]rune(m20 + m21 + m22 + m23 + m24),
	[]rune(m30 + m31 + m32 + m33 + m34),
	[]rune(m40 + m41 + m42 + m43 + m44),
}

func SimplePlot(size int, data []float64) string {
	if len(data) == 0 {
		return ""
	}
	lo, hi := minMax(data)
	log.Debg.Printf("hi: %.2f lo: %.2f\n", hi, lo)
	maxRange := math.Abs(lo) + math.Abs(hi)
	dot := maxRange / (float64(size * 4))
	log.Debg.Printf("1 dot of 4 lines fits: %.3f", dot)
	log.Debg.Println("high:", maxRange)

	ps := pairs(data)
	log.Debg.Println("all ps:", ps)

	// make all values positive
	if lo < 0 {
		delta := lo * -1
		for _, p := range ps {
			p[0], p[1] = p[0]+delta, p[1]+delta
		}
	}
	log.Debg.Println("pos ps:", ps)

	// normalize values (to show only dynamic range)
	for _, p := range ps {
		p[0], p[1] = max(0, p[0]-lo)*float64(size*4), max(0, p[1]-lo)*float64(size*4)
	}
	log.Debg.Println("nor ps:", ps)

	// converts values to amount of dots
	ps = dots(ps, dot)
	log.Debg.Println("dot ps:", ps)

	// add value bounds
	vis := fmt.Sprintf("%.2f\n%s%.2f", hi, render(size, ps), lo)

	return vis
}

func render(size int, ps [][]float64) string {
	plot := make([][]string, size)
	for r := range plot {
		plot[r] = make([]string, len(ps))
	}

	for c, p := range ps {
		fst, snd := int(p[0]), int(p[1])
		for r := len(plot) - 1; r >= 0; r, fst, snd = r-1, fst-4, snd-4 {
			currFst, currSnd := max(0, fst), max(0, snd)
			currFst, currSnd = min(4, currFst), min(4, currSnd)
			plot[r][c] = string(bps[currFst][currSnd])
		}
	}

	sb := strings.Builder{}
	for _, line := range plot {
		sb.WriteString(strings.Join(line, "") + "\n")
	}

	return sb.String()
}

func dots(ps [][]float64, dot float64) [][]float64 {
	for _, p := range ps {
		p[0], p[1] = math.Round(p[0]/dot), math.Round(p[1]/dot)
	}

	return ps
}

func pairs(xs []float64) [][]float64 {
	if len(xs)%2 != 0 {
		xs = append(xs, 0)
	}

	if len(xs) <= 2 {
		return [][]float64{xs}
	}

	return append([][]float64{xs[:2]}, pairs(xs[2:])...)
}

func minMax(xs []float64) (float64, float64) {
	minimum := math.MaxFloat64
	maximum := math.SmallestNonzeroFloat64
	for _, x := range xs {
		if x < minimum {
			minimum = x
		}
		if x > maximum {
			maximum = x
		}
	}

	return minimum, maximum
}
