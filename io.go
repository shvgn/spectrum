package xy

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// FromFile reads data from the passed TSV file path and returns a new dataset
// cols is an abritrary array argument containing columns indexes starting from 1.
// If cols is not passed, then X defaults to 1 and Y defaults to 2 as in ordinar
// 2-column ASCII TSV file. If len(cols)>2, the error is returned
// The elements are interpreted as follows
// 		FromFile(fname string, xcol, ycol)
// 		FromFile(fname string, ycol)
// 		FromFile(fname string)
func FromFile(fname string, cols ...int) (*XY, error) {
	// So we received cols. Now we decide which numbers of columns to take into
	// the dataset. We keep numbers starting from 1 im order to print these
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
		log.Fatal("Incorrect number of entries in FromFile")
	}
	if xcol < 1 || ycol < 1 {
		return nil, fmt.Errorf(
			"Column indexes mut be positive, received xcol=%d ycol=%d",
			xcol, ycol)
	}
	xcol--
	ycol--
	// Now the columns numbers are hypothetic.

	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	// We read with the TSV reader.
	s, err := readFromTSV(fi, xcol, ycol)
	if err != nil {
		// Try read the dataset in another way.
		// FIXME Why we use TSV if this can handle all cases?
		var rawdata []byte
		rawdata, err = ioutil.ReadFile(fname)
		if err != nil {
			// fmt.Println("Cannot read file", fname, err.Error())
			return nil, err
		}
		s, err = parseDataSet(rawdata, xcol, ycol)
	}
	return s, err
}

// newTSVReader constructs a reader for TSV files
func newTSVReader(r io.Reader) *csv.Reader {
	csvr := csv.NewReader(r)
	csvr.Comma = '\t'
	csvr.Comment = commentPrefix
	csvr.FieldsPerRecord = 0 // The expected number of columns is derived from the first line
	csvr.LazyQuotes = false
	csvr.TrimLeadingSpace = true
	return csvr
}

// readFromTSV reads TSV file and takes columns for X and Y numbered from 1
func readFromTSV(r io.Reader, xcol, ycol int) (*XY, error) {
	tsvreader := newTSVReader(r)
	records, err := tsvreader.ReadAll() // [][]string
	if err != nil {
		return nil, err
	}

	if len(records[0]) < ycol+1 {
		return nil, csv.ErrFieldCount
	}
	var (
		data [][2]float64
		i    int
	)
	meta := make(map[string]string)
	entry := [2]string{}

	for _, e := range records {
		entry[0] = e[xcol]
		entry[1] = e[ycol]
		x, xerr := parseFloat(entry[0])
		y, yerr := parseFloat(entry[1])
		if xerr != nil || yerr != nil {
			meta[entry[0]] = entry[1]
		} else {
			data = append(data, [2]float64{x, y})
			i++
		}
	}
	spec := &XY{}
	spec.data = data
	spec.meta = meta
	return spec, nil
}
