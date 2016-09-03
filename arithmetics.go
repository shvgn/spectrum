package xy

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

// SortByX sorts the data by X in scending order
func (s *XY) SortByX() { sort.Sort(dataSorterX(s.data)) }

// SortByY sorts the data by Y in scending order
func (s *XY) SortByY() { sort.Sort(dataSorterY(s.data)) }

// Cut chooses borders X1 and X2 to cut the spectrum X range
func (s *XY) Cut(x1, x2 float64) {
	// FIXME what about one-side cut?
	i1, i2, err := FindBordersIndexes(s.data, x1, x2)
	if err != nil {
		log.Fatal("x1 cannot be bigger than x2")
	}
	s.data = s.data[i1 : i2+1]
}

// ModifyX applies arbitrary function to X and ensures sorted X after the modification
func (s *XY) ModifyX(f func(x float64) float64) {
	for i := range s.data {
		s.data[i][0] = f(s.data[i][0])
	}
	s.SortByX()
}

// ModifyY applies arbitrary function to Y
func (s *XY) ModifyY(f func(x float64) float64) {
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

// Function for arithmetic operation over Y's of two datasets (ds1 and ds2). If their
// X values don't coincide, this function will interpolatie Y of ds2 over X of ds1
// Thus the result always inherits X of ds1.
func doArithOperation(ds1, ds2 *XY, op rune) error {
	ol, err := newOverlap(ds1, ds2)

	if err != nil {
		log.Println("Error in overlap")
		log.Println("Headers of s1:")
		log.Println(ds1.meta)
		log.Println("Headers of s2:")
		log.Println(ds2.meta)
		return err
	}

	f := arithOpFunc(op)
	l1 := ol.i1r - ol.i1l + 1
	l2 := ol.i2r - ol.i2l + 1

	data := make([][2]float64, l1, l1) // The result size will be the one of ds1
	data1 := ds1.data[ol.i1l : ol.i1r+1]
	data2 := ds2.data[ol.i2l : ol.i2r+1]

	// First we shall see if X axes coincise and spectra can be operated. This
	// is useful for data obtained on one setup. If l1 == l2 then X1 and X2
	// must coincise but they can still be shifted in their indexes
	if l1 == l2 {
		ok := true // The things-go-fine indicator
		for j, p := range data1 {
			x1, y1 := p[0], p[1]
			x2, y2 := data2[j][0], data2[j][1]
			if x1 != x2 {
				// They don't coincise. Clear the data and go another way
				data = make([][2]float64, 0, l1)
				ok = false
				break
			}
			data[j] = [2]float64{x1, f(y1, y2)}
		}
		if ok {
			ds1.data = data // Here we cut s1
			return nil
		}
	}

	// If X ranges do not coincise Y2 is reduced to the interpolated over X1
	// Filling slices #1
	xa1 := make([]float64, l1, l1)
	ya1 := make([]float64, l1, l1)
	for i, p := range data1 {
		xa1[i] = p[0]
		ya1[i] = p[1]
	}

	// Filling slices #2
	xa2 := make([]float64, l2, l2)
	ya2 := make([]float64, l2, l2)
	for i, p := range data2 {
		xa2[i] = p[0]
		ya2[i] = p[1]
	}

	// Cubic spline
	var cb *spline.Cubic
	cb = spline.NewCubic(xa2, ya2)
	ya2 = cb.Evaluate(xa1)

	data = make([][2]float64, l1, l1) // Cleaned data
	for i, x := range xa1 {
		data[i] = [2]float64{x, f(ya1[i], ya2[i])}
	}
	ds1.data = data
	return nil
}

// Add spectrum to the current one
func (s *XY) Add(ss *XY) error {
	return doArithOperation(s, ss, '+')
}

// Subtract spectrum from the current one
func (s *XY) Subtract(ss *XY) error {
	return doArithOperation(s, ss, '-')
}

// Multiply spectrum by the current one
func (s *XY) Multiply(ss *XY) error {
	return doArithOperation(s, ss, '*')
}

// Divide spectrum by the current one
func (s *XY) Divide(ss *XY) error {
	return doArithOperation(s, ss, '/')
}
