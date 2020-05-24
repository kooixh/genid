package main

import (
	"fmt"
	ap "github.com/akamensky/argparse"
	core "github.com/kooixh/genid/pkg/app"
	c "github.com/kooixh/genid/pkg/constants"
	"os"
)

func main() {
	fmt.Println("Genid Id Genertion System")
	parser := ap.NewParser("main", "Configuration provided for core genid")
	calibrate := parser.Flag("c", "calibrate", &ap.Options{Required: false, Help: "Initiate calibration"})
	refill := parser.Flag("rf", "refill", &ap.Options{Required: false, Help: "Refill ids"})
	initial := parser.Int("it", "initial", &ap.Options{Required: false, Help: "Initial starting number for id",
		Default: 100})
	total  := parser.Int("t", "total", &ap.Options{Required: false, Help: "Total number stored each refill",
		Default: 100})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	refillChannel := make(chan int)
	if *calibrate {
		core.Calibrate(c.InitialOffset, *initial, *total, refillChannel)
		exit := <-refillChannel
		os.Exit(exit)
		return
	} else if *refill {
		go core.Refill(refillChannel)
		exit := <-refillChannel
		os.Exit(exit)
		return
	}
	idChannel := make(chan string)
	go core.GenerateNewId(idChannel, refillChannel)
	fmt.Println("id generated is " + <-idChannel)
	exit := <-refillChannel
	os.Exit(exit)
}
