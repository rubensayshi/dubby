package srcutils

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	SLOT_IDX_UNIT    = -1
	SLOT_IDX_SYSTEM  = -2
	SLOT_IDX_LIBRARY = -3
)

type ElementType struct {
	ClassName string
	Filters   []string
}

func NewElementType(className string, filters ...string) ElementType {
	return ElementType{
		ClassName: className,
		Filters:   filters,
	}
}

// abstract, can be "inherited"
var enterable = NewElementType("", "enter(id)", "leave(id)")
var pressable = NewElementType("", "pressed()", "released()")

// @TODO: not sure if this should be a map or a list ... keys seem arbitrary?
var ElementTypes = map[string]ElementType{
	"antigravityGenerator": NewElementType("AntiGravityGeneratorUnit"),
	"core":                 NewElementType("CoreUnit"),
	"databank":             NewElementType("DataBank"),
	"fuelContainer":        NewElementType("FuelContainer"),
	"industry":             NewElementType("IndustryUnit", "completed()", "statusChanged(status)"),
	"gyro":                 NewElementType("GyroUnit"),
	"radar":                NewElementType("RadarUnit", enterable.Filters...),
	"pvpRadar":             NewElementType("RadarPVPUnit", enterable.Filters...),
	"screen":               NewElementType("ScreenUnit", "mouseDown(x,y)", "mouseUp(x,y)"),
	"laserDetector":        NewElementType("LaserDetectorUnit", "laserHit()", "laserRelease()"),
	"receiver":             NewElementType("ReceiverUnit", "receive(channel,message)"),
	"weapon":               NewElementType("WeaponUnit"),
	"control":              NewElementType("Control Unit", "start()", "stop()", "tick(timerId)"),
	"system":               NewElementType("System", "actionStart(action)", "actionStop(action)", "actionLoop(action)", "update()", "flush()"),
	"library":              NewElementType("Library"),
}

// FilterSignatures = {"tick": "tick(timerId)", ...}
var FilterSignatures = map[string]string{}

// FiltersBySignature = {"tick(timerId)": "tick", ...}
var FiltersBySignature = map[string]string{}

func init() {
	for _, elementType := range ElementTypes {
		for _, filterSignature := range elementType.Filters {
			fn, _, err := ParseFilterCall(filterSignature)
			if err != nil {
				panic(fmt.Sprintf("%+v", errors.WithStack(err)))
			}

			FilterSignatures[fn] = filterSignature
			FiltersBySignature[filterSignature] = fn
		}
	}
}
