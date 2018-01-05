package main

import (
	"fmt"
	"os"
	"math"
	"sort"
	"flag"

	"github.com/takatoh/jscale/wave"
	"github.com/takatoh/jscale/util"
	"github.com/takatoh/jscale/fft"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s [option] <wave.csv>
Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_check := flag.Bool("check", false, "Check FFt result.")
	flag.Parse()

	var x, y, z []complex128
	var c []float64

	csvfile := flag.Args()[0]

	waves := wave.LoadCSV(csvfile)
	dt := waves[0].Dt
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

	// FFT で周波数領域へ
	x = fft.FFT(x, nn, false)
	y = fft.FFT(y, nn, false)
	z = fft.FFT(z, nn, false)

	if *opt_check {
		x = util.FFT(x, nn, true)
		for i := 0; i < n; i++ {
			fmt.Printf("%d,%f,%f\n", i, real(x[i]), imag(x[i]))
		}
		os.Exit(0)
	}

	// フィルタをかける
	x = util.Filter(x, dt, nn)
	y = util.Filter(y, dt, nn)
	z = util.Filter(z, dt, nn)

	// FFT で時間領域に戻す
	x = fft.FFT(x, nn, true)
	y = fft.FFT(y, nn, true)
	z = fft.FFT(z, nn, true)

	for i := 0; i < nn; i++ {
		xr := real(x[i])
		yr := real(y[i])
		zr := real(z[i])
		c = append(c, math.Sqrt(xr * xr + yr * yr + zr * zr))
	}
	sort.Slice(c, func(i, j int) bool {	return c[i] > c[j] })
	a := c[int(0.3 / dt)]
	I := 2 * math.Log10(a) + 0.94
	I = math.Floor(math.Floor(I * 100.0 + 0.5) / 10.0) / 10.0

	fmt.Printf("計測震度 %.1f\n", I)
	if I < 0.5 {
		fmt.Println("震度0")
	} else if I < 1.5 {
		fmt.Println("震度1")
	} else if I < 2.5 {
		fmt.Println("震度2")
	} else if I < 3.5 {
		fmt.Println("震度3")
	} else if I < 4.5 {
		fmt.Println("震度4")
	} else if I < 5.0 {
		fmt.Println("震度5弱")
	} else if I < 5.5 {
		fmt.Println("震度5強")
	} else if I < 6.0 {
		fmt.Println("震度6弱")
	} else if I < 6.5 {
		fmt.Println("震度6強")
	} else {
		fmt.Println("震度7")
	}
}