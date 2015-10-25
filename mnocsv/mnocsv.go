package main

import (
	"flag"
	"os"
	"github.com/ArchieT/mno/wyniki"
)

var data, miejsce string
var czas, czasplus int
var filename string

func main() {
	flag.StringVar(&data, "d", "", "data")
	flag.StringVar(&miejsce, "m", "", "miejsce")

	flag.IntVar(&czas, "t", 0, "czas podstawowy")
	flag.IntVar(&czasplus, "p", 0, "czas nisko naliczany")

	flag.StringVar(&filename, "f", "", "nazwa pliku csv")

	flag.Parse()

	file,err := os.Open(filename)
	if err!=nil {
		panic(err)
	}

	w := wyniki.Daj(file)

	var z wyniki.Zawody
	z.Czas = czas
	z.CzasPlus = czasplus
	z.Data = data
	z.Miejsce = miejsce
	z.Wyniki = &w

	write := os.Stdout

	z.Present(write)
}
