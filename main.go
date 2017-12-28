package main

import (
	"fmt"
	"os"
	"math"
	"sort"

	"github.com/takatoh/jscale/wave"
	"github.com/takatoh/jscale/util"
)

func main() {
	var x, y, z []complex128
	var c []float64

	csvfile := os.Args[1]

	waves := wave.LoadCSV(csvfile)
	dt := waves[0].Dt
//	fmt.Println(dt)
	n := len(waves[0].Data)
	nn := int(math.Pow(2.0, math.Ceil(math.Log10(float64(n)) / math.Log10(2.0))))

	for i := 0; i < n; i++ {
		x = append(x, complex(waves[0].Data[i], 0.0))
		y = append(y, complex(waves[1].Data[i], 0.0))
		z = append(z, complex(waves[2].Data[i], 0.0))
	}
	for i := n; i < nn; i++ {
		x = append(x, complex(0.0, 0.0))
		y = append(y, complex(0.0, 0.0))
		z = append(z, complex(0.0, 0.0))
	}

//	x = util.FFT(x, nn, -1)
//	y = util.FFT(y, nn, -1)
//	z = util.FFT(z, nn, -1)
	x = util.DFT(x, nn, -1)
	y = util.DFT(y, nn, -1)
	z = util.DFT(z, nn, -1)

//	filter waves.
	x = util.Filter(x, dt, nn)
	y = util.Filter(y, dt, nn)
	z = util.Filter(z, dt, nn)

//	x = util.FFT(x, nn, 1)
//	y = util.FFT(y, nn, 1)
//	z = util.FFT(z, nn, 1)
	x = util.DFT(x, nn, 1)
	y = util.DFT(y, nn, 1)
	z = util.DFT(z, nn, 1)

//	for i := 0; i < n; i++ {
//		fmt.Printf("x=%f, %f\n", real(x[i]), imag(x[i]))
//	}

	for i := 0; i < nn; i++ {
		xr := real(x[i])
		yr := real(y[i])
		zr := real(z[i])
		c = append(c, math.Sqrt(xr * xr + yr * yr + zr * zr))
	}
	sort.Slice(c, func(i, j int) bool {
		return c[i] > c[j]
	})
//	fmt.Println(int(0.3 / dt))
	a := c[int(0.3 / dt)]
	I := 2 * math.Log10(a) + 0.94
	I = math.Floor(math.Floor(I * 100.0 + 0.5) / 10.0) / 10.0

	fmt.Printf("計測震度= %.1f\n", I)
}