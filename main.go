package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/takatoh/jscale/intensity"
	"github.com/takatoh/jscale/lpgm"
	"github.com/takatoh/seismicwave"
)

const (
	progVersion = "v1.5.0"
)

func main() {
	progName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s [-long-period] <wavefile.csv>
  %s [-long-period] -jma <wavefile.txt>
  %s [-long-period] -knet <wavefile>
  %s [-long-period] -fixed <input.toml>

Options:
`, progName, progName, progName, progName)
		flag.PrintDefaults()
	}
	opt_jma := flag.Bool("jma", false, "Load JMA waves.")
	opt_knet := flag.Bool("knet", false, "Load KNET waves.")
	opt_fixed := flag.Bool("fixed", false, "Load fixed format waves.")
	opt_lpgm := flag.Bool("long-period", false, "Caclulate intensity scale on long-period ground motion.")
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	filename := flag.Arg(0)

	var waves []*seismicwave.Wave
	var err error
	if *opt_jma {
		waves, err = seismicwave.LoadJMA(filename)
	} else if *opt_knet {
		waves, err = seismicwave.LoadKNETSet(filename)
	} else if *opt_fixed {
		waves, err = seismicwave.LoadFixedFormatWithTOML(filename)
	} else {
		waves, err = seismicwave.LoadCSV(filename)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if *opt_lpgm {
		fmt.Println("長周期地震動階級を計算します。")
		sv := lpgm.Calc(waves[0], waves[1])
		maxSv := 0.0
		for i := 0; i < len(sv); i++ {
			maxSv = math.Max(maxSv, sv[i])
		}
		fmt.Printf("Max Sv = %.1f\n", maxSv)
	} else {
		I := intensity.Calc(waves[0], waves[1], waves[2])

		fmt.Printf("計測震度 %.1f\n", I)
		fmt.Println(intensity.Scale(I))
	}
}
