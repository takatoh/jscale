package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/takatoh/jscale/intensity"
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
  %s <wavefile.csv>
  %s -jma <wavefile.txt>
  %s -knet <wavefile>
  %s -fixed <input.toml>

Options:
`, progName, progName, progName, progName)
		flag.PrintDefaults()
	}
	opt_jma := flag.Bool("jma", false, "Load JMA waves.")
	opt_knet := flag.Bool("knet", false, "Load KNET waves.")
	opt_fixed := flag.Bool("fixed", false, "Load fixed format waves.")
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

	I := intensity.Calc(waves[0], waves[1], waves[2])

	fmt.Printf("計測震度 %.1f\n", I)
	fmt.Println(intensity.Scale(I))
}
