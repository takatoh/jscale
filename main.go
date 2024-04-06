package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

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
	opt_spectrum := flag.Bool("spectrum", false, "Output spectrum to file.")
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
		Sva := lpgm.Calc(waves[0], waves[1])
		maxSva := 0.0
		for i := 0; i < len(Sva); i++ {
			maxSva = math.Max(maxSva, Sva[i])
		}

		if *opt_spectrum {
			outputSpectra(lpgm.Periods(), Sva, "spec-lpgm.csv")
		}
		fmt.Printf("絶対速度応答スペクトルの最大値 %.1f cm/sec\n", maxSva)
		fmt.Println(lpgm.Scale(maxSva))
	} else {
		I := intensity.Calc(waves[0], waves[1], waves[2])

		fmt.Printf("計測震度 %.1f\n", I)
		fmt.Println(intensity.Scale(I))
	}
}

func outputSpectra(periods, responses []float64, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic("Error! Can not open the file")
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{"period", "Sva"}
	err = w.Write(header)
	if err != nil {
		panic("Error! Can not write to the file")
	}
	for i := 0; i < len(periods); i++ {
		period := strconv.FormatFloat(periods[i], 'f', 1, 64)
		response := strconv.FormatFloat(responses[i], 'f', 3, 64)
		record := []string{period, response}
		err := w.Write(record)
		if err != nil {
			panic("Error! Can not write to the file")
		}
	}
}
