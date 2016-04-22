// Package xy is a simple library for manipulation of X,Y data
package xy

import (
	"math"
	"testing"
)

func TestRoundFloat64(t *testing.T) {
	nums := []struct {
		orig float64
		prec int
		want float64
	}{
		{orig: math.E * 1000, prec: 4, want: 2718.2818},
		{orig: math.E * 1000, prec: 3, want: 2718.2820},
		{orig: math.E * 1000, prec: 2, want: 2718.2800},
		{orig: math.E * 1000, prec: 1, want: 2718.3000},
		{orig: math.E * 1000, prec: 0, want: 2718.0000},
		{orig: math.E * 1000, prec: -1, want: 2720.0000},
		{orig: math.E * 1000, prec: -2, want: 2700.0000},
		{orig: math.E * 1000, prec: -3, want: 3000.0000},
		{orig: math.E * 1000, prec: -4, want: 0000.0000},
	}
	for _, p := range nums {
		got := roundFloat64(p.orig, p.prec)
		if got != p.want {
			t.Errorf("roundFloat64(%q, %q) == %q, want %q",
				p.orig, p.prec, got, p.want)
		}

	}
}
