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
	ol.xl = math.Max(x1l, x2l)        // Maximum of the minima
	ol.xr = math.Min(x1r, x2r)        // Minimum of the maxima
	if ol.xl > ol.xr {
		ol.xl, ol.xr = 0.0, 0.0
		ol.err = errors.New("X ranges do not overlap")
		return ol
	}
	ol.i1l, ol.i1r, _ = FindBordersIndexes(s1.data, ol.xl, ol.xr)
	ol.i2l, ol.i2r, _ = FindBordersIndexes(s2.data, ol.xl, ol.xr)
	return ol
}

// Find indices of data borders. The data is assumed to be sorted. xLeft and
// xRight shall not be swapped automatically. If x1 is less then X minimum, i1
// shall be 0, and in the same manner i2 will be the last index of the data if
// x2 is bigger than X maximum.
func FindBordersIndexes(data [][2]float64, x1, x2 float64) (int, int, error) {
	var err error
	if x1 >= x2 {
		return -1, -1, errors.New("FindBorderIndexes: x1 >= x2")
	}
	var i1, i2 int
	var found1, found2 bool
	xmin := data[0][0]
	xmax := data[len(data)-1][0]
	if x1 <= xmin {
		i1, found1 = 0, true
	}
	if x2 >= xmax {
		i2, found2 = len(data)-1, true
		if found1 {
			return i1, i2, nil
		}
	}
	for i, p := range data {
		if !found1 && p[0] >= x1 {
			i1, found1 = i, true
			if found2 {
				break
			}
		}
		if found1 && !found2 && p[0] > x2 {
			i2, found2 = i-1, true
			break
		}
	}
	return i1, i2, err
}
