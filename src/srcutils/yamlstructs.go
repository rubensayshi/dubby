package srcutils

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

const uniqueSlotsPadding = "DUSLOT__%s__DUSLOT%d"

// regexp to find the padding we added to make slots keys unique
var reUniqueSlotsPadding = regexp.MustCompile(`DUSLOT__(.+?)__DUSLOT([0-9]+)`)

const argsPadding = "DUARGS__%s__DUARGS"

// regexp to find the padding we added to args
var reArgsPadding = regexp.MustCompile(`['"]?DUARGS__(.+?)__DUARGS['"]?`)

type scriptExportYaml struct {
	Slots    map[string]*slotYaml `yaml:"slots"`
	Handlers yaml.MapSlice        `yaml:"handlers"`
}

type filterYaml struct {
	Args string `yaml:"args"`
	Code string `yaml:"lua"`
}

type slotYaml struct {
	Class  string  `yaml:"class"`
	Select *string `yaml:"select"`
}

func MarshalAutoConf(e *ScriptExport) ([]byte, error) {
	out, err := yaml.Marshal(e)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// strip the padding, autoconf isn't real yaml and can have duplicates
	out = reUniqueSlotsPadding.ReplaceAll(out, []byte("$1"))

	// strip the padding, autoconf isn't real yaml and can have duplicates
	out = reArgsPadding.ReplaceAll(out, []byte("$1"))

	return out, nil
}

func UnmarshalAutoConf(input []byte) (*ScriptExport, error) {
	e := &ScriptExport{}
	err := yaml.Unmarshal(input, e)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return e, nil
}

func (e *ScriptExport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tmp := &scriptExportYaml{}
	err := unmarshal(tmp)
	if err != nil {
		return errors.WithStack(err)
	}

	slots := make(map[int]*Slot, len(tmp.Slots))
	k := 0
	for slotName, slot := range tmp.Slots {
		k++

		slots[k] = &Slot{
			Name:     slotName,
			Type:     NewType(),
			AutoConf: NewSlotAutoConf(slot.Class),
		}

		if slot.Select != nil {
			slots[k].AutoConf.Select = slot.Select
		}
	}

	if slots[SLOT_IDX_UNIT] == nil {
		slots[SLOT_IDX_UNIT] = &Slot{
			Name:     "unit",
			Type:     NewType(),
			AutoConf: nil,
		}
	}
	if slots[SLOT_IDX_SYSTEM] == nil {
		slots[SLOT_IDX_SYSTEM] = &Slot{
			Name:     "system",
			Type:     NewType(),
			AutoConf: nil,
		}
	}

	if slots[SLOT_IDX_LIBRARY] == nil {
		slots[SLOT_IDX_LIBRARY] = &Slot{
			Name:     "library",
			Type:     NewType(),
			AutoConf: nil,
		}
	}

	slotKeyIdx := 0
	handlers := make([]*Handler, 0, len(tmp.Handlers))
	for _, s := range tmp.Handlers {
		slot := s.Key.(string)
		fmt.Printf("slot: %+v \n", slot)

		slotKey := 0
		if slot == "unit" {
			slotKey = SLOT_IDX_UNIT
		} else if slot == "system" {
			slotKey = SLOT_IDX_SYSTEM
		} else if slot == "library" {
			slotKey = SLOT_IDX_LIBRARY
		} else {
			slotKeyIdx++
			slotKey = slotKeyIdx
		}

		for _, ss := range s.Value.(yaml.MapSlice) {
			k := ss.Key.(string)

			args := []Arg{}
			lua := ""

			for _, v := range ss.Value.(yaml.MapSlice) {
				switch v.Key.(string) {
				case "args":
					value, ok := v.Value.([]interface{})
					if !ok {
						return errors.Errorf("unsupported type for args (%T)[%+v]", v.Value, v.Value)
					}

					args = make([]Arg, len(value))
					for i, a := range value {
						arg, ok := a.(string)
						if !ok {
							return errors.Errorf("unsupported type for arg (%T)[%+v]", a, a)
						}
						args[i] = Arg{Value: arg}
					}
				case "lua":
					lua = v.Value.(string)
				default:
					return errors.Errorf("unknown key [%s]", v.Key.(string))
				}
			}

			// the key should either by the filter name
			//  or a parsable signature of which the function is the  filter name
			fn := ""
			if FilterSignatures[k] != "" {
				fn = k
			} else {
				fnName, _, err := ParseFilterCall(k)
				if err != nil {
					return errors.WithStack(err)
				}
				fn = fnName
			}

			filter, ok := FilterSignatures[fn]
			if !ok {
				return errors.Errorf("Unknown filter [%s] (from %s)", fn, k)
			}

			handlers = append(handlers, &Handler{
				Code: lua,
				Filter: &Filter{
					Args:      args,
					Signature: filter,
					SlotKey:   slotKey,
				},
				Key: len(handlers) + 1,
			})
		}
	}

	e.Slots = slots
	e.Handlers = handlers

	return nil
}

func (e *ScriptExport) MarshalYAML() (interface{}, error) {
	slots := make(map[string]*slotYaml, len(e.Slots))
	handlers := make(yaml.MapSlice, 0, len(e.Slots))
	handlersBySlotKey := make(map[int]map[string]*filterYaml, len(e.Slots))

	slotKeys := make(sort.IntSlice, 0, len(slots))
	for k, _ := range e.Slots {
		slotKeys = append(slotKeys, k)
	}

	slotKeys.Sort()

	for _, slotKey := range slotKeys {
		slot := e.Slots[slotKey]

		if slot.AutoConf != nil {
			slots[slot.Name] = &slotYaml{
				Class:  slot.AutoConf.Class,
				Select: slot.AutoConf.Select,
			}
		}

		filters := make(map[string]*filterYaml)
		handlers = append(handlers, yaml.MapItem{
			Key:   slot.Name,
			Value: filters,
		})

		handlersBySlotKey[slotKey] = filters
	}

	for k, v := range e.Handlers {
		slotKey := v.Filter.SlotKey

		fn, _, err := ParseFilterCall(v.Filter.Signature)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		_, signatureIsValid := FilterSignatures[fn]
		if !signatureIsValid {
			return nil, errors.Errorf("invalid signature: %s (from %s)", fn, v.Filter.Signature)
		}

		// add padding to make the slots unique
		fn = fmt.Sprintf(uniqueSlotsPadding, fn, k)

		args := make([]string, len(v.Filter.Args))
		for i, arg := range v.Filter.Args {
			if strings.Contains(arg.Value, ",") {
				return nil, errors.Errorf("arg contains a `,`, which probably isn't allow ...")
			}
			args[i] = arg.Value
		}

		argsstr := "[]"
		if len(args) > 0 {
			argsstr = "[" + strings.Join(args, ",") + "]"
		}

		handlersBySlotKey[slotKey][fn] = &filterYaml{
			Args: fmt.Sprintf(argsPadding, argsstr),
			Code: v.Code,
		}
	}

	tmp := &scriptExportYaml{
		Slots:    slots,
		Handlers: handlers,
	}

	return tmp, nil
}
