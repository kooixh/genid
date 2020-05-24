package main

import (
	"fmt"
	ap "github.com/akamensky/argparse"
	core "github.com/kooixh/genid/pkg/app"
	c "github.com/kooixh/genid/pkg/constants"
	g "github.com/kooixh/genid/pkg/generators"
	"os"
)

func main() {
	fmt.Println("Genid Id Genertion System")
	parser := ap.NewParser("main", "Configuration provided for core genid")
	suffix := parser.String("s", "suffix", &ap.Options{Required: false, Help: "String to print", Default: ""})
	prefix := parser.String("p", "prefix", &ap.Options{Required: false, Help: "String to print", Default: ""})
	calibrate := parser.Flag("c", "calibrate", &ap.Options{Required: false, Help: "Initiate calibration"})
	initial := parser.Int("i", "initial", &ap.Options{Required: false, Help: "Initial starting number for id",
		Default: 100})
	total  := parser.Int("t", "total", &ap.Options{Required: false, Help: "Total number stored each refill",
		Default: 100})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	if *calibrate {
		core.Calibrate(c.InitialOffset, *initial, *total)
		os.Exit(0)
		return
	}

	gen := g.NumericIdGenerator{Prefix: *prefix, Suffix: *suffix}
	gen.Generate()

	finished := make(chan int)
	go core.Refill(finished)
	exit := <- finished
	fmt.Println("genid done generation")
	os.Exit(exit)
}
