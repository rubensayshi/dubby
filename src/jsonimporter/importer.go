package jsonimporter

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rubensayshi/dubby/src/srcutils"

	"github.com/pkg/errors"
)

func Import(inputFile string) (*srcutils.ScriptExport, error) {
	i := NewImporter()

	scriptExport, err := i.ReadFrom(inputFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return scriptExport, nil
}

type Importer struct {
}

func NewImporter() *Importer {
	return &Importer{}
}

func (i *Importer) ReadFrom(inputFile string) (*srcutils.ScriptExport, error) {
	f, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	export := &srcutils.ScriptExport{}
	err = json.Unmarshal(f, export)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return export, nil
}
