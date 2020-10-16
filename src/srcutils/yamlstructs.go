package srcutils

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const uniqueSlotsPadding = "DUSLOT__%s__DUSLOT%d"

// regexp to find the padding we added to make slots keys unique
var reUniqueSlotsPadding = regexp.MustCompile(`DUSLOT__(.+?)__DUSLOT([0-9]+)`)

type scriptExportYaml struct {
	Name     string        `yaml:"name"`
	Slots    yaml.MapSlice `yaml:"slots"`    // using MapSlice to maintain order
	Handlers yaml.MapSlice `yaml:"handlers"` // using MapSlice to maintain order
}

type filterYaml struct {
	Args []string `yaml:"args,omitempty,flow"`
	Code string   `yaml:"lua"`
}

type slotYaml struct {
	Class  string `yaml:"class"`
	Select string `yaml:"select,omitempty"`
}

func MarshalAutoConf(e *ScriptExport) ([]byte, error) {
	out, err := yaml.Marshal(e)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// strip the padding, autoconf isn't real yaml and can have duplicates
	out = reUniqueSlotsPadding.ReplaceAll(out, []byte("$1"))

	return out, nil
}

func UnmarshalAutoConf(input []byte) (*ScriptExport, error) {
	e := NewScriptExport()
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

	// slot counter, starting from 0 (note, -3, -2, -1 already exist)
	k := 0

	for _, v := range tmp.Slots {
		slotName := v.Key.(string)

		// juggling yaml back and forth is easier than digging through the nested yaml.MapItem
		slotRaw, err := yaml.Marshal(v.Value)
		if err != nil {
			return errors.WithStack(err)
		}
		slot := &slotYaml{}
		err = yaml.Unmarshal(slotRaw, slot)
		if err != nil {
			return errors.WithStack(err)
		}

		k++

		e.Slots[k] = &Slot{
			Name:     slotName,
			Type:     NewType(),
			AutoConf: NewSlotAutoConf(slot.Class),
		}

		if slot.Select != "" {
			e.Slots[k].AutoConf.Select = slot.Select
		}
	}

	slotKeyIdx := 0
	handlers := make([]*Handler, 0, len(tmp.Handlers))
	for _, s := range tmp.Handlers {
		slot := s.Key.(string)

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

			// juggle back and forth between yaml
			filterRaw, err := yaml.Marshal(ss.Value)
			if err != nil {
				return errors.WithStack(err)
			}
			filter := &filterYaml{}
			err = yaml.Unmarshal(filterRaw, filter)
			if err != nil {
				return errors.WithStack(err)
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

			signature, ok := FilterSignatures[fn]
			if !ok {
				return errors.Errorf("Unknown filter [%s] (from %s)", fn, k)
			}

			handlers = append(handlers, &Handler{
				Code: filter.Code,
				Filter: &Filter{
					Args:      filter.Args,
					Signature: signature,
					SlotKey:   slotKey,
				},
				Key: len(handlers) + 1,
			})
		}
	}

	e.AutoConfName = tmp.Name
	e.Handlers = handlers

	return nil
}

func NewAutoConfConfig(e *ScriptExport) *AutoConfConfig {
	conf := &AutoConfConfig{
		Name: e.AutoConfName,
	}

	slots := make(map[string]*slotYaml, len(e.Slots))

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
	}

	conf.Slots = slots

	return conf
}

func (e *ScriptExport) MarshalYAML() (interface{}, error) {
	slots := make(yaml.MapSlice, 0, len(e.Slots)-3)
	handlers := make(yaml.MapSlice, 0, len(e.Slots))
	handlersBySlotKey := make(map[int]map[string]*filterYaml, len(e.Slots))

	// map looses the sorting of our keys, so we need to resort the keys
	slotKeys := make(sort.IntSlice, 0, len(slots))
	for k, _ := range e.Slots {
		slotKeys = append(slotKeys, k)
	}
	slotKeys.Sort()

	for _, slotKey := range slotKeys {
		slot := e.Slots[slotKey]

		if slot.AutoConf != nil {
			slots = append(slots, yaml.MapItem{
				Key: slot.Name,
				Value: &slotYaml{
					Class:  slot.AutoConf.Class,
					Select: slot.AutoConf.Select,
				},
			})
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
			if strings.Contains(arg, ",") {
				return nil, errors.Errorf("arg contains a `,`, which probably isn't allow ...")
			}
			args[i] = arg
		}

		handlersBySlotKey[slotKey][fn] = &filterYaml{
			Args: args,
			Code: v.Code,
		}
	}

	// get rid if handlers without filters
	finalHandlers := make(yaml.MapSlice, 0, len(handlers))
	for _, v := range handlers {
		if len(v.Value.(map[string]*filterYaml)) > 0 {
			finalHandlers = append(finalHandlers, v)
		}
	}

	// alphabetical sort, for testing mostly ...
	sort.Slice(slots, func(i, j int) bool {
		return strings.Compare(slots[i].Key.(string), slots[j].Key.(string)) == -1
	})

	tmp := &scriptExportYaml{
		Name:     e.AutoConfName,
		Slots:    slots,
		Handlers: finalHandlers,
	}

	return tmp, nil
}
