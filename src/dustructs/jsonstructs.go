package dustructs

import (
	"encoding/json"
	"strconv"
	"fmt"

	"github.com/pkg/errors"
)

const (
	SLOT_IDX_UNIT    = -1
	SLOT_IDX_SYSTEM  = -2
	SLOT_IDX_LIBRARY = -3
)

type ScriptExport struct {
	Slots    map[int]*Slot
	Handlers []*Handler
	Methods  []*Method
	Events   []*Event
}

func NewScriptExport() *ScriptExport {
	s := &ScriptExport{
		Slots:    make(map[int]*Slot),
		Handlers: make([]*Handler, 0),
		Methods:  make([]*Method, 0),
		Events:   make([]*Event, 0),
	}

	s.Slots[SLOT_IDX_UNIT] = NewSlot("unit")
	s.Slots[SLOT_IDX_SYSTEM] = NewSlot("system")
	s.Slots[SLOT_IDX_LIBRARY] = NewSlot("library")

	return s
}

type scriptExportJson struct {
	Slots    map[string]*Slot `json:"slots"` // keys are quoted numbers
	Handlers []*handlerRaw    `json:"handlers"`
	Methods  []*Method        `json:"methods"`
	Events   []*Event         `json:"events"`
}

func (e *ScriptExport) UnmarshalJSON(d []byte) error {
	tmp := &scriptExportJson{}
	err := json.Unmarshal(d, tmp)
	if err != nil {
		return errors.WithStack(err)
	}

	slots := make(map[int]*Slot, len(tmp.Slots))
	for k, v := range tmp.Slots {
		v := v // we're referencing this so need to declare inside the loop
		kint, err := strconv.Atoi(k)
		if err != nil {
			return errors.WithStack(err)
		}

		slots[kint] = v
	}

	handlers := make([]*Handler, len(tmp.Handlers))
	for k, v := range tmp.Handlers {
		slotKey, _ := v.Filter.SlotKey.Int64()
		key, _ := v.Key.Int64()

		handlers[k] = &Handler{
			Code: v.Code,
			Filter: &Filter{
				Args: v.Filter.Args,
				Signature: v.Filter.Signature, 
				SlotKey: int(slotKey),
			}, 
			Key: int(key),
		}
	}

	e.Slots = slots
	e.Handlers = handlers
	e.Methods = tmp.Methods
	e.Events = tmp.Events

	return nil
}

func (e *ScriptExport) MarshalJSON() ([]byte, error) {
	slots := make(map[string]*Slot, len(e.Slots))
	for k, v := range e.Slots {
		kstr := strconv.Itoa(k)
		slots[kstr] = v
	}

	handlers := make([]*handlerRaw, len(e.Handlers))
	for k, v := range e.Handlers {
		handlers[k] = &handlerRaw{
			Code: v.Code,
			Filter: &filterRaw{
				Args: v.Filter.Args,
				Signature: v.Filter.Signature,
				SlotKey: json.Number(fmt.Sprintf("%d", v.Filter.SlotKey)),
			},
			Key: json.Number(fmt.Sprintf("%d", v.Key)),
		}
	}

	tmp := &scriptExportJson{
		Slots:    slots,
		Handlers: handlers,
		Methods:  e.Methods,
		Events:   e.Events,
	}
	res, err := json.Marshal(tmp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

type Slot struct {
	Name string `json:"name"`
	Type *Type  `json:"type"`
}

func NewSlot(name string) *Slot {
	return &Slot{
		Name: name,
		Type: NewType(),
	}
}

type Type struct {
	Events  []Event  `json:"events"`
	Methods []Method `json:"methods"`
}

func NewType() *Type {
	return &Type{
		Events:  make([]Event, 0),
		Methods: make([]Method, 0),
	}
}

type Event struct {
}

type Method struct {
}

type Handler struct {
	Code   string `json:"code"`
	Filter *Filter `json:"filter"`
	Key    int    `json:"key,string"`
}

type handlerRaw struct {
	Code   string `json:"code"`
	Filter *filterRaw `json:"filter"`
	Key    json.Number    `json:"key"` // can be quoted and unquoted
}

type Filter struct {
	Args      []Arg  `json:"args"`
	Signature string `json:"signature"`
	SlotKey   int    `json:"slotKey,string"`
}

type filterRaw struct {
	Args      []Arg  `json:"args"`
	Signature string `json:"signature"`
	SlotKey   json.Number    `json:"slotKey"`  // can be quoted and unquoted
}

type Arg struct {
	Value string `json:"value"`
}
