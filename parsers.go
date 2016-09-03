package xy

import (
	"encoding/csv"
	"strings"
)

const (
	commentPrefix = '#' // Lines starting with this rune will be ignored
)

// Text parser (FIXME must take reader)
func parseDataSet(b []byte, xcol, ycol int) (*XY, error) {
	var (
		lines = strings.Split(string(b), "\n")
		data  = make([][2]float64, 0, len(lines))
		meta  = map[string]string{}
	)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" ||
			strings.HasPrefix(line, string(commentPrefix)) {
			continue // Skip empty lines and comments
		}

		fields := strings.Fields(line)
		if len(fields) < ycol+1 {
			return nil, csv.ErrFieldCount
		}
		x, errx := parseFloat(fields[xcol])
		y, erry := parseFloat(fields[ycol])
		if errx != nil || erry != nil {
			// Not a point x,y hence must be a header
			name, value := parseHeader(line)
			meta[name] = value
			continue
		}
		// valid float64's
		data = append(data, [2]float64{x, y})
	}

	s := &XY{meta: meta, data: data}
	s.SortByX()
	return s, nil
}

// parseHeader parses a header string
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
