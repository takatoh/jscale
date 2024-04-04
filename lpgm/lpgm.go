package lpgm

import (
	"math"

	"github.com/takatoh/sdof/directintegration"
	"github.com/takatoh/seismicwave"
)

const (
	dumping = 0.05
)

func Calc(ns, ew *seismicwave.Wave) []float64 {
	n := ns.NData()
	accNs := ns.Data
	accEw := ew.Data
	ddy := make([]float64, n)
	for i := 0; i < n; i++ {
		ddy[i] = math.Sqrt(accNs[i]*accNs[i] + accEw[i]*accEw[i])
	}

	var periods []float64
	for t := 16; t <= 78; t += 2 {
		periods = append(periods, float64(t)/10.0)
	}
	ts := len(periods)
	sv := make([]float64, ts)
	for i := 0; i < ts; i++ {
		t := 1.6
		w := 2.0 * math.Pi / t
		_, dx, _ := directintegration.Nigam(dumping, w, ns.DT(), n, ddy)
		vel := seismicwave.Make("vel", ns.DT(), dx)
		sv[i] = vel.AbsMax()
	}

	return sv
}
