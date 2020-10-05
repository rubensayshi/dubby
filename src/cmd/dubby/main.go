package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/rubensayshi/dubby/src/jsonimporter"
	"github.com/rubensayshi/dubby/src/srcreader"
	"github.com/rubensayshi/dubby/src/srcwriter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "dubby"
	app.Usage = ""

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 1)
		return nil
	}

	app.Commands = []*cli.Command{{
		Name:      "parse-to-src",
		Aliases:   []string{},
		Usage:     "parse a json file into a source directory",
		ArgsUsage: "inputfile srcdir",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			inputfile := c.Args().Get(0)
			if inputfile == "" {
				cli.ShowCommandHelpAndExit(c, "parse-to-src", 1)
				return nil
			}
			_, err := os.Stat(inputfile)
			if os.IsNotExist(err) {
				return errors.Errorf("can't open json file: %s", inputfile)
			}

			srcdir := c.Args().Get(0)
			if srcdir == "" {
				cli.ShowCommandHelpAndExit(c, "parse-to-src", 1)
				return nil
			}

			return parseToSrc(inputfile, srcdir)
		},
	}, {
		Name:      "export-to-json",
		Aliases:   []string{},
		Usage:     "compile a source directory and export to json",
		ArgsUsage: "srcdir outputfile",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "minify",
			},
		},
		Action: func(c *cli.Context) error {
			srcdir := c.Args().Get(0)
			if srcdir == "" {
				cli.ShowCommandHelpAndExit(c, "parse-to-src", 1)
				return nil
			}

			outputfile := c.Args().Get(1)
			if outputfile == "" {
				cli.ShowCommandHelpAndExit(c, "parse-to-src", 1)
				return nil
			}

			return exportToJson(srcdir, outputfile, c.Bool("minify"))
		},
	}}

	err := app.Run(os.Args)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func parseToSrc(inputfile string, srcdir string) error {
	scriptExport, err := jsonimporter.Import(inputfile)
	if err != nil {
		return errors.WithStack(err)
	}

	w := srcwriter.NewSrcWriter(scriptExport)
	err = w.WriteTo(srcdir)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func exportToJson(srcdir string, outputfile string, minify bool) error {
	reader := srcreader.NewSrcReader(srcdir, minify)

	err := reader.Read()
	if err != nil {
		return errors.WithStack(err)
	}

	scriptExport := reader.ScriptExport()

	res, err := json.Marshal(scriptExport)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(outputfile, res, 0666)
	if err != nil {
		return errors.WithStack(err)
	}

	report := reader.Report()
	if minify {
		p := float64(report.SrcLen-report.MinifiedLen) / float64(report.SrcLen) * 100
		fmt.Printf("minified %d bytes of lua -> %d (%.1f%% saved) \n", report.SrcLen, report.MinifiedLen, p)
	}

	return nil
}
