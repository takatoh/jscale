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
	accNsHPF := HPF(accNs)
	accEwHPF := HPF(accEw)
	dyNs := integrate(accNsHPF, dt)
	dyEw := integrate(accEwHPF, dt)

	var periods []float64
	for t := 16; t <= 78; t += 2 {
		periods = append(periods, float64(t)/10.0)
	}
	ts := len(periods)
	sva := make([]float64, ts)
	for i := 0; i < ts; i++ {
		t := periods[i]
		w := 2.0 * math.Pi / t
		dxNs := RespSv(dumping, w, dt, n, accNsHPF)
		dxEw := RespSv(dumping, w, dt, n, accEwHPF)
		dxa := make([]float64, n)
		for j := 0; j < n; j++ {
			vaNs := dxNs[j] + dyNs[j]
			vaEw := dxEw[j] + dyEw[j]
			dxa[j] = math.Sqrt(vaNs*vaNs + vaEw*vaEw)
		}
		vel := seismicwave.Make("vel", dt, dxa)
		sva[i] = vel.AbsMax()
	}

	return sva
}

func HPF(acc []float64) []float64 {
	a1 := -2.0
	a2 := 1.0
	b1 := -1.995438545842
	b2 := 0.995448925627
	g0 := 0.997721867867

	accFilterd := make([]float64, len(acc))
	accFilterd[0], accFilterd[1] = acc[0], acc[1]
	y := make([]float64, len(acc))
	y[0], y[1] = 0.0, 0.0
	for t := 2; t < len(acc); t++ {
		y[t] = acc[t] + a1*acc[t-1] + a2*acc[t-2] - b1*y[t-1] - b2*y[t-2]
		accFilterd[t] = g0 * y[t]
	}

	return accFilterd
}

func integrate(ddy []float64, dt float64) []float64 {
	dy, _ := integral.Iacc(ddy, dt)
	return dy
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
