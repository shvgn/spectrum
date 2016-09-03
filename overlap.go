// Package xy is a simple library for manipulation of X,Y data
package xy

import (
	"errors"
	"math"
)

type overlap struct {
	xl, xr             float64
	i1l, i1r, i2l, i2r int // Indices of the boundary elements
}

// newOverlap checks whether two spectra overlap. If they do, the function
// returns an overlap structure or returns an error otherwise
func newOverlap(s1, s2 *XY) (ol *overlap, err error) {
	ol = new(overlap)

	x1l := s1.data[0][0] // x1 left
	x2l := s2.data[0][0] // x2 left

	x1r := s1.data[len(s1.data)-1][0] // x1 right
	x2r := s2.data[len(s2.data)-1][0] // x2 right

	ol.xl = math.Max(x1l, x2l) // Maximum of minima
	ol.xr = math.Min(x1r, x2r) // Minimum of maxima

	if ol.xl > ol.xr {
		return nil, errors.New("X ranges do not overlap")
	}

	ol.i1l, ol.i1r, err = FindBordersIndexes(s1.data, ol.xl, ol.xr)
	if err != nil {
		return nil, err
	}

	ol.i2l, ol.i2r, err = FindBordersIndexes(s2.data, ol.xl, ol.xr)
	if err != nil {
		return nil, err
	}
	return ol, nil
}

// FindBordersIndexes finds borders of data. The data is assumed to be sorted. If x1 is
// less then X minimum, i1 shall be 0, and in the same manner i2 will be the
// last index of the data if x2 is bigger than X maximum.
func FindBordersIndexes(data [][2]float64, x1, x2 float64) (i1, i2 int, err error) {
	var found1, found2 bool

	if x1 >= x2 {
		err = errors.New("Invalid border values: x1 >= x2")
		return -1, -1, err
	}

	xmin := data[0][0]
	xmax := data[len(data)-1][0]

	// First we check whether the values exceed data borders
	if x1 <= xmin {
		i1, found1 = 0, true
	}
	if x2 >= xmax {
		i2, found2 = len(data)-1, true
		if found1 {
			return i1, i2, nil
		}
	}

	// Now search
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

	// I dont't how this can happen
	if !found1 || !found2 {
		err = errors.New("Cannot find border values")
		i1, i2 = -1, -1
	}
	return i1, i2, err
}
