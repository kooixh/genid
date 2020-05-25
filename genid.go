package main

import (
	"fmt"
	ap "github.com/akamensky/argparse"
	core "github.com/kooixh/genid/pkg/app"
	c "github.com/kooixh/genid/pkg/constants"
	"github.com/kooixh/genid/utils/redis"
	"os"
)

func main() {
	fmt.Println("Genid Id Generation System")

	parser := ap.NewParser("main", "Configuration provided for core genid")
	calibrate := parser.Flag("c", "calibrate", &ap.Options{Required: false, Help: "Initiate calibration"})
	refill := parser.Flag("r", "refill", &ap.Options{Required: false, Help: "Refill ids"})
	info := parser.Flag("i", "info", &ap.Options{Required: false, Help: "Show info of app"})
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

	res := redis.Ping()
	if res.Err() != nil {
		fmt.Println("error connecting to Redis database")
		os.Exit(1)
	}

	refillChannel := make(chan int)
	if *calibrate {
		fmt.Println("Starting calibration")
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
		fmt.Println("Calibration completed")
		os.Exit(exit)
	} else if *refill {
		fmt.Println("Starting refill")
		go core.Refill(refillChannel)
		exit := <-refillChannel
		fmt.Println("Refill done")
		os.Exit(exit)
	} else if *info {
		appInfo, err := core.ReturnAppInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Status:", appInfo.Status)
		fmt.Println("Version:", appInfo.Version)
		if appInfo.Status == c.Calibrated {
			fmt.Println("Next Id:", appInfo.NextId)
			settings := appInfo.Settings
			fmt.Println("Initial Id:", settings.Initial)
			fmt.Println("Id Type:", settings.IdType)
			fmt.Println("Offset:", settings.Offset)
			fmt.Println("Refill Amount:", settings.Total)
		}
		os.Exit(0)
	}

	idChannel := make(chan string)
	go core.GenerateNewId(idChannel, refillChannel)
	fmt.Println("id generated is " + <-idChannel)
	exit := <-refillChannel
	os.Exit(exit)
	
}
