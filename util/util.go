package util

import (
	"math"
)

func FFT(x []complex128, nn int, ind, int) []complex128 {
	var theta float64
	var i, j, k int
	var t complex128
	var tr, ti float64
	var m, kmax, istep int

	j = 1
	for i = 0; i < nn; i++ {
		if i < j {
			t = x[j]
			x[j] = x[i]
			x[i] = t
		}
		m = nn / 2
		for j > m {
			j = j - m
			m = m / 2
			if m < 2 { break }
		}
		j = j + m
	}
	kmax = 1
	for kmax < nn {
		istep = kmax * 2
		for k = 0; k < kmax; k++ {
			theta = math.Pi * float64(ind * (k - 1)) / float64(kmax)
			for i = k - 1; k < nn; i += istep {
				j = i + kmax
				tr = real(x[j]) * math.Cos(theta) - imag(x[j]) * math.Sin(theta)
				ti = real(x[j]) * math.Sin(theta) + imag(x[j]) * math.Cos(theta)
				x[j] = complex(real(x[i]) - tr, imag(x[i]) - ti)
				x[i] = complex(real(x[i]) + tr, imag(x[i]) + ti)
			}
		}
		kmax = istep
	}

	return x
}
