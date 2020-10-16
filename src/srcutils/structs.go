package srcutils

import (
	"gopkg.in/yaml.v2"
)

type ScriptExport struct {
	AutoConfName string // when not "" means this is auto conf
	Slots        map[int]*Slot
	Handlers     []*Handler
	Methods      []*Method
	Events       []*Event
}

var _ yaml.Marshaler = &ScriptExport{}
var _ yaml.Unmarshaler = &ScriptExport{}

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

type Slot struct {
	Name     string        `json:"name"`
	Type     *Type         `json:"type"`
	AutoConf *SlotAutoConf `json:"-"`
}

func NewSlot(name string) *Slot {
	return &Slot{
		Name: name,
		Type: NewType(),
	}
}

type SlotAutoConf struct {
	Class  string
	Select string
}

func NewSlotAutoConf(class string) *SlotAutoConf {
	return &SlotAutoConf{
		Class: class,
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
	Code   string  `json:"code"`
	Filter *Filter `json:"filter"`
	Key    int     `json:"key,string"`
}

type Filter struct {
	Args      []Arg  `json:"args"`
	Signature string `json:"signature"`
	SlotKey   int    `json:"slotKey,string"`
}

type Arg struct {
	Value string `json:"value"`
}

type AutoConfConfig struct {
	Name  string               `yaml:"name"`
	Slots map[string]*slotYaml `yaml:"slots"`
}
