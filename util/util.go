package util

import (
	"math"
	"fmt"
)

func FFT(x []complex128, nn int, ind int) []complex128 {
//	fmt.Println(nn)

	var theta float64
	var i, j, k int
	var t complex128
	var tr, ti float64
	var m, kmax, istep int

	j = 0
	for i = 0; i < nn; i++ {
		fmt.Printf("i=%d, j=%d\n", i, j)
		if i < j {
			fmt.Println(j)
			t = x[j]
			x[j] = x[i]
			x[i] = t
		}
		m = nn / 2
//		fmt.Printf("m=%d\n", m)
		for j > m {
			j = j - m
			m = m / 2
			if m < 2 { break }
		}
		j = j + m
//		fmt.Printf("j=%d\n", j)
	}
	kmax = 1
	for kmax < nn {
//		fmt.Printf("kmax=%d\n", kmax)
		istep = kmax * 2
		for k = 0; k < kmax; k++ {
			theta = math.Pi * float64(ind * (k - 1)) / float64(kmax)
			for i = k; i < nn; i += istep {
				j = i + kmax - 1
//				fmt.Printf("j=%d\n", j)
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
