package main

import (
	"fmt"
	"os"

	"github.com/takatoh/jscale/wave"
)

func main() {
	csvfile := os.Args[1]

	waves := wave.LoadCSV(csvfile)
	x := waves[0]
	y := waves[1]
	z := waves[2]

	n := len(x.Data)
	fmt.Printf("%s,%s,%s\n", x.Name, y.Name, z.Name)
	for i := 0; i < n; i++ {
		fmt.Printf("%f,%f,%f\n", x.Data[i], y.Data[i], z.Data[i])
	}
}