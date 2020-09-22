package srcutils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rubensayshi/duconverter/src/dustructs"
)

var sigRegex = regexp.MustCompile(`^(?P<fn>[a-zA-Z0-9_-]+)\((?P<args>.*)\)$`)

func SignatureWithArgs(signature string, args []dustructs.Arg) (string, error) {
	res := sigRegex.FindStringSubmatch(signature)

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

	resArgs := make([]string, len(fnargs))
	for k, _ := range resArgs {
		resArgs[k] = fmt.Sprintf("\"%s\"", args[k].Value)
	}

	return fmt.Sprintf("%s(%s)", fn, strings.Join(resArgs, ", ")), nil
}

func ArgsFromFileHeader(header string) (string, []dustructs.Arg, error) {
	res := sigRegex.FindStringSubmatch(header)

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

	args := make([]dustructs.Arg, len(fnargs))

	for k, fnarg := range fnargs {
		if strings.HasPrefix(fnarg, `"`) && strings.HasSuffix(fnarg, `"`) {
			fnarg = fnarg[1 : len(fnarg)-1]
		}

		args[k] = dustructs.Arg{Value: fnarg}
	}

	return fn, args, nil
}
