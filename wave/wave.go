package wave

import (
	"encoding/csv"
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"regexp"
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

func LoadCSV(filename string) ([]*Wave, error) {
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

	read_file, err := os.Open(filename)
	if err != nil {
		return waves, err
	}
	defer read_file.Close()

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
			return waves, nil
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

func LoadKNET(basename string) ([]*Wave, error) {
	var waves []*Wave
	var dirs = []string{ "NS", "EW", "UD" }

	for _, dir := range dirs {
		wave, err := loadKnetWave(basename, dir)
		if err != nil {
			return waves, err
		}
		waves = append(waves, wave)
	}

	return waves, nil
}

func loadKnetWave(basename, dir string) (*Wave, error) {
	var dt float64
	var scaleFactor float64
	wave := newWave()
	data := make([]float64, 0)

	f, err := os.Open(basename + "." + dir)
	if err != nil {
		return wave, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "Memo") == 0 {
			break
		}
		if strings.Index(line, "Sampling Freq(Hz)") == 0 {
			s := regexp.MustCompile(" +").Split(line, 3)
			s2 := strings.Trim(s[2], "Hz")
			f, _ := strconv.ParseFloat(s2, 64)
			dt = 1.0 / f
		}
		if strings.Index(line, "Scale Factor") == 0 {
			s := regexp.MustCompile(" +").Split(line, 3)
			s2 := regexp.MustCompile(`\(gal\)/`).Split(s[2], 2)
			f1, _ := strconv.ParseFloat(s2[0], 64)
			f2, _ := strconv.ParseFloat(s2[1], 64)
			scaleFactor = f1 / f2
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		dataStrings := regexp.MustCompile(" +").Split(line, 8)
		for _, s := range dataStrings {
			d, _ := strconv.ParseFloat(s, 64)
			data = append(data, d * scaleFactor)
		}
	}

	wave.Name = dir
	wave.Dt = dt
	wave.Data = data
	return wave, nil
}

func LoadJMA(filename string) ([]*Wave, error) {
	var waves []*Wave
	var ns, ew, ud *Wave
	var dt float64
	var dataNS, dataEW, dataUD []float64

	ns = newWave()
	ns.Name = "NS"
	ew = newWave()
	ew.Name = "EW"
	ud = newWave()
	ud.Name = "UD"

	f, err := os.Open(filename)
	if err != nil {
		return waves, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, " NS,EW,UD") == 0 {
			break
		}
		if strings.Index(line, " SAMPLING RATE") == 0 {
			srb := regexp.MustCompile(`\d+`).Find([]byte(line))
			sr, _ := strconv.ParseFloat(string(srb), 64)
			dt = 1.0 / sr
		}
	}
	ns.Dt = dt
	ew.Dt = dt
	ud.Dt = dt

	reader := csv.NewReader(f)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			ns.Data = dataNS
			waves = append(waves, ns)
			ew.Data = dataEW
			waves = append(waves, ew)
			ud.Data = dataEW
			waves = append(waves, ud)
			break
		}
		d0, _ := strconv.ParseFloat(row[0], 64)
		dataNS = append(dataNS, d0)
		d1, _ := strconv.ParseFloat(row[1], 64)
		dataEW = append(dataEW, d1)
		d2, _ := strconv.ParseFloat(row[2], 64)
		dataUD = append(dataUD, d2)
	}

	return waves, nil
}
