package srcwriter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rubensayshi/duconverter/src/dustructs"
	"github.com/rubensayshi/duconverter/src/srcutils"
)

var libHeaderRegex = regexp.MustCompile(`-- !DU\[lib]: (.*?)\n\n?`)

type SrcWriter struct {
	scriptExport dustructs.ScriptExport
}

func NewSrcWriter(scriptExport *dustructs.ScriptExport) *SrcWriter {
	return &SrcWriter{
		scriptExport: *scriptExport,
	}
}

func (i *SrcWriter) WriteTo(outputDir string) error {
	err := os.RemoveAll(outputDir)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(outputDir, 0777)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(path.Join(outputDir, "slots"), 0777)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(path.Join(outputDir, "lib"), 0777)
	if err != nil {
		return errors.WithStack(err)
	}

	for k, slot := range i.scriptExport.Slots {
		slotPath := path.Join(outputDir, "slots", fmt.Sprintf("%d.%s", k, slot.Name))
		err = os.MkdirAll(slotPath, 0777)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	libKey := 0
	handlerKeyOffset := 0

	for _, handler := range i.scriptExport.Handlers {
		code := handler.Code

		if strings.HasPrefix(code, "-- !DU[lib]: ") {
			libHeaders := libHeaderRegex.FindAllString(code, -1)
			libs := libHeaderRegex.Split(code, -1)[1:]

			if len(libs) != len(libHeaders) {
				return errors.Errorf("Lib header mismatch! libs=%d != headers=%d", len(libs), len(libHeaders))
			}

			for k, libCode := range libs {
				libHeader := libHeaders[k]

				libHeaderMatch := libHeaderRegex.FindStringSubmatch(libHeader)
				libName := libHeaderMatch[1]

				libPath := path.Join(outputDir, "lib", fmt.Sprintf("%d.%s.lua", libKey, libName))
				libKey += 1

				err = ioutil.WriteFile(libPath, []byte(libCode), 0666)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			handlerKeyOffset -= 1

		} else {
			slot := i.scriptExport.Slots[handler.Filter.SlotKey]
			slotPath := path.Join(outputDir, "slots", fmt.Sprintf("%d.%s", handler.Filter.SlotKey, slot.Name))

			handlerKey := handler.Key + handlerKeyOffset

			handlerPath := path.Join(slotPath, fmt.Sprintf("%d.%s.lua", handlerKey, handler.Filter.Signature))

			sig, err := srcutils.SignatureWithArgs(handler.Filter.Signature, handler.Filter.Args)
			if err != nil {
				return errors.WithStack(err)
			}

			// add signature as first line
			code = fmt.Sprintf("-- !DU: %s\n", sig) + code

			err = ioutil.WriteFile(handlerPath, []byte(code), 0666)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}
