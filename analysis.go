// Package xy is a simple library for manipulation of X,Y data

package xy

import (
	"log"
	"math"
	"sort"
)

const (
	precisionOrder int = 4 // Precision for rounding and calculations
)

// Rounding for float64
func roundFloat64(f float64, prec int) float64 {
	shift := math.Pow(10, float64(prec))
	return math.Floor(f*shift+0.5) / shift
}

// Noise naively calculates noise level of the dataset according to its minimum Y values
// distribution. The more noise in the dataset the better it is calculated
func (s *XY) Noise() float64 {
	ydist := map[float64]int{}
	// Choose the precision
	_, ymax, eps := s.MaxAndEps()
	prec := int(math.Log10(ymax / eps))
	if prec < precisionOrder {
		prec = precisionOrder
	}
	// Populate the distribution
	for _, p := range s.data {
		y := roundFloat64(p[1], prec)
		if _, found := ydist[y]; found {
			ydist[y]++
		} else {
			ydist[y] = 1
		}
	}
	distdata := make([][2]float64, len(ydist))
	var i int
	// y0 stands for the y that has the biggest counter value c0,
	// namely the peak of the noise distribution
	var c0, y0 float64
	for y, c := range ydist {
		cf := float64(c)
		if cf > c0 {
			c0 = cf
			y0 = y
		}
		distdata[i][0] = y
		distdata[i][1] = cf
		i++
	}
	// Now we have a distribution data [Y,C] sorted by Y's. We are going to find the
	// full width at half-maximum (l for left and r fot right) and the counts
	// peak center is what we are looking for.
	sort.Sort(dataSorterX(distdata))
	cl := c0
	cr := cl
	var yl, yr float64
	for _, p := range distdata {
		y, c := p[0], p[1]
		zeroSeek := math.Abs(c0/2 - c)
		if y < y0 {
			if math.Abs(cl-c0/2) > zeroSeek {
				cl = c
				yl = y
			}
		}
		if y > y0 {
			if math.Abs(cr-c0/2) > zeroSeek {
				cr = c
				yr = y
			}
		}
	}
	// Return the center of the FWHM of the [Y,C] peak
	return 0.5 * (yl + yr)
}

// Calculate area under the spectrum with the trapezoidal method
func (s *XY) Area() float64 {
	l := len(s.data)
	data := make([][2]float64, l)
	copy(s.data, data)
	sort.Sort(dataSorterX(data))
	var area float64
	for i := 0; i < l-1; i++ {
		x1 := data[i][0]
		x2 := data[i+1][0]
		y1 := data[i][1]
		y2 := data[i+1][1]
		area += (x2 - x1) * (y1 + y2)
	}
	return area / 2

}

// Get X and Y of the first spectrum point
func (s *XY) FirstPoint() (float64, float64) {
	return s.data[0][0], s.data[0][1]
}

// Get X and Y of the last spectrum point
func (s *XY) LastPoint() (float64, float64) {
	k := len(s.data) - 1
	return s.data[k][0], s.data[k][1]
}

// The spectrum maximum point X Y
func (s *XY) MaxY() (float64, float64) {
	x, y, _ := s.MaxAndEps()
	return x, y
}

// MaxAndEps returns the maximum point credentials
func (s *XY) MaxAndEps() (xmax, ymax, eps float64) {
	l := len(s.data)
	if l == 0 {
		log.Fatal("Empty data")
	}
	// Make a copy to modify freely
	sortedByY := make([][2]float64, l)
	copy(sortedByY, s.data)
	sort.Sort(dataSorterY(sortedByY))
	xmax = sortedByY[l-1][0]
	ymax = sortedByY[l-1][1]
	eps = ymax
	// Search for the mimimum step for Y values
	for _, point := range sortedByY {
		y := math.Abs(point[1])
		if y < eps && y != 0 {
			eps = y
		}
	}
	return xmax, ymax, eps
}

// FWHM is full width at half-maximum near the given X
// Not implemented yet
func (s *XY) FWHM(x float64) float64 {
	// Here I must calculate derivatives with noise-ignorant method such as
	// Savitsky-Golay filter or Holoborodko's method
	return 0.0
}
