package srcwriter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rubensayshi/dubby/src/srcutils"
)

var libHeaderRegex = regexp.MustCompile(`-- !DU\[lib]: (.*?)\n\n?`)

type SrcWriter struct {
	scriptExport srcutils.ScriptExport
}

func NewSrcWriter(scriptExport *srcutils.ScriptExport) *SrcWriter {
	return &SrcWriter{
		scriptExport: *scriptExport,
	}
}

type SlotSrc struct {
	key      int
	name     string
	mainCode []string
	handlers []*SlotSrcHandler
}

type SlotSrcHandler struct {
	code []string
	sig  string
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

	libKey := 0

	// create intermediate struct to hold data per slot, because we'll write the aggregate in 1 file
	slots := make(map[int]*SlotSrc, len(i.scriptExport.Slots))
	for i, slot := range i.scriptExport.Slots {
		slots[i] = &SlotSrc{
			key:      i,
			name:     slot.Name,
			mainCode: []string{},
			handlers: []*SlotSrcHandler{},
		}
	}

	for _, handler := range i.scriptExport.Handlers {
		code := handler.Code

		// if marked as lib then we place it in the libs folder
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

		} else {
			slotSrc := slots[handler.Filter.SlotKey]

			filterCall, err := srcutils.MakeFilterCallFromSignature(handler.Filter.Signature, handler.Filter.Args)
			if err != nil {
				return errors.WithStack(err)
			}

			// expand the code into lines, ignore 1 trailing blank line
			lines := strings.Split(strings.TrimSuffix(code, "\n"), "\n")

			// main code block in start() filter is special
			if filterCall == "start()" && lines[0] == "-- !DU: main" {
				// trim off the marker
				lines = lines[1:]
				// trim off 1 blank line
				if lines[0] == "" {
					lines = lines[1:]
				}
				slotSrc.mainCode = append(slotSrc.mainCode, lines...)
			} else {
				slotSrc.handlers = append(slotSrc.handlers, &SlotSrcHandler{
					code: lines,
					sig:  filterCall,
				})
			}
		}
	}

	for _, slotSrc := range slots {
		if len(slotSrc.handlers) == 0 {
			continue
		}

		slotPath := path.Join(outputDir, "slots", fmt.Sprintf("%d.%s.lua", slotSrc.key, slotSrc.name))

		out := make([]string, 0)

		// add main code block first
		out = append(out, slotSrc.mainCode...)
		out = append(out, "")

		// then add the handlers
		for _, handler := range slotSrc.handlers {
			// open our code block with `do` and its marker
			out = append(out, fmt.Sprintf("do -- !DU: %s", handler.sig))

			// indent the code @TODO: some crazy people like 2 spaces or tabs ...
			indented := make([]string, len(handler.code))
			for k, l := range handler.code {
				if l != "" {
					indented[k] = "    " + handler.code[k]
				}
			}

			out = append(out, indented...)

			// close the block
			out = append(out, fmt.Sprintf("end -- !DU: end"), "")
		}

		err = ioutil.WriteFile(slotPath, []byte(strings.Join(out, "\n")), 0666)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
