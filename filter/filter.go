package filter

import (
	"math"
)

func Filter(x []complex128, dt float64, nn int) []complex128 {
	var nfold, i int
	var f, y float64
	var f1, f2, f3 float64

	nfold = nn / 2 + 1

	x[0] = complex(0.0, 0.0)
	for i = 1; i < nfold; i++ {
		f = float64(i) / float64(nn) / dt
		y = f / 10.0
		f1 = filter1(f)
		f2 = filter2(y)
		f3 = filter3(f)
		x[i] = complex(f1 * f2 * f3 * real(x[i]), f1 * f2 * f3 * imag(x[i]))
		x[nn - i] = complex(real(x[i]), -1.0 * imag(x[i]))
	}

	f = float64(nfold - 1) / float64(nn) / dt
	y = f / 10.0
	f1 = filter1(f)
	f2 = filter2(y)
	f3 = filter3(f)
	x[nfold] = complex(f1 * f2 * f3 * real(x[nfold]), f1 * f2 * f3 * imag(x[nfold]))

	return x
}

func filter1(f float64) float64 {
	return math.Sqrt(1.0 / f)
}

func filter2(y float64) float64 {
	return 1.0 / math.Sqrt(1.0 +
		0.694 * math.Pow(y, 2.0) +
		0.241 * math.Pow(y, 4.0) +
		0.0557 * math.Pow(y, 6.0) +
		0.009664 * math.Pow(y, 8.0) +
		0.00134 * math.Pow(y, 10.0) +
		0.000155 * math.Pow(y, 12.0))
}

func filter3(f float64) float64 {
	return math.Sqrt(1.0 - math.Exp(-1.0 * math.Pow(f / 0.5, 3.0)))
}
