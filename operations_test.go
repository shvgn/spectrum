// Package xy is a simple library for manipulation of X,Y data

package xy

import "testing"

var sp1 XY = XY{
	[][2]float64{
		[2]float64{1.0, 5.34047454734599},
		[2]float64{2.0, 5.56134762834134},
		[2]float64{3.0, 5.88921617459386},
		[2]float64{4.0, 6.35335283236613},
		[2]float64{5.0, 6.97898699083615},
		[2]float64{6.0, 7.78037300453194},
		[2]float64{7.0, 8.753110988514},
		[2]float64{8.0, 9.86752255959972},
		[2]float64{9.0, 11.0653065971263},
		[2]float64{10.0, 12.2614903707369},
		[2]float64{11.0, 13.3527021141127},
		[2]float64{12.0, 14.2311634638664},
		[2]float64{13.0, 14.8019867330676},
		[2]float64{14.0, 15},
		[2]float64{15.0, 14.8019867330676},
		[2]float64{16.0, 14.2311634638664},
		[2]float64{17.0, 13.3527021141127},
		[2]float64{18.0, 12.2614903707369},
		[2]float64{19.0, 11.0653065971263},
		[2]float64{20.0, 9.86752255959972}}, nil}

// Shifted X of sp1, the same function
var sp1s XY = XY{
	[][2]float64{
		[2]float64{1.431, 5.42442917648117},
		[2]float64{2.431, 5.68780485174319},
		[2]float64{3.431, 6.07091119799867},
		[2]float64{4.431, 6.60202727323746},
		[2]float64{5.431, 7.30257917016815},
		[2]float64{6.431, 8.17970960358249},
		[2]float64{7.431, 9.21879639534797},
		[2]float64{8.431, 10.3779640489827},
		[2]float64{9.431, 11.5868146663982},
		[2]float64{10.431, 12.7510620523878},
		[2]float64{11.431, 13.7634518399934},
		[2]float64{12.431, 14.5195718399246},
		[2]float64{13.431, 14.9354569906083},
		[2]float64{14.431, 14.9629167289098},
		[2]float64{15.431, 14.598721120586},
		[2]float64{16.431, 13.885225783046},
		[2]float64{17.431, 12.9022685224626},
		[2]float64{18.431, 11.7524801902834},
		[2]float64{19.431, 10.5437428259392},
		[2]float64{20.431, 9.372915167086}}, nil}

var sp2 XY = XY{
	[][2]float64{
		[2]float64{1.0, 5.13945037358082},
		[2]float64{2.0, 5.35474279421912},
		[2]float64{3.0, 5.81654107234731},
		[2]float64{4.0, 6.70064254215682},
		[2]float64{5.0, 8.20493034257966},
		[2]float64{6.0, 10.465057082495},
		[2]float64{7.0, 13.4322089931297},
		[2]float64{8.0, 16.7722255225981},
		[2]float64{9.0, 19.8712110291761},
		[2]float64{10.0, 21.9982671886829},
		[2]float64{11.0, 22.5805929315872},
		[2]float64{12.0, 21.452539296646},
		[2]float64{13.0, 18.931661287365},
		[2]float64{14.0, 15.6743989266113},
		[2]float64{15.0, 12.4003882231},
		[2]float64{16.0, 9.64233128392267},
		[2]float64{17.0, 7.63504664116284},
		[2]float64{18.0, 6.35335283236613},
		[2]float64{19.0, 5.62893290550818},
		[2]float64{20.0, 5.26446496710796}}, nil}

// Check for data equality in two spectra
func equal(s1, s2 XY) bool {
	for i, p1 := range s1.data {
		p2 := s2.data[i]
		if p1 != p2 {
			return false
		}
	}
	return true
}

func EqualTest(t *testing.T) {
	if !equal(sp1, sp1) {
		t.Errorf("equal() doesn't work")
	}
}

func TestAdd(t *testing.T) {
	// names := []struct {
	// 	fname string
	// 	sfx   string
	// 	want  string
	// }{
	// 	{"data1.txt", "ev", "data1.ev.txt"},
	// 	{"data1.dat", "nm", "data1.nm.dat"},
	// 	{"data1.nm", "ev", "data1.ev"},
	// 	{"data1.ev", "nm", "data1.nm"},
	// 	{"data1.nm", "nm", "data1.nm"},
	// 	{"data1.ev", "ev", "data1.ev"},
	// 	{"d.nm.a.ev.t.a.ev.dat", "ev", "d.nm.a.ev.t.a.ev.dat"},
	// 	{"d.nm.a.ev.t.a.ev.dat", "nm", "d.nm.a.ev.t.a.nm.dat"},
	// 	{"d.nm.a.ev.t.a.ev", "ev", "d.nm.a.ev.t.a.ev"},
	// 	{"d.nm.a.ev.t.a.ev", "nm", "d.nm.a.ev.t.a.nm"},
	// }
	// for _, s := range names {
	// 	got := addPreSuffix(s.fname, s.sfx)
	// 	if got != s.want {
	// 		t.Errorf("addPreSuffix(%q, %q) == %q, want %q",
	// 			s.fname, s.sfx, got, s.want)
	// 	}

	// }
}
