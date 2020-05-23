package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	g "github.com/kooixh/genid/pkg/generators"
	"os"
)

func main() {
	fmt.Println("Genid Id Genertion System")
	parser := argparse.NewParser("config", "Configuration provided for generation")
	suffix := parser.String("s", "suffix", &argparse.Options{Required: false, Help: "String to print", Default: ""})
	prefix := parser.String("p", "prefix", &argparse.Options{Required: false, Help: "String to print", Default: ""})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	gen := g.NumericIdGenerator{Prefix: *prefix, Suffix: *suffix}
	fmt.Println(gen.Generate())
}
