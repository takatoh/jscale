package main

import (
	"fmt"
	"os"
	"math"
	"sort"
	"flag"

	"github.com/takatoh/fft"
	"github.com/takatoh/jscale/wave"
	"github.com/takatoh/jscale/filter"
)

const (
	progVersion = "v1.1.0"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
`Usage:
  %s <wavefile.csv>
  %s -jma <wavefile.txt>
  %s -knet <wavefile>

Options:
`, os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	opt_jma := flag.Bool("jma", false, "Load JMA waves.")
	opt_knet := flag.Bool("knet", false, "Load KNET waves.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	filename := flag.Args()[0]

	var waves []*wave.Wave
	var err error
	if *opt_jma {
		waves, err = wave.LoadJMA(filename)
	} else if *opt_knet {
		waves, err = wave.LoadKNET(filename)
	} else {
		waves, err = wave.LoadCSV(filename)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	I := calcIntensity(waves[0], waves[1], waves[2])

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

func calcIntensity(ns, ew, ud *wave.Wave) float64 {
	var x, y, z []complex128
	var v []float64
	var dt float64
	var n, nn int
	var a, I float64

	dt = ns.Dt
	n = len(ns.Data)
	nn = int(math.Pow(2.0, math.Ceil(math.Log10(float64(n)) / math.Log10(2.0))))

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
		v = append(v, math.Sqrt(xr * xr + yr * yr + zr * zr))
	}
	sort.Slice(v, func(i, j int) bool { return v[i] > v[j] })
//	fmt.Println(int(0.3 / dt) - 1)
	a = v[int(0.3 / dt) - 1]
	I = 2.0 * math.Log10(a) + 0.94
	I = math.Floor(math.Floor(I * 100.0 + 0.5) / 10.0) / 10.0

	return I
}
