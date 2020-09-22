package srcwriter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/rubensayshi/duconverter/src/dustructs"
	"github.com/rubensayshi/duconverter/src/srcutils"
)

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

	for k, slot := range i.scriptExport.Slots {
		slotPath := path.Join(outputDir, "slots", fmt.Sprintf("%d.%s", k, slot.Name))
		err = os.MkdirAll(slotPath, 0777)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	for _, handler := range i.scriptExport.Handlers {
		slot := i.scriptExport.Slots[handler.Filter.SlotKey]
		slotPath := path.Join(outputDir, "slots", fmt.Sprintf("%d.%s", handler.Filter.SlotKey, slot.Name))

		handlerPath := path.Join(slotPath, fmt.Sprintf("%d.%s.lua", handler.Key, handler.Filter.Signature))

		code := handler.Code

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

	return nil
}
