package jsonimporter

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rubensayshi/dubby/src/dustructs"
)

func Import(inputFile string) (*dustructs.ScriptExport, error) {
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

func (i *Importer) ReadFrom(inputFile string) (*dustructs.ScriptExport, error) {
	f, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	export := &dustructs.ScriptExport{}
	err = json.Unmarshal(f, export)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return export, nil
}
