package main

import (
	"fmt"
	"flag"
	"github.com/ArchieT/mno/wyniki"
)

func main() {
	var data,miejsce string
	var czas,czasplus int

	flag.StringVar(&data,"d","","data")
	flag.StringVar(&miejsce,"m","","miejsce")

	flag.IntVar(&czas,"t",0,"czas podstawowy")
	flag.IntVar(&czasplus,"p",0,"czas nisko naliczany")

	filename := flag.Args()[0]

	
}
