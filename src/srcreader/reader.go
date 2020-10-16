package srcreader

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
	"github.com/rubensayshi/dubby/src/luamin"
	"github.com/rubensayshi/dubby/src/srcutils"
)

var badHandlerStartRegexp = regexp.MustCompile(`^.*-- ?!DU:.*$`)
var handlerStartRegexp = regexp.MustCompile(`^(do)? *-- ?!DU: *((?P<fn>[a-zA-Z0-9_-]+)\(\[?(?P<args>.*?)\]?\)) *$`)
var handlerEndRegexp = regexp.MustCompile(`^(end)? *-- ?!DU: end *$`)

func Read(srcDir string) (*srcutils.ScriptExport, error) {
	r := NewSrcReader(srcDir, false)
	err := r.Read()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return r.scriptExport, nil
}

type SrcReader struct {
	srcDir       string
	minify       bool
	scriptExport *srcutils.ScriptExport
	report       *Report
}

type Report struct {
	SrcLen      int
	MinifiedLen int
}

func NewSrcReader(srcDir string, minify bool) *SrcReader {
	return &SrcReader{
		srcDir:       srcDir,
		minify:       minify,
		scriptExport: srcutils.NewScriptExport(),
		report:       &Report{},
	}
}

func (r *SrcReader) ScriptExport() *srcutils.ScriptExport {
	return r.scriptExport
}

func (r *SrcReader) Report() *Report {
	return r.report
}

func (r *SrcReader) Read() error {
	return r.readFromSrcDir(r.srcDir)
}

func (r *SrcReader) readFromSrcDir(dir string) error {
	slots, err := os.Stat(path.Join(dir, "slots"))
	if err != nil {
		return errors.WithStack(err)
	}
	if !slots.IsDir() {
		return errors.Errorf("slots is a file, expected a directory")
	}

	err = r.readAutoConfConfig(path.Join(dir, srcutils.AUTOCONF_CONFIGFILE))
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.readFromSlotsDir(path.Join(dir, slots.Name()))
	if err != nil {
		return errors.WithStack(err)
	}

	lib, err := os.Stat(path.Join(dir, "lib"))
	if !os.IsNotExist(err) {
		if err != nil {
			return errors.WithStack(err)
		}
		if !lib.IsDir() {
			return errors.Errorf("lib is a file, expected a directory")
		}

		err = r.readFromLibDir(path.Join(dir, lib.Name()))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *SrcReader) readAutoConfConfig(autoConfFile string) error {
	autoConfYml, err := ioutil.ReadFile(autoConfFile)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.WithStack(err)
	}

	autoConf := &srcutils.AutoConfConfig{}
	err = yaml.Unmarshal(autoConfYml, autoConf)
	if err != nil {
		return errors.WithStack(err)
	}

	r.scriptExport.AutoConfName = autoConf.Name

	slotKey := 0
	for k, slotConfig := range autoConf.Slots {
		slotKey++

		slot := srcutils.NewSlot(k)
		slot.AutoConf = &srcutils.SlotAutoConf{
			Class:  slotConfig.Class,
			Select: slotConfig.Select,
		}
		r.scriptExport.Slots[slotKey] = slot
	}

	return nil
}

