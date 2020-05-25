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
		os.Exit(initCalibration(*initial, *total, *idType, refillChannel))
	} else if *refill {
		os.Exit(initRefill(refillChannel))
	} else if *info {
		os.Exit(initPrintInfo())
	}

	if !core.CheckIsCalibrated() {
		fmt.Println("Genid is not calibrated, calibrate before using")
		os.Exit(1)
	}
	idChannel := make(chan string)
	go core.GenerateNewId(idChannel, refillChannel)
	fmt.Println("id generated is " + <-idChannel)
	exit := <-refillChannel
	os.Exit(exit)
}


func initCalibration(initial int, total int, idType string, refillChannel chan int) int {
	settings := core.CalibrationSettings{
		Initial: initial,
		Offset: c.InitialOffset,
		Total: total,
		IdType: idType,
	}
	fmt.Println("Starting Genid calibration with settings", settings)
	validate, msg := settings.Validate()
	if !validate {
		fmt.Println(msg)
		return 1
	}
	core.Calibrate(settings, refillChannel)
	exit := <-refillChannel
	fmt.Println("Completed Genid Calibration")
	return exit
}

func initRefill(refillChannel chan int) int {
	fmt.Println("Refilling the lists of Ids")
	if !core.CheckIsCalibrated() {
		fmt.Println("Genid is not calibrated, calibrate before using")
		return 1
	}
	go core.Refill(refillChannel)
	exit := <-refillChannel
	fmt.Println("Successfully refilled the lists of Ids")
	return exit
}

func initPrintInfo() int {
	appInfo, err := core.ReturnAppInfo()
	if err != nil {
		fmt.Println(err)
		return 1
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
	return 0
}