package intensity

import (
	"math"
	"sort"

	"github.com/takatoh/fft"
	"github.com/takatoh/jscale/filter"
	"github.com/takatoh/seismicwave"
)

func Calc(ns, ew, ud *seismicwave.Wave) float64 {
	var x, y, z []complex128
	var v []float64
	var dt float64
	var n, nn int
	var a, I float64

	dt = ns.Dt
	n = len(ns.Data)
	nn = 2
	for nn < n {
		nn = nn * 2
	}

	for i := 0; i < n; i++ {
		x = append(x, complex(ns.Data[i], 0.0))
		y = append(y, complex(ew.Data[i], 0.0))
		z = append(z, complex(ud.Data[i], 0.0))
	}
	for i := n; i < nn; i++ {
		x = append(x, complex(0.0, 0.0))
		y = append(y, complex(0.0, 0.0))
		z = append(z, complex(0.0, 0.0))
	}

	// FFT で周波数領域へ
	x = fft.FFT(x, nn)
	y = fft.FFT(y, nn)
	z = fft.FFT(z, nn)

	// フィルタをかける
	x = filter.Filter(x, dt, nn)
	y = filter.Filter(y, dt, nn)
	z = filter.Filter(z, dt, nn)

	// FFT で時間領域に戻す
	x = fft.IFFT(x, nn)
	y = fft.IFFT(y, nn)
	z = fft.IFFT(z, nn)

	for i := 0; i < n; i++ {
		xr := real(x[i])
		yr := real(y[i])
		zr := real(z[i])
		v = append(v, math.Sqrt(xr*xr+yr*yr+zr*zr))
	}
	sort.Slice(v, func(i, j int) bool { return v[i] > v[j] })
	a = v[int(0.3/dt)-1]
	I = 2.0*math.Log10(a) + 0.94
	I = math.Floor(math.Floor(I*100.0+0.5)/10.0) / 10.0

	return I
}