func (r *SrcReader) readFromSlotsDir(slotsDir string) error {
	slotFiles, err := ioutil.ReadDir(slotsDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, slotFile := range slotFiles {
		slotFilePath := path.Join(slotsDir, slotFile.Name())

		if slotFile.IsDir() {
			return errors.Errorf("slotFile is a directory, expected a file: %s", slotFilePath)
		}

		s := strings.Split(slotFile.Name(), ".")
		if len(s) < 2 {
			return errors.Errorf("slotFile should its slot key as index [%s]", slotFile.Name())
		}
		slotKeyStr := s[0]
		slotName := strings.Join(s[1:], ".")

		slotKey, err := strconv.Atoi(slotKeyStr)
		if err != nil {
			return errors.WithStack(err)
		}

		// init the slot
		if r.scriptExport.Slots[slotKey] == nil {
			r.scriptExport.Slots[slotKey] = srcutils.NewSlot(slotName)
		}

		err = r.readFromSlotFile(slotFilePath, slotKey)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *SrcReader) readFromSlotFile(filePath string, slotKey int) error {
	handlers := make([]*srcutils.Handler, 0)

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.WithStack(err)
	}
	content := string(buf)
	// @TODO: for now, get rid of windows line endings, should be configurable ...
	content = strings.ReplaceAll(content, "\r\n", "\n")

	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		var handler *srcutils.Handler

		mainCode := make([]string, 0)
		handlerCode := make([]string, 0)

		for k, line := range lines {
			if handlerStartRegexp.MatchString(line) {
				header, err := extractHeaderFromLine(line)
				if err != nil {
					return errors.WithStack(err)
				}

				fnname, args, err := srcutils.ParseFilterCall(header)
				if err != nil {
					return errors.WithStack(err)
				}

				if srcutils.FilterSignatures[fnname] == "" {
					return errors.Errorf("unknown filter signature: [%s] -> [%s] (from: [%d] %s)",
						header, fnname, k, line)
				}

				signature := srcutils.FilterSignatures[fnname]
				header, _ = srcutils.MakeFilterCallFromSignature(signature, args)

				handler = &srcutils.Handler{
					Filter: &srcutils.Filter{
						Signature: signature,
						Args:      args,
						SlotKey:   slotKey,
					},
				}
			} else if handlerEndRegexp.MatchString(line) {
				if handler == nil {
					// @TODO: could be warning?
					return errors.Errorf("end marker without start: [%d][%s]", k, line)
				}

				// trim off any (consistent) indenting
				handlerCode = srcutils.TrimConsistentIndenting(handlerCode)

				code := strings.Join(handlerCode, "\n")
				r.report.SrcLen += len(code)
				if r.minify {
					minified, err := luamin.LuaMin([]byte(code))
					if err != nil {
						return errors.WithStack(err)
					}

					code = string(minified)
					r.report.MinifiedLen += len(code)
				}

				// flush handler
				handler.Code = code
				handlers = append(handlers, handler)

				// reset state
				handler = nil
				handlerCode = []string{}

			} else if badHandlerStartRegexp.MatchString(line) {
				// @TODO: could be warning?
				return errors.Errorf("bad marker: [%d][%s]", k, line)
			} else {
				// append code to handler or to main block
				if handler != nil {
					handlerCode = append(handlerCode, line)
				} else {
					mainCode = append(mainCode, line)
				}
			}
		}

		if handler != nil || len(handlerCode) > 0 {
			// @TODO: could be warning?
			return errors.Errorf("unclosed state")
		}

		// if we have a main block then we need to add a handler for it
		if len(mainCode) > 0 {
			justWhitelines := true
			for _, l := range mainCode {
				if l != "" {
					justWhitelines = false
					break
				}
			}

			if !justWhitelines {
				// main block needs marker
				mainCode = append([]string{"-- !DU: main"}, mainCode...)

				// trim of 2 trailing lines, these keep being added
				if mainCode[len(mainCode)-1] == "" {
					mainCode = mainCode[:len(mainCode)-1]
				}
				if mainCode[len(mainCode)-1] == "" {
					mainCode = mainCode[:len(mainCode)-1]
				}

				code := strings.Join(mainCode, "\n")
				r.report.SrcLen += len(code)
				if r.minify {
					minified, err := luamin.LuaMin([]byte(code))
					if err != nil {
						return errors.WithStack(err)
					}

					code = string(minified)
					r.report.MinifiedLen += len(code)
				}

				mainHandler := &srcutils.Handler{
					Code: code,
					Filter: &srcutils.Filter{
						Signature: "start()",
						Args:      []string{},
						SlotKey:   slotKey,
					},
				}
				handlers = append([]*srcutils.Handler{mainHandler}, handlers...)
			}
		}
	}

	// fix the handler keys
	key := len(r.scriptExport.Handlers)
	for _, handler := range handlers {
		key++
		handler.Key = key

		r.scriptExport.Handlers = append(r.scriptExport.Handlers, handler)

	}

	return nil
}

func (r *SrcReader) readFromLibDir(libDir string) error {
	files, err := ioutil.ReadDir(libDir)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(files) == 0 {
		return nil
	}

	libContent := make([]string, 0)

	for _, file := range files {
		filePath := path.Join(libDir, file.Name())

		if file.IsDir() {
			return errors.Errorf("file is a directory, expected a file: %s", filePath)
		}

		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			return errors.WithStack(err)
		}
		content := string(buf)
		content = strings.ReplaceAll(content, "\r\n", "\n")

		handlerName := strings.TrimSuffix(file.Name(), ".lua")
		libName := handlerName

		libContent = append(libContent, "-- !DU[lib]: "+libName+"\n\n"+content)
	}

	// shift all handlers 1 slot forward
	for key, handler := range r.scriptExport.Handlers {
		r.scriptExport.Handlers[key].Key = handler.Key + 1
	}

	code := strings.Join(libContent, "\n")
	r.report.SrcLen += len(code)
	if r.minify {
		minified, err := luamin.LuaMin([]byte(code))
		if err != nil {
			return errors.WithStack(err)
		}

		code = string(minified)
		r.report.MinifiedLen += len(code)
	}

	handler := &srcutils.Handler{
		Code: code,
		Filter: &srcutils.Filter{
			Args:      []string{},
			Signature: "start()",
			SlotKey:   srcutils.SLOT_IDX_UNIT,
		},
		Key: 1, // @TODO: maybe we can do 0?
	}

	r.scriptExport.Handlers = append([]*srcutils.Handler{handler}, r.scriptExport.Handlers...)

	return nil
}

func extractHeaderFromLine(line string) (string, error) {
	res := handlerStartRegexp.FindStringSubmatch(line)
	if res == nil || len(res) < 4 {
		return "", errors.Errorf("Header does not match expected pattern: %s", line)
	}

	return res[2], nil
}
