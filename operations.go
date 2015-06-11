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
	// FIXME what about one-side cut?
	i1, i2, err := FindBordersIndexes(s.data, x1, x2)
	if err != nil {
		log.Fatal("X1 cannot be bigger than X2 in Cut() method")
	}
	s.data = s.data[i1 : i2+1]
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
// coincide the interpolation of the second specrum is used
func doArithOperation(s1, s2 *Spectrum, op rune) error {
	ol := newOverlap(s1, s2)
	if ol.err != nil {
		log.Println("Error in overlap")
		log.Println("Headers os s1:")
		log.Println(s1.meta)
		log.Println("Headers os s2:")
		log.Println(s2.meta)
		return ol.err
	}

	f := arithOpFunc(op)
	l1 := ol.i1r - ol.i1l + 1
	l2 := ol.i2r - ol.i2l + 1

	data := make([][2]float64, 0, l1) // The result size will be the one of s1
	data1 := s1.data[ol.i1l : ol.i1r+1]
	data2 := s2.data[ol.i2l : ol.i2r+1]

	// First we shall see if X axes coincise and spectra can be operated. This
	// is useful for data obtained on one setup. If l1 == l2 then X1 and X2
	// must coincise but they can still be shifted in their indexes
	if l1 == l2 {
		ok := true
		for j, p := range data1 {
			x1, y1 := p[0], p[1]
			x2, y2 := data2[j][0], data2[j][1]
			if x1 != x2 {
				// They don't coincise. Clear the data and go another way
				data = make([][2]float64, 0, l1)
				ok = false
				break
			}
			data = append(data, [2]float64{x1, f(y1, y2)})
		}
		if ok {
			s1.data = data // Here we cut s1
			return nil
		}
	}

	// If X ranges do not coincise Y2 is reduced to the interpolated over X1
	// Filling slices #1
	xa1 := make([]float64, 0, l1)
	ya1 := make([]float64, 0, l1)
	for _, p := range data1 {
		xa1 = append(xa1, p[0])
		ya1 = append(ya1, p[1])
	}

	// Filling slices #2
	xa2 := make([]float64, 0, l2)
	ya2 := make([]float64, 0, l2)
	for _, p := range data2 {
		xa2 = append(xa2, p[0])
		ya2 = append(ya2, p[1])
	}

	// Cubic spline
	var cb *spline.Cubic
	cb = spline.NewCubic(xa2, ya2)
	ya2 = cb.Evaluate(xa1)

	for i, x := range xa1 {
		data = append(data, [2]float64{x, f(ya1[i], ya2[i])})
	}
	s1.data = data
	return nil
}

// Adds spectrum to the current one
func (s *Spectrum) Add(ss *Spectrum) error {
	return doArithOperation(s, ss, '+')
}

// Subtract spectrum from the current one
func (s *Spectrum) Subtract(ss *Spectrum) error {
	return doArithOperation(s, ss, '-')
}

// Multiply spectrum by the current one
func (s *Spectrum) Multiply(ss *Spectrum) error {
	return doArithOperation(s, ss, '*')
}

// Divide spectrum by the current one
func (s *Spectrum) Divide(ss *Spectrum) error {
	return doArithOperation(s, ss, '/')
}
