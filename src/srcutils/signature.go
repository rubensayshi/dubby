package srcutils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var reFilterCall = regexp.MustCompile(`^ *(?P<fn>[a-zA-Z0-9_-]+)\(\[?(?P<args>.*?)\]?\) *$`)

func MakeFilterCallFromSignature(signature string, args []string) (string, error) {
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

	return MakeFilterCallFromFn(fn, args)
}

func MakeFilterCallFromFn(fn string, args []string) (string, error) {
	if len(args) == 0 {
		return fmt.Sprintf("%s()", fn), nil
	}

	// marshall the args as if they're yaml, because that's the format they should be in
	fnargs := struct {
		Args []string `yaml:"args,flow"`
	}{
		Args: args,
	}
	out, err := yaml.Marshal(fnargs)
	if err != nil {
		return "", errors.WithStack(err)
	}

	argsstr := strings.TrimSuffix(strings.TrimPrefix(string(out), "args: "), "\n")

	return fmt.Sprintf("%s(%s)", fn, argsstr), nil
}

func ParseFilterCall(filterCall string) (string, []string, error) {
	// extract the args from the call
	res := reFilterCall.FindStringSubmatch(filterCall)
	if res == nil || len(res) < 2 {
		return "", nil, errors.Errorf("Filter does not match expected pattern: %s", filterCall)
	}

	// parse the args from the signature
	fn := res[1]
	argstr := res[2]

	// wrap args in [] if they're not already, because we're gonna parse it as a yaml list
	if argstr != "" && !strings.HasPrefix(argstr, "[") {
		argstr = "[" + argstr + "]"
	}

	// parse the args as if they're yaml, because that's the format they should be in
	fnargs := struct {
		Args []string `yaml:"args,flow"`
	}{}
	err := yaml.Unmarshal([]byte("args: "+argstr), &fnargs)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return fn, fnargs.Args, nil
}
