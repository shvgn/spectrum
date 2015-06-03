// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015
package spectrum

import (
	"log"
	"sort"

	"github.com/ready-steady/spline"
)

type dataSorterX [][2]float64
type dataSorterY [][2]float64

func (d dataSorterX) Len() int           { return len(d) }
func (d dataSorterY) Len() int           { return len(d) }
func (d dataSorterX) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d dataSorterY) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d dataSorterX) Less(i, j int) bool { return d[i][0] < d[j][0] }
func (d dataSorterY) Less(i, j int) bool { return d[i][1] < d[j][1] }

func (s *Spectrum) SortByX() { sort.Sort(dataSorterX(s.data)) }
func (s *Spectrum) SortByY() { sort.Sort(dataSorterY(s.data)) }

// Choose borders X1 and X2 to cut the spectrum X range
func (s *Spectrum) Cut(x1, x2 float64) {
	if x1 > x2 {
		log.Fatal("X1 cannot be bigger than X2 in Filter() method")
	}
	var i1, i2 int
	startIndexKnown := false
	for i, p := range s.data {
		if !startIndexKnown && p[0] >= x1 {
			i1 = i
			startIndexKnown = true
		}
		if startIndexKnown && p[0] > x2 {
			i2 = i
			break
		}
	}
	s.data = s.data[i1:i2]
}

// Modifies X with arbitrary function, ensures sorted X after the modification
func (s *Spectrum) ModifyX(f func(x float64) float64) {
	for i := range s.data {
		s.data[i][0] = f(s.data[i][0])
	}
	s.SortByX()
}

// Modifies Y with arbitrary function
func (s *Spectrum) ModifyY(f func(x float64) float64) {
	for i := range s.data {
		s.data[i][1] = f(s.data[i][1])
	}

}

// Arithmetic operation for two floats
func arithOpFunc(sym rune) func(float64, float64) float64 {
	switch sym {
	case '+':
		return func(a, b float64) float64 { return a + b }
	case '-':
		return func(a, b float64) float64 { return a - b }
	case '*':
		return func(a, b float64) float64 { return a * b }
	case '/':
		return func(a, b float64) float64 { return a / b }
	default:
		log.Fatal("Unknown arithmetic operation")
	}
	return nil
}

// Function for arithmetic operation over two spectra. If X values do not
// coincide, interpolation of the second specrum is used
func doArithOperation(s1, s2 *Spectrum, op rune) error {
	ol := newOverlap(s1, s2)
	if ol.err != nil {
		return ol.err
	}

	f := arithOpFunc(op)
	l1 := ol.i1r - ol.i1l + 1
	l2 := ol.i2r - ol.i2l + 1
	data := make([][2]float64, l1)

	// First we shall see if X axes coincise and spectra can be operated.
	// If l1 == l2 then X1 and X2 must coincise but they can still be shifted
	// in their indexes
	if l1 == l2 {
		for j := 0; j < l1; j++ {
			data[j][0] = s1.data[j+ol.i1l][0]                          // x
			data[j][1] = f(s1.data[j+ol.i1l][1], s2.data[j+ol.i2l][1]) // y
		}
		s1.data = data
		return nil
	}

	// If X ranges do not coincise Y2 is reduced to the interpolated over X1

	// Filling slices #1
	x1slc := make([]float64, l1)
	y1slc := make([]float64, l1)
	for j := ol.i1l; j <= ol.i1r; j++ {
		x1slc[j-ol.i1l] = s1.data[j][0]
		y1slc[j-ol.i1l] = s1.data[j][1]
	}

	// Filling slices #2
	x2slc := make([]float64, l2)
	y2slc := make([]float64, l2)
	for j := ol.i2l; j <= ol.i2r; j++ {
		x2slc[j-ol.i2l] = s2.data[j][0]
		y2slc[j-ol.i2l] = s2.data[j][1]
	}

	// Cubic spline
	var cb *spline.Cubic
	cb = spline.NewCubic(x2slc, y2slc)
	y2slc = cb.Evaluate(x1slc)

	for i, x := range x1slc {
		data[i] = [2]float64{x, f(y1slc[i], y2slc[i])}
	}
	s1.data = data
	return nil
}

// Adds spectrum to the current one
func (s *Spectrum) Add(ss *Spectrum) {
	doArithOperation(s, ss, '+')
}

// Subtract spectrum from the current one
func (s *Spectrum) Subtract(ss *Spectrum) {
	doArithOperation(s, ss, '-')
}

// Multiply spectrum by the current one
func (s *Spectrum) Multiply(ss *Spectrum) {
	doArithOperation(s, ss, '*')
}

// Divide spectrum by the current one
func (s *Spectrum) Divide(ss *Spectrum) {
	doArithOperation(s, ss, '/')
}
