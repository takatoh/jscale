package lpgm

import (
	"math"

	"github.com/takatoh/jscale/integral"
	"github.com/takatoh/sdof/directintegration"
	"github.com/takatoh/seismicwave"
)

const (
	dumping = 0.05
)

func Calc(ns, ew *seismicwave.Wave) []float64 {
	n := ns.NData()
	dt := ns.DT()
	accNs := ns.Data
	accEw := ew.Data
	ddy := make([]float64, n)
	for i := 0; i < n; i++ {
		ddy[i] = math.Sqrt(accNs[i]*accNs[i] + accEw[i]*accEw[i])
	}
	dy, _ := integral.Iacc(ddy, dt)

	var periods []float64
	for t := 16; t <= 78; t += 2 {
		periods = append(periods, float64(t)/10.0)
	}
	ts := len(periods)
	sva := make([]float64, ts)
	for i := 0; i < ts; i++ {
		t := periods[i]
		w := 2.0 * math.Pi / t
		dx := RespSv(dumping, w, dt, n, ddy)
		dxa := make([]float64, n)
		for j := 0; j < n; j++ {
			dxa[j] = dx[j] + dy[j]
		}
		vel := seismicwave.Make("vel", dt, dxa)
		sva[i] = vel.AbsMax()
	}

	return sva
}

func RespSv(h, w, dt float64, n int, ddy []float64) []float64 {
	_, dx, _ := directintegration.Nigam(h, w, dt, n, ddy)
	return dx
}

func Scale(intensity float64) string {
	var scale string
	if intensity < 5.0 {
		scale = "長周期地震動階級0"
	} else if intensity < 15.0 {
		scale = "長周期地震動階級1"
	} else if intensity < 50.0 {
		scale = "長周期地震動階級2"
	} else if intensity < 100.0 {
		scale = "長周期地震動階級3"
	} else {
		scale = "長周期地震動階級4"
	}
	return scale
}
