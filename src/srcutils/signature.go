package srcutils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var reFilterCall = regexp.MustCompile(`^ *(?P<fn>[a-zA-Z0-9_-]+)\(\[?(?P<args>.*?)\]?\) *$`)

func MakeFilterCallFromSignature(signature string, args []Arg) (string, error) {
	res := reFilterCall.FindStringSubmatch(signature)

	// parse the args from the signature
	fn := res[1]
	argstr := res[2]
	var fnargs []string
	if strings.TrimSpace(argstr) == "" {
		fnargs = []string{}
	} else {
		fnargs = strings.Split(res[2], ",")
		for k, v := range fnargs {
			fnargs[k] = strings.TrimSpace(v)
		}
	}

	if len(fnargs) != len(args) {
		return "", errors.Errorf("Wrong number of args, expected %d: %s", len(fnargs), signature)
	}

	if len(fnargs) == 0 {
		return fmt.Sprintf("%s()", fn), nil
	}

	resArgs := make([]string, len(fnargs))
	for k, _ := range resArgs {
		resArgs[k] = args[k].Value
	}

	return fmt.Sprintf("%s([%s])", fn, strings.Join(resArgs, ", ")), nil
}

func MakeFilterCallFromFn(fn string, args []Arg) (string, error) {
	resArgs := make([]string, len(args))
	for k, _ := range resArgs {
		resArgs[k] = args[k].Value
	}

	return fmt.Sprintf("%s([%s])", fn, strings.Join(resArgs, ", ")), nil
}

func ParseFilterCall(filterCall string) (string, []Arg, error) {
	res := reFilterCall.FindStringSubmatch(filterCall)
	if res == nil || len(res) < 2 {
		return "", nil, errors.Errorf("Filter does not match expected pattern: %s", filterCall)
	}

	// parse the args from the signature
	fn := res[1]
	argstr := res[2]
	var fnargs []string
	if strings.TrimSpace(argstr) == "" {
		fnargs = []string{}
	} else {
		fnargs = strings.Split(res[2], ",")
		for k, v := range fnargs {
			fnargs[k] = strings.TrimSpace(v)
		}
	}

	args := make([]Arg, len(fnargs))

	for k, fnarg := range fnargs {
		if strings.HasPrefix(fnarg, `"`) && strings.HasSuffix(fnarg, `"`) {
			fnarg = fnarg[1 : len(fnarg)-1]
		}

		args[k] = Arg{Value: fnarg}
	}

	return fn, args, nil
}
