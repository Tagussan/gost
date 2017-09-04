package main

import (
	"bufio"
	"fmt"
	goopt "github.com/droundy/goopt"
	"math"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
)

var fNoHead = goopt.Flag([]string{"-n", "--no-header"}, []string{}, "don't display header", "")
var fComplete = goopt.Flag([]string{"--complete"}, []string{}, "everything", "")
var fStrict = goopt.Flag([]string{"--strict"}, []string{}, "throws error for invalid input", "")

func main() {
	goopt.Description = func() string {
		return "simple statistics from command line"
	}
	goopt.Parse(nil)
	data := make([]float64, 0, 128)
	loadData(&data)
	dispData(data)
}

func dispData(data []float64) {
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 6, 2, ' ', 0)
	var format string
	if *fComplete {
		format = "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n"
	} else {
		format = "%v\t%v\t%v\t%v\t%v\t%v\t\n"
	}

	if !*fNoHead {
		if *fComplete {
			fmt.Fprintf(tw, format, "N", "min", "q1", "med", "q3", "max", "sum", "avg", "stddev", "stderr")
		} else {
			fmt.Fprintf(tw, format, "N", "min", "max", "avg", "stddev", "stderr")
		}
	}
	q1, med, q3 := quantile(data)
	if *fComplete {
		fmt.Fprintf(tw, format, len(data), min(data), q1, med, q3, max(data), sum(data), average(data), stddev(data), stderr(data))
	} else {
		fmt.Fprintf(tw, format, len(data), min(data), max(data), average(data), stddev(data), stderr(data))
	}
	tw.Flush()
}

func loadData(data *[]float64) {
	if len(goopt.Args) > 0 {
		for _, name := range goopt.Args {
			fp, err := os.Open(name)
			if err != nil {
				panic(err)
			}
			defer fp.Close()
			scanNumbers(fp, data)
		}
	} else {
		scanNumbers(os.Stdin, data)
	}
	if len(*data) == 0 {
		panic("No numbers given")
	}

}

func scanNumbers(fp *os.File, data *[]float64) {
	sc := bufio.NewScanner(fp)
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		text := sc.Text()
		r, err := strconv.ParseFloat(text, 64)
		if err != nil {
			if *fStrict {
				panic(fmt.Sprintf("invalid input %s", text))
			}
			fmt.Printf("warning: invalid input %s\n", text)
			continue
		}
		*data = append(*data, r)
	}
}

func sum(data []float64) float64 {
	x := float64(0)
	for _, r := range data {
		x += r
	}
	return x
}

func average(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func min(data []float64) float64 {
	x := data[0]
	for _, r := range data {
		if x > r {
			x = r
		}
	}
	return x
}

func max(data []float64) float64 {
	x := data[0]
	for _, r := range data {
		if x < r {
			x = r
		}
	}
	return x
}

func variance(data []float64) float64 {
	avg := average(data)
	x := float64(0)
	for _, r := range data {
		x += (r - avg) * (r - avg)
	}
	return x / float64(len(data))
}

func stddev(data []float64) float64 {
	return math.Sqrt(variance(data))
}

func stderr(data []float64) float64 {
	return stddev(data) / math.Sqrt(float64(len(data)))
}

func medSorted(data []float64) float64 {
	if len(data) == 1 {
		return data[0]
	}
	if len(data)%2 == 1 { //odd
		return data[(len(data)-1)/2]
	} else { //even
		return (data[len(data)/2-1] + data[len(data)/2]) / 2
	}
}

func quantile(data []float64) (float64, float64, float64) {
	sort.Float64s(data)
	halfL := data[:len(data)/2]
	halfH := data[len(data)/2:]
	return medSorted(halfL), medSorted(data), medSorted(halfH)
}
