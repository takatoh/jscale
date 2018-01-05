package fft

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
