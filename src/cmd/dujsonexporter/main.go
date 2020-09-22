package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/rubensayshi/duconverter/src/srcreader"
)

func main() {
	flag.Parse()

	srcDir := flag.Arg(0)
	outputFile := flag.Arg(1)

	scriptExport, err := srcreader.Read(srcDir)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	res, err := json.Marshal(scriptExport)

	err = ioutil.WriteFile(outputFile, res, 0777)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
