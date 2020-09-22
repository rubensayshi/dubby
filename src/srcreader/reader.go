package srcreader

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rubensayshi/duconverter/src/dustructs"
	"github.com/rubensayshi/duconverter/src/srcutils"
)

func Read(srcDir string) (*dustructs.ScriptExport, error) {
	r := NewSrcReader(srcDir)
	err := r.Read()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return r.scriptExport, nil
}

type SrcReader struct {
	srcDir       string
	scriptExport *dustructs.ScriptExport
}

func NewSrcReader(srcDir string) *SrcReader {
	return &SrcReader{
		srcDir: srcDir,
		scriptExport: &dustructs.ScriptExport{
			Slots:    make(map[int]*dustructs.Slot),
			Handlers: make([]dustructs.Handler, 0),
			Methods:  make([]dustructs.Method, 0),
			Events:   make([]dustructs.Event, 0),
		},
	}
}

func (i *SrcReader) Read() error {
	return readSrcDirInto(i.srcDir, i.scriptExport)
}

func readSrcDirInto(dir string, scriptExport *dustructs.ScriptExport) error {
	slots, err := os.Stat(path.Join(dir, "slots"))
	if err != nil {
		return errors.WithStack(err)
	}
	if !slots.IsDir() {
		return errors.Errorf("slots is a file, expected a directory")
	}

	err = readSlotsDirInto(path.Join(dir, slots.Name()), scriptExport)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func readSlotsDirInto(slotsDir string, scriptExport *dustructs.ScriptExport) error {
	slotDirs, err := ioutil.ReadDir(slotsDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, slotDir := range slotDirs {
		slotDirPath := path.Join(slotsDir, slotDir.Name())

		if !slotDir.IsDir() {
			return errors.Errorf("slotDir is a file, expected a directory: %s", slotDirPath)
		}

		s := strings.Split(slotDir.Name(), ".")
		slotKeyStr := s[0]
		slotName := strings.Join(s[1:], ".")

		slotKey, err := strconv.Atoi(slotKeyStr)
		if err != nil {
			return errors.WithStack(err)
		}

		// init the slot
		if scriptExport.Slots[slotKey] == nil {
			scriptExport.Slots[slotKey] = &dustructs.Slot{
				Name: slotName,
				Type: dustructs.Type{
					Methods: make([]dustructs.Method, 0),
					Events:  make([]dustructs.Event, 0),
				},
			}
		}

		err = readSlotDirInto(slotDirPath, slotKey, scriptExport)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func readSlotDirInto(slotDir string, slotKey int, scriptExport *dustructs.ScriptExport) error {
	files, err := ioutil.ReadDir(slotDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		filePath := path.Join(slotDir, file.Name())

		if file.IsDir() {
			return errors.Errorf("file is a directory, expected a file: %s", filePath)
		}

		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			return errors.WithStack(err)
		}
		content := string(buf)

		handlerName := strings.TrimSuffix(file.Name(), ".lua")
		s := strings.Split(handlerName, ".")
		keyStr := s[0]
		signature := strings.Join(s[1:], ".")

		key, err := strconv.Atoi(keyStr)
		if err != nil {
			return errors.WithStack(err)
		}

		var args []dustructs.Arg

		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			if strings.HasPrefix(lines[0], "-- !DU: ") {
				content = strings.Join(lines[1:], "\n")

				header := strings.TrimPrefix(lines[0], "-- !DU: ")

				_, argsFromHeader, err := srcutils.ArgsFromFileHeader(header)
				if err != nil {
					return errors.WithStack(err)
				}

				args = argsFromHeader
			}

			handler := dustructs.Handler{
				Code: content,
				Filter: dustructs.Filter{
					Args:      args,
					Signature: signature,
					SlotKey:   slotKey,
				},
				Key: key,
			}

			scriptExport.Handlers = append(scriptExport.Handlers, handler)
		}
	}

	return nil
}
