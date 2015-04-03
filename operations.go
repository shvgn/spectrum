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
)

// type spectrumSorterX Spectrum
// type spectrumSorterY Spectrum

// Sort interface
// func (ss spectrumSorterX) Len() int           { return len(ss.data) }
// func (ss spectrumSorterY) Len() int           { return len(ss.data) }
// func (ss spectrumSorterX) Swap(i, j int)      { ss.data[i], ss.data[j] = ss.data[j], ss.data[i] }
// func (ss spectrumSorterY) Swap(i, j int)      { ss.data[i], ss.data[j] = ss.data[j], ss.data[i] }
// func (ss spectrumSorterX) Less(i, j int) bool { return ss.data[i][0] < ss.data[j][0] }
// func (ss spectrumSorterY) Less(i, j int) bool { return ss.data[i][1] < ss.data[j][1] }

// func (s *Spectrum) sort()    { sort.Sort(spectrumSorterX(*s)) }
// func (s *Spectrum) sortByY() { sort.Sort(spectrumSorterY(*s)) }

type dataSorterX [][2]float64
type dataSorterY [][2]float64

func (d dataSorterX) Len() int           { return len(d) }
func (d dataSorterY) Len() int           { return len(d) }
func (d dataSorterX) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d dataSorterY) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d dataSorterY) Less(i, j int) bool { return d[i][1] < d[j][1] }
func (d dataSorterX) Less(i, j int) bool { return d[i][0] < d[j][0] }

func (s *Spectrum) SortByX() { sort.Sort(dataSorterX(s.data)) }
func (s *Spectrum) SortByY() { sort.Sort(dataSorterY(s.data)) }

// Choose borders X1 and X2 to cut the spectrum X range
func (s *Spectrum) Cut(x1, x2 float64) {
	if x1 > x2 {
		log.Fatal("X1 cannot be bigger than X2 in Filter() method")
	}
	filtered := [][2]float64{}
	for _, p := range s.data {
		if p[0] >= x1 && p[0] <= x2 {
			filtered = append(filtered, p)
		}
	}
	s.data = filtered

}

// Modify X with arbitrary function
func (s *Spectrum) ModifyX(f func(x float64) float64) {
	for i := range s.data {
		s.data[i][0] = f(s.data[i][0])
	}
	s.SortByX()
}

// Modify Y with arbitrary function
func (s *Spectrum) ModifyY(f func(x float64) float64) {
	for i := range s.data {
		s.data[i][1] = f(s.data[i][1])
	}

}
