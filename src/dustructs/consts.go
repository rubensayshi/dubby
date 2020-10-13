package dustructs

var FilterSignatures = map[string]string{
	"start":       "start()",
	"stop":        "stop()",
	"flush":       "flush()",
	"update":      "update()",
	"tick":        "tick(timerId)",
	"actionStart": "actionStart(action)",
	"actionStop":  "actionStop(action)",
	"actionLoop":  "actionLoop(action)",
}
