package main

import (
	"flag"
	"fmt"

	"github.com/rubensayshi/duconverter/src/jsonimporter"
	"github.com/rubensayshi/duconverter/src/srcwriter"
)

func main() {
	flag.Parse()

	inputFile := flag.Arg(0)
	outputDir := flag.Arg(1)

	scriptExport, err := jsonimporter.Import(inputFile)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	w := srcwriter.NewSrcWriter(scriptExport)
	err = w.WriteTo(outputDir)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
