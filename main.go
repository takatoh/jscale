package main

import (
	"fmt"
	"os"
	"math"
	"sort"

	"github.com/takatoh/jscale/wave"
)

func main() {
	var x, y, z []complex128
	var c []float64

	csvfile := os.Args[1]

	waves := wave.LoadCSV(csvfile)
//	dt := waves[0].Dt
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

	for i := 0; i < nn; i++ {
		xr := real(x[i])
		yr := real(y[i])
		zr := real(z[i])
		c = append(c, math.Sqrt(xr * xr + yr * yr + zr * zr))
	}
	sort.Slice(c, func(i, j int) bool {
		return c[i] > c[j]
	})

	for i := 0; i < nn; i++ {
		fmt.Println(c[i])
	}
}