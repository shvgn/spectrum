package xy

import (
	"math"
	"strconv"
	"strings"
)

// roundFloat64 rounds the float64 with the desired precision
func roundFloat64(f float64, prec int) float64 {
	shift := math.Pow(10, float64(prec))
	return math.Floor(f*shift+0.5) / shift
}

// parseFloat parses float64 from a string where the delimiter
// is either a dot or a comma
func parseFloat(s string) (float64, error) {
	ns := strings.TrimSpace(s)
	ns = strings.Replace(ns, ",", ".", 1)
	return strconv.ParseFloat(ns, 64)
}
