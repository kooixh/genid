package main

import (
	"fmt"
	ap "github.com/akamensky/argparse"
	core "github.com/kooixh/genid/pkg/app"
	c "github.com/kooixh/genid/pkg/constants"
	"os"
)

func main() {
	fmt.Println("Genid Id Generation System")

	parser := ap.NewParser("main", "Configuration provided for core genid")
	calibrate := parser.Flag("c", "calibrate", &ap.Options{Required: false, Help: "Initiate calibration"})
	refill := parser.Flag("r", "refill", &ap.Options{Required: false, Help: "Refill ids"})
	initial := parser.Int("", "initial", &ap.Options{Required: false, Help: "Initial starting number for id",
		Default: 100})
	total  := parser.Int("", "total", &ap.Options{Required: false, Help: "Total number stored each refill",
		Default: 100})
	idType := parser.String("t", "type", &ap.Options{Required: false, Help: "Type of id to generate (alphanum, num). Default alphanum",
		Default: c.GeneratorTypeAlphaNum})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	refillChannel := make(chan int)
	if *calibrate {
		settings := core.CalibrationSettings{
			Initial: *initial,
			Offset: c.InitialOffset,
			Total: *total,
			IdType: *idType,
		}
		validate, msg := settings.Validate()
		if !validate {
			fmt.Println(msg)
			os.Exit(1)
		}
		core.Calibrate(settings, refillChannel)
		exit := <-refillChannel
		os.Exit(exit)
	} else if *refill {
		go core.Refill(refillChannel)
		exit := <-refillChannel
		os.Exit(exit)
	}

	idChannel := make(chan string)
	go core.GenerateNewId(idChannel, refillChannel)
	fmt.Println("id generated is " + <-idChannel)
	exit := <-refillChannel
	os.Exit(exit)
}
