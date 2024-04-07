package intensity

import (
	"math"
	"math/cmplx"
	"sort"

	"github.com/takatoh/fft"
	"github.com/takatoh/seismicwave"
)

func Calc(ns, ew, ud *seismicwave.Wave) float64 {
	var x, y, z []complex128
	var v []float64
	var dt float64
	var n, nn int
	var a, I float64

	dt = ns.DT()
	n = ns.NData()
	x, nn = fft.MakeComplexData(ns.Data)
	y, _ = fft.MakeComplexData(ew.Data)
	z, _ = fft.MakeComplexData(ud.Data)

	// FFT で周波数領域へ
	x = fft.FFT(x, nn)
	y = fft.FFT(y, nn)
	z = fft.FFT(z, nn)

	// フィルタをかける
	x = Filter(x, dt, nn)
	y = Filter(y, dt, nn)
	z = Filter(z, dt, nn)

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

func Filter(x []complex128, dt float64, nn int) []complex128 {
	var nfold, i int
	var f, y float64
	var f1, f2, f3 float64

	nfold = nn / 2

	x[0] = complex(0.0, 0.0)
	for i = 1; i < nfold; i++ {
		f = float64(i) / float64(nn) / dt
		y = f / 10.0
		f1 = filter1(f)
		f2 = filter2(y)
		f3 = filter3(f)
		x[i] = complex(f1*f2*f3, 0.0) * x[i]
		x[nn-i] = cmplx.Conj(x[i])
	}

	f = float64(nfold) / float64(nn) / dt
	y = f / 10.0
	f1 = filter1(f)
	f2 = filter2(y)
	f3 = filter3(f)
	x[nfold] = complex(f1*f2*f3, 0.0) * x[nfold]

	return x
}

func filter1(f float64) float64 {
	return math.Sqrt(1.0 / f)
}

func filter2(y float64) float64 {
	return 1.0 / math.Sqrt(1.0+
		0.694*math.Pow(y, 2.0)+
		0.241*math.Pow(y, 4.0)+
		0.0557*math.Pow(y, 6.0)+
		0.009664*math.Pow(y, 8.0)+
		0.00134*math.Pow(y, 10.0)+
		0.000155*math.Pow(y, 12.0))
}

func filter3(f float64) float64 {
	return math.Sqrt(1.0 - math.Exp(-1.0*math.Pow(f/0.5, 3.0)))
}

func Scale(intensity float64) string {
	var scale string
	if intensity < 0.5 {
		scale = "震度0"
	} else if intensity < 1.5 {
		scale = "震度1"
	} else if intensity < 2.5 {
		scale = "震度2"
	} else if intensity < 3.5 {
		scale = "震度3"
	} else if intensity < 4.5 {
		scale = "震度4"
	} else if intensity < 5.0 {
		scale = "震度5弱"
	} else if intensity < 5.5 {
		scale = "震度5強"
	} else if intensity < 6.0 {
		scale = "震度6弱"
	} else if intensity < 6.5 {
		scale = "震度6強"
	} else {
		scale = "震度7"
	}
	return scale
}
