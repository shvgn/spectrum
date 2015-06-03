// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code in any way you wish.
//
// Evgeny Shevchenko
// shvgn@protonmail.ch
// 2015

package spectrum

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/* Spectrum type
 * data is [][2]float64, like
 * [
 *    [x1, y1],
 *    [x2, y2],
 *    [x3, y3]
 * ]
 * meta is a map of strings
 */

const (
	commentPrefix = '#' // Lines starting with this rune will be ignored
)

type Spectrum struct {
	data [][2]float64
	meta map[string]string
}

// Constructor
func NewSpectrum(capacity int) *Spectrum {
	spec := Spectrum{make([][2]float64, capacity), make(map[string]string)}
	return &spec
}

// Number of points in the spectrum data
func (s *Spectrum) Len() int {
	return len(s.data)
}

// String representation // FIXME the order of headers must not be randomized
func (s *Spectrum) String() string {
	var buf bytes.Buffer
	var lines []string

	for xstr, ystr := range s.meta {
		lines = append(lines, fmt.Sprintf("%s\t%s\n", xstr, ystr))
	}

	sort.Strings(lines) // For consistent order  XXX WHY?
	for _, line := range lines {
		buf.WriteString(line)
	}

	for _, xy := range s.data {
		buf.WriteString(fmt.Sprintf("%f\t%f\n", xy[0], xy[1]))
	}
	return buf.String()
}

// Parse a header string
func parseHeader(line string) (string, string) {
	header := strings.TrimSpace(line)
	index := strings.IndexAny(header, "\t:=")
	if index > 0 && index < len(line) {
		value := header[index+1:]
		header = header[:index]
		return header, value
	}
	parts := strings.Fields(line)
	header = parts[0]
	value := strings.Join(parts[1:], " ")
	return header, value
}

// Read data file and return a new spectrum
func SpectrumFromFile(fname string, cols ...int) (*Spectrum, error) {
	// So we received cols. Now we decide which numbers of columns to take into
	// the spectrum. We keep numbers starting from 1 im order to print these
	// numbers in the below error if it occurs.
	var xcol, ycol int
	switch len(cols) {
	case 0:
		xcol = 1
		ycol = 2
	case 1:
		xcol = 1
		ycol = cols[0]
	case 2:
		xcol = cols[0]
		ycol = cols[1]
	default:
		log.Fatal("Incorrect number of entries in ReadFromFile")
	}
	if xcol < 0 || ycol < 0 {
		return nil, errors.New(
			fmt.Sprintf(
				"Column indexes mut be positive, received xcol=%d ycol=%d",
				xcol, ycol))
	}
	xcol--
	ycol--
	// Now the columns numbers are hypothetic.

	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	// We read we the TSV reader.
	s, err := ReadFromTSV(fi, xcol, ycol)
	if err != nil {
		// Try read the spectrum in another way. Why we use TSV if this can
		// handle all cases?
		var rawdata []byte
		rawdata, err = ioutil.ReadFile(fname)
		if err != nil {
			fmt.Println("Cannot read file", fname, err.Error())
			return nil, err
		}
		s, err = parseSpectrum(rawdata, xcol, ycol)
	}
	return s, err
}

// Reader for TSV files
func NewTSVReader(r io.Reader) *csv.Reader {
	csvr := csv.NewReader(r)
	csvr.Comma = '\t'
	csvr.Comment = commentPrefix
	csvr.FieldsPerRecord = 0 // The expected number of columns is derived from the first line
	csvr.LazyQuotes = false
	csvr.TrimLeadingSpace = true
	return csvr
}

// Reading from TSV file, cols must contain numbers of columns to take into
// account. If cols consists of one integer, the integer value is taken as
// number of the Y column. If cols consists of two integers, they are taken as
// numbers of X and Y columns in the passed TSV. If cols is not passed, then X
// defaults to 1 and Y defaults to 2 as in ordinar 2-column ASCII TSV file. If
// len(cols)>2, the error is returned
func ReadFromTSV(r io.Reader, xcol, ycol int) (*Spectrum, error) {

	tsvreader := NewTSVReader(r)
	records, err := tsvreader.ReadAll() // [][]string
	if err != nil {
		return nil, err
	}

	if len(records[0]) < ycol+1 {
		return nil, csv.ErrFieldCount
	}
	data := make([][2]float64, len(records))
	meta := make(map[string]string)
	entry := [2]string{}
	var i int = 0
	for _, e := range records {
		entry[0] = e[xcol]
		entry[1] = e[ycol]
		x, xerr := ParseFloat(entry[0])
		y, yerr := ParseFloat(entry[1])
		if xerr != nil || yerr != nil {
			meta[entry[0]] = entry[1]
		} else {
			data[i] = [2]float64{x, y}
			i++
		}
	}
	spec := NewSpectrum(len(data))
	spec.data = data
	spec.meta = meta
	return spec, nil
}

// Parser for the the read data
func parseSpectrum(data []byte, xcol, ycol int) (*Spectrum, error) {

	lines := strings.Split(string(data), "\n")
	datamap := make(map[float64]float64)
	metamap := make(map[string]string)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" ||
			strings.HasPrefix(line, string(commentPrefix)) {
			continue // Skip empty lines and comments
		}

		fields := strings.Fields(line)
		if len(fields) < ycol+1 {
			return nil, csv.ErrFieldCount
		}
		x, errx := ParseFloat(fields[xcol])
		y, erry := ParseFloat(fields[ycol])
		if errx != nil || erry != nil {
			// Not a point X,Y hence must be a header
			header, value := parseHeader(line)
			metamap[header] = value
			continue
		}

		// Ok here, x and y are valid float64's
		datamap[x] = y
	}

	// Make sorted slices of x and y
	length := len(datamap)
	x_range := make([]float64, length)
	index := 0
	dataslice := make([][2]float64, length)

	for x := range datamap {
		x_range[index] = x
		index++
	}
	sort.Float64s(x_range)
	for i, x := range x_range {
		dataslice[i] = [...]float64{x, datamap[x]}
	}

	spec := NewSpectrum(length)
	spec.meta = metamap
	spec.data = dataslice
	return spec, nil
}

// Parse a float64
func ParseFloat(s string) (float64, error) {
	ns := strings.TrimSpace(s)
	ns = strings.Replace(ns, ",", ".", 1)
	return strconv.ParseFloat(ns, 64)
}

// Write spectrum to a file
func (s *Spectrum) WriteToFile(file string) error {
	err := ioutil.WriteFile(file, []byte(s.String()), 0600)
	if err != nil {
		// fmt.Println("Cannot write to file", file, err.Error())
		return err
	}
	return nil
}

// --------------------------------------------------------------------------
// TO IMPLEMENT
// Smoothing
// Splicing
// Peaks analisys: FWHM, Gauss/Lorenz fitting... maybe
