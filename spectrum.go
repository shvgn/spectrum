// The code is provided "as is" without any warranty and shit.
// You are free to copy, use and redistribute the code as you wish.
//
// Evgenii Shevchenko
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

type Spectrum struct {
	data [][2]float64
	meta map[string]string
}

// --------------------------------------------------------------------------
// A Spectrum constructor
func NewSpectrum(capacity int) *Spectrum {
	spec := Spectrum{make([][2]float64, capacity), make(map[string]string)}
	return &spec
}

// Reader for TSV files
func NewTSVReader(r io.Reader) *csv.Reader {
	csvr := csv.NewReader(r)
	csvr.Comma = '\t'
	csvr.Comment = '#'
	csvr.FieldsPerRecord = 0
	csvr.LazyQuotes = false
	csvr.TrimLeadingSpace = true
	return csvr
}

// Reading from TSV file, cols must contain numbers of columns to take into
// account. If cols consists of one integer, the integer value is taken as Y.
// If cols consists of two integers, they are taken as numbers of X and Y
// columns in the passed TSV. If cols is not passed, then X defaults to 1 and Y
// defaults to 2 as in ordinar 2-column ASCII TSV file. If len(cols)>2, the
// error is returned
func ReadFromTSV(r io.Reader, xcol, ycol int) (*Spectrum, error) {

	tsvreader := NewTSVReader(r)
	records, err := tsvreader.ReadAll()
	if err != nil {
		return nil, err
	}

	data := make([][2]float64, len(records))
	meta := make(map[string]string)
	entry := [2]string{}
	var i int = 0
	for _, e := range records {
		entry[0] = e[xcol]
		entry[1] = e[ycol]
		x, xerr := parseFloat(entry[0])
		y, yerr := parseFloat(entry[1])
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

// Parse a slice of srings and make them float64s
func parseFloat(str string) (float64, error) {
	newstr := strings.TrimSpace(str)
	newstr = strings.Replace(newstr, ",", ".", 1)
	return strconv.ParseFloat(newstr, 64)
}

// --------------------------------------------------------------------------
// String representation // FIXME the order of headers must not be randomized
func (spec *Spectrum) String() string {
	var buf bytes.Buffer
	var lines []string

	for xstring, ystring := range spec.meta {
		lines = append(lines, fmt.Sprintf("%s\t%s\n", xstring, ystring))
	}

	sort.Strings(lines) // For consistent order  XXX WHY?
	for _, line := range lines {
		buf.WriteString(line)
	}

	for _, xy := range spec.data {
		buf.WriteString(fmt.Sprintf("%f\t%f\n", xy[0], xy[1]))
	}
	return buf.String()
}

// --------------------------------------------------------------------------
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

// --------------------------------------------------------------------------
// Parser for the the read data
func parseSpectrum(data []byte, xcol, ycol int) *Spectrum {

	lines := strings.Split(string(data), "\n")
	datamap := make(map[float64]float64)
	metamap := make(map[string]string)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}

		fields := strings.Fields(line)

		x, errx := parseFloat(fields[xcol])
		y, erry := parseFloat(fields[ycol])
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
	return spec
}

// --------------------------------------------------------------------------
// Read data file and return a new spectrum
func SpectrumFromFile(fname string) (*Spectrum, error) {
	spec, err := ReadFromFile(fname)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// --------------------------------------------------------------------------
// Read data into an existing spectrum
func ReadFromFile(fname string, cols ...int) (*Spectrum, error) {
	// So we received cols. Now we decide which numbers of columns to take into
	// the spectrum
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
	xcol--
	ycol--
	if xcol < 0 || ycol < 0 {
		return nil, errors.New(fmt.Sprintf("Column indexes mut be positive, received xcol=%d ycol=%d", xcol+1, ycol+1))
	}

	fi, err := os.Open(fname)
	if err != nil {
		fmt.Printf("Cannot open file <%s>\n", fname)
		return nil, err
	}
	spec, err := ReadFromTSV(fi, xcol, ycol)
	if err != nil {
		// Try to parse in another way
		var rawdata []byte
		rawdata, err = ioutil.ReadFile(fname)
		if err != nil {
			fmt.Println("Cannot read file", fname, err.Error())
			return nil, err
		}
		spec = parseSpectrum(rawdata, xcol, ycol)
	}
	return spec, err
}

//
// --------------------------------------------------------------------------
// Write spectrum to a file
func (spec *Spectrum) WriteToFile(file string) error {
	err := ioutil.WriteFile(file, []byte(spec.String()), 0600)
	if err != nil {
		fmt.Println("Cannot write to file", file, err.Error())
		return err
	}
	return nil
}

// --------------------------------------------------------------------------
// Calculate noise level of the spectrum according to its minimum Y values distribution
func (spec *Spectrum) Noise() float64 {
	fmt.Println("WARNING! Noise() method is no implemented yet.")
	return 0.0
}

// --------------------------------------------------------------------------
// Calculate are under the spectrum
func (spec *Spectrum) Area() float64 {
	fmt.Println("WARNING! Area() method is no implemented yet.")
	return 0.0

}

// --------------------------------------------------------------------------
// Choose borders X1 and X2 to cut the spectrum X range
func (spec *Spectrum) Cut(x1, x2 float64) {
	fmt.Println("WARNING! Cut() method is no implemented yet.")

}

// --------------------------------------------------------------------------
// Modify X with arbitrary function
func (spec *Spectrum) ModifyX(modifier func(x float64) float64) {
	fmt.Println("WARNING! ModifyX() method is no implemented yet.")

}

// --------------------------------------------------------------------------
// Modify Y with arbitrary function
func (spec *Spectrum) ModifyY(modifier func(x float64) float64) {
	fmt.Println("WARNING! ModifyY() method is no implemented yet.")

}

// --------------------------------------------------------------------------
// Take position of the spectrum maximum
func (spec *Spectrum) MaxY() (float64, float64) {
	fmt.Println("WARNING! MaxY() method is no implemented yet.")
	return 0.0, 0.0
}

// --------------------------------------------------------------------------
// --------------------------------------------------------------------------
// TO IMPLEMENT
// Spectra multiplication and division, merging/averaging
// Spectra addition and subtraction
// Smoothing
// Splicing
// Peaks analisys: FWHM, Gauss/Lorenz fitting... maybe
