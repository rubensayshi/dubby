package importer

import (
	"encoding/json"
	"io/ioutil"
	"strings"

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

	if strings.HasSuffix(inputFile, ".json") {
		export := &srcutils.ScriptExport{}
		err = json.Unmarshal(f, export)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return export, nil

	} else if strings.HasSuffix(inputFile, ".conf") {
		export, err := srcutils.UnmarshalAutoConf(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return export, nil

	} else {
		return nil, errors.Errorf("Can only parse .json or .conf files")
	}
}
