package dustructs

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

const (
	SLOT_IDX_UNIT    = -1
	SLOT_IDX_SYSTEM  = -2
	SLOT_IDX_LIBRARY = -3
)

type ScriptExport struct {
	Slots    map[int]*Slot
	Handlers []Handler
	Methods  []Method
	Events   []Event
}

type scriptExportJson struct {
	Slots    map[string]Slot `json:"slots"` // keys are quoted numbers
	Handlers []Handler       `json:"handlers"`
	Methods  []Method        `json:"methods"`
	Events   []Event         `json:"events"`
}

func (e ScriptExport) UnmarshalJSON(d []byte) error {
	tmp := &scriptExportJson{}
	err := json.Unmarshal(d, tmp)
	if err != nil {
		return errors.WithStack(err)
	}

	slots := make(map[int]*Slot, len(tmp.Slots))
	for k, v := range tmp.Slots {
		kint, err := strconv.Atoi(k)
		if err != nil {
			return errors.WithStack(err)
		}

		slots[kint] = &v
	}

	e.Slots = slots
	e.Handlers = tmp.Handlers
	e.Methods = tmp.Methods
	e.Events = tmp.Events

	return nil
}

func (e *ScriptExport) MarshalJSON() ([]byte, error) {
	slots := make(map[string]Slot, len(e.Slots))
	for k, v := range e.Slots {
		kstr := strconv.Itoa(k)

		slots[kstr] = *v
	}

	tmp := &scriptExportJson{
		Slots:    slots,
		Handlers: e.Handlers,
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
	Type Type   `json:"type"`
}

type Type struct {
	Events  []Event  `json:"events"`
	Methods []Method `json:"methods"`
}

type Event struct {
}

type Method struct {
}

type Handler struct {
	Code   string `json:"code"`
	Filter Filter `json:"filter"`
	Key    int    `json:"key,string"`
}

type Filter struct {
	Args      []Arg  `json:"args"`
	Signature string `json:"signature"`
	SlotKey   int    `json:"slotKey,string"`
}

type Arg struct {
	Value string `json:"value"`
}
