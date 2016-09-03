// Package xy is a simple library for manipulation of columned X,Y data
package xy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
)

const (
	headerDelimiter = "\t"
	numberDelimiter = "\t"
)

/* dataset type
 * data is [][2]float64, like
 * [
 *    [x1, y1],
 *    [x2, y2],
 *    [x3, y3]
 * ]
 * meta is a map of strings
 */

// XY represents the `data` containig points (x,y) as slices of floats and metainfo
// in the `meta` strings map
type XY struct {
	data [][2]float64
	meta map[string]string
}

// NewXY creates XY instance for the given capacity
func NewXY(capacity int) *XY {
	spec := XY{make([][2]float64, capacity), make(map[string]string)}
	return &spec
}

// Len returns the number of points
func (s *XY) Len() int {
	return len(s.data)
}

// String representation
func (s *XY) String() string {
	var buf bytes.Buffer
	var lines []string

	// FIXME the order of headers must not be randomized
	for name, value := range s.meta {
		lines = append(lines, fmt.Sprintf("%s%s%s\n", name, headerDelimiter, value))
	}

	sort.Strings(lines) // For consistent order  XXX WHY?
	for _, line := range lines {
		buf.WriteString(line)
	}

	for _, xy := range s.data {
		buf.WriteString(fmt.Sprintf("%f%s%f\n", xy[0], numberDelimiter, xy[1]))
	}
	return buf.String()
}

// WriteToFile writes the data to a file (FIXME must take writer)
func (s *XY) WriteToFile(file string) error {
	err := ioutil.WriteFile(file, []byte(s.String()), 0600)
	if err != nil {
		return err
	}
	return nil
}

/* --------------------------------------------------------------------------
TODO:
    Smoothing
    Splicing
    Peaks analisys: FWHM, Gauss/Lorenz fitting... maybe
*/
