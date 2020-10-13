package dustructs

import (
	"github.com/pkg/errors"
)

type scriptExportYaml struct {
	Slots    map[string]*slotYaml              `yaml:"slots"`
	Handlers map[string]map[string]*filterYaml `yaml:"handlers"` // @TODO: technically we can have duplicates in DUs "yaml"
}

type filterYaml struct {
	Args []Arg  `yaml:"args"`
	Code string `yaml:"lua"`
}

type slotYaml struct {
	Class  string  `yaml:"class"`
	Select *string `yaml:"select"`
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

	slotKeyIdx := 0
	handlers := make([]*Handler, 0, len(tmp.Handlers))
	for slot, filters := range tmp.Handlers {
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

		for k, v := range filters {
			handlers = append(handlers, &Handler{
				Code: v.Code,
				Filter: &Filter{
					Args:      v.Args,
					Signature: FilterSignatures[k],
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
	handlers := make(map[string]map[string]*filterYaml, len(e.Slots))

	for _, v := range e.Slots {
		if v.AutoConf != nil {
			slots[v.Name] = &slotYaml{
				Class:  v.AutoConf.Class,
				Select: v.AutoConf.Select,
			}
		}

		handlers[v.Name] = make(map[string]*filterYaml)
	}

	for _, v := range e.Handlers {
		slot := e.Slots[v.Filter.SlotKey]

		fn := v.Filter.Signature // @TODO: sig to fn

		handlers[slot.Name][fn] = &filterYaml{
			Args: v.Filter.Args,
			Code: v.Code,
		}
	}

	tmp := &scriptExportYaml{
		Slots:    slots,
		Handlers: handlers,
	}

	return tmp, nil
}
