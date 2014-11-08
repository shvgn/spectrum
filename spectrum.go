// The code is provided "as is" without any warranty and shit.
// You are free do anything you want with it.
//
// Evgenii Shevchenko a.k.a @shvgn
// 2014

package spectrum

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type Spectrum struct {
	// data    map[float64]float64 // X -> Y
	x       []float64         // X
	y       []float64         // Y
	headers map[string]string // arbitrary metainfo
}

// --------------------------------------------------------------------------
// A Spectrum constructor
func NewSpectrum(capacity int) *Spectrum {
	// Map
	// spec := Spectrum{make(map[float64]float64), make(map[string]string)}

	// Map and x and y arrays
	// spec := Spectrum{	make(map[float64]float64), make([]float64, 100),
	// 						make([]float64, 100), make(map[string]string)}

	// x and y arrays
	spec := Spectrum{make([]float64, capacity), make([]float64, capacity), make(map[string]string)}
	return &spec
}

// --------------------------------------------------------------------------
// Make a string representation of the Spectrum
func (spec *Spectrum) String() string {
	var buf bytes.Buffer
	var lines []string

	for header, value := range spec.headers {
		lines = append(lines, fmt.Sprintf("%s\t%s\n", header, value))
	}
	sort.Strings(lines) // For consistent order
	for _, line := range lines {
		buf.WriteString(line)
	}

	for i, x := range spec.x {
		buf.WriteString(fmt.Sprintf("%f\t%f\n", x, spec.y[i]))
	}
	return buf.String()
}

// --------------------------------------------------------------------------
// Parse a header string
func parseHeader(line string) (string, string) {
	header := strings.TrimSpace(line)
	index := strings.IndexAny(header, ":=")
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
func (spec *Spectrum) parseSpectrum(data []byte) {

	lines := strings.Split(string(data), "\n")
	datamap := make(map[float64]float64)
	headersmap := make(map[string]string)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}

		fields := strings.Fields(line)

		// FIXME there can be more valid columns, especially useless ones such
		// as data from WinSpec software where the second column contains frame
		// number which is useless for a spectrum:
		// 238	  1    5643.43

		if len(fields) != 2 {
			// Not a point X,Y, must be a header
			header, value := parseHeader(line)
			headersmap[header] = value
			continue
		}

		x, errx := strconv.ParseFloat(fields[0], 64)
		if errx != nil {
			// This line seems to be a header, too
			header, value := parseHeader(line)
			headersmap[header] = value
			continue
		}

		y, erry := strconv.ParseFloat(fields[1], 64)
		if erry != nil {
			// This line seems to be unknown stuff, since x is float
			fmt.Println("Cannot parse line", line)
			continue
		}

		// Ok here, x and y are valid float64's
		datamap[x] = y
	}

	spec.headers = headersmap

	// Make sorted slices of x and y
	length := len(datamap)
	xrange := make([]float64, length)
	yrange := make([]float64, length)
	index := 0

	for k, _ := range datamap {
		xrange[index] = k
		index++
	}

	sort.Float64s(xrange)

	for i, x := range xrange {
		yrange[i] = datamap[x]
	}

	spec.x = xrange
	spec.y = yrange
}

// --------------------------------------------------------------------------
// Read data file and return a new spectrum
func SpectrumFromFile(file string) (*Spectrum, error) {
	spec := &Spectrum{}
	err := spec.ReadFromFile(file)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// --------------------------------------------------------------------------
// Read data into an existing spectrum
func (spec *Spectrum) ReadFromFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Cannot read file", file, err.Error())
		return err
	}
	spec.parseSpectrum(data)
	return nil
}

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
// Peaks analisys: FWHM, Gauss/Lorenz fitting
