package lpgm

import (
	"math"

	"github.com/takatoh/sdof/directintegration"
	"github.com/takatoh/seismicwave"
)

func Calc(ns, ew *seismicwave.Wave) float64 {
	n := ns.NData()
	accNs := ns.Data
	accEw := ew.Data
	ddy := make([]float64, n)
	for i := 0; i < n; i++ {
		ddy[i] = math.Sqrt(accNs[i]*accNs[i] + accEw[i]*accEw[i])
	}

	t := 1.6
	w := 2.0 * math.Pi / t
	_, dx, _ := directintegration.Nigam(0.05, w, ns.DT(), n, ddy)
	vel := seismicwave.Make("vel", ns.DT(), dx)
	sv := vel.AbsMax()

	return sv
}
