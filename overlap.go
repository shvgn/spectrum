// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package spectrum

import (
	"errors"
	"math"
)

type overlap struct {
	xl, xr             float64
	err                error
	i1l, i1r, i2l, i2r int // Indices of the boundary elements
}

// Check if two spectrum overlap with each other
// func newOverlap(s1, s2 *Spectrum) (xl, xr float64, err error) {
func newOverlap(s1, s2 *Spectrum) *overlap {
	ol := &overlap{}
	x1l := s1.data[0][0]              // x1 left
	x1r := s1.data[len(s1.data)-1][0] // x1 right
	x2l := s2.data[0][0]              // x2 left
	x2r := s2.data[len(s2.data)-1][0] // x2 right
	ol.xl = math.Max(x1l, x2l)        // Maximum of minimums
	ol.xr = math.Min(x1r, x2r)        // Minimum of maximums
	if ol.xl > ol.xr {
		ol.xl, ol.xr = 0.0, 0.0
		ol.err = errors.New("X ranges do not overlap")
		return ol
	}
	ol.i1l, ol.i1r = findBordersIndexes(s1.data, ol.xl, ol.xr)
	ol.i2l, ol.i2r = findBordersIndexes(s2.data, ol.xl, ol.xr)
	return ol
}

// Find indices of data borders
func findBordersIndexes(data [][2]float64, xLeft, xRight float64) (int, int) {
	var i1, i2 int
	foundLeft, foundRight := false, false
	for i, p := range data {
		if foundLeft && foundRight {
			break
		}
		if !foundLeft && p[0] >= xLeft {
			i1 = i
			foundLeft = true
		}
		if p[0] <= xRight {
			i2 = i
		} else {
			foundRight = true
		}
	}
	return i1, i2

}
