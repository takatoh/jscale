package wave

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"math"
)

type Wave struct {
	Name string
	Dt   float64
	Data []float64
}

func newWave() *Wave {
	p := new(Wave)
	return p
}

func LoadCSV(filename string) []*Wave {
	var waves []*Wave
	var reader *csv.Reader
	var columns []string
	var err error
	var ns, ew, ud *Wave
	var t1, t2, d1, d2, d3 float64
	var dataNs, dataEw, dataUd []float64

	ns = newWave()
	ew = newWave()
	ud = newWave()
	t1 = 0.0
	t2 = 0.0

	read_file, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
	reader = csv.NewReader(read_file)

	columns, err = reader.Read()
	ns.Name = columns[1]
	ew.Name = columns[2]
	ud.Name = columns[3]
	for {
		columns, err = reader.Read()
		if err == io.EOF {
			dt := round(t2 - t1, 2)
			ns.Dt = dt
			ns.Data = dataNs
			waves = append(waves, ns)
			ew.Dt = dt
			ew.Data = dataEw
			waves = append(waves, ew)
			ud.Dt = dt
			ud.Data = dataUd
			waves = append(waves, ud)
			return waves
		}
		t1 = t2
		t2, _ = strconv.ParseFloat(columns[0], 64)
		d1, _ = strconv.ParseFloat(columns[1], 64)
		d2, _ = strconv.ParseFloat(columns[2], 64)
		d3, _ = strconv.ParseFloat(columns[3], 64)
		dataNs = append(dataNs, d1)
		dataEw = append(dataEw, d2)
		dataUd = append(dataUd, d3)
	}
}

func round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func LoadKNET(basename string) []*Wave {
	var waves []*Wave
	var dirs = []string{ "EW", "NS", "UD" }

	for _, dir := range dirs {
		waves = append(waves, loadKnetWave(basename, dir))
	}

	return waves
}
