package util

import (
	"math"
	"math/cmplx"
)

func FFT(x []complex128, nn int, inv bool) []complex128 {
	if inv {
		for i := 0; i < nn; i++ {
			x[i] = complex(imag(x[i]), real(x[i]))
		}
		x = fft(x, nn)
		s := 1.0 / float64(nn)
		for i := 0; i < nn; i++ {
			x[i] = complex(imag(x[i]) * s, real(x[i]) * s)
		}
	} else {
		x = fft(x, nn)
	}
	return x
}

func fft(x []complex128, nn int) []complex128 {
	if nn == 1 {
		return x
	}
	nh := nn / 2
	even := make([]complex128, nh)
	odd := make([]complex128, nh)
	for i := 0; i < nh; i++ {
		even[i] = x[i] + x[i + nh]
		odd[i] = (x[i] - x[i + nh]) *
			cmplx.Exp(complex(0, 2 * float64(i) * math.Pi / float64(nn)))
	}
	even = fft(even, nh)
	odd = fft(odd, nh)
	for i := 0; i < nh; i++ {
		x[2 * i] = even[i]
		x[2 * i + 1] = odd[i]
	}
	return x
}

func DFT(x []complex128, nn int, inv int) []complex128 {
	ret := make([]complex128, 0)
	a := -2.0 * math.Pi / float64(nn) * float64(inv)

	for i := 0; i < nn; i++ {
		ret = append(ret, complex(0.0, 0.0))
		for j := 0; j < nn; j++ {
			fi := float64(i)
			fj := float64(j)
			ret[i] = ret[i] + (x[j] * complex(math.Cos(a * fi * fj), math.Sin(a * fi * fj)))
		}
		if inv > 0 {
			ret[i] = ret[i] * complex(1.0 / float64(nn), 0.0)
		}
	}

	return ret
}

func Filter(x []complex128, dt float64, nn int) []complex128 {
	var nfold, i int
	var f, y float64
	var f1, f2, f3 float64

	nfold = nn / 2 + 1

	x[0] = complex(0.0, 0.0)
	for i = 1; i < nfold; i++ {
		f = float64(i) / float64(nn) / dt
		y = f / 10.0
		f1 = math.Sqrt(1.0 / f)
		f2 = 1.0 / math.Sqrt(1.0 + 0.694 * math.Pow(y, 2.0) + 0.241 * math.Pow(y, 4.0) + 0.0557 * math.Pow(y, 6.0) + 0.009664 * math.Pow(y, 8.0) + 0.00134 * math.Pow(y, 10.0) + 0.000155 * math.Pow(y, 12.0))
		f3 = math.Sqrt(1.0 - math.Exp(-1.0 * math.Pow((2.0 * f), 3.0)))
		x[i] = complex(f1 * f2 * f3 * real(x[i]), f1 * f2 * f3 * imag(x[i]))
		x[nn - i] = complex(real(x[i]), -1.0 * imag(x[i]))
	}

	f = float64(nfold - 1) / float64(nn) / dt
	y = f / 10.0
	f1 = math.Sqrt(1.0 / f)
	f2 = 1.0 / math.Sqrt(1.0 + 0.694 * math.Pow(y, 2.0) + 0.241 * math.Pow(y, 4.0) + 0.0557 * math.Pow(y, 6.0) + 0.009664 * math.Pow(y, 8.0) + 0.00134 * math.Pow(y, 10.0) + 0.000155 * math.Pow(y, 12.0))
	f3 = math.Sqrt(1.0 - math.Exp(-1.0 * math.Pow(2.0 * f, 3.0)))
	x[nfold] = complex(f1 * f2 * f3 * real(x[nfold]), f1 * f2 * f3 * imag(x[nfold]))

	return x
}
