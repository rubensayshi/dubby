package luamin

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const LUAMIN_CMD = "luamin"

var versionRegex = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

var isSupported *bool = nil

func IsSupported() bool {
	if isSupported == nil {
		cmd := exec.Command(LUAMIN_CMD, "-v")
		res, _ := cmd.Output()

		v := strings.TrimSuffix(strings.TrimSuffix(string(res), "\n"), "\r\n")

		ok := versionRegex.MatchString(v)
		isSupported = &ok
	}

	return *isSupported
}

func LuaMin(lua []byte) ([]byte, error) {
	if !IsSupported() {
		return nil, errors.Errorf("LuaMin not supported, couldn't detect `%s`", LUAMIN_CMD)
	}

	cmd := exec.Command(LUAMIN_CMD, "-c")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, _ := cmd.StdinPipe()
	err := cmd.Start()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = stdin.Write(lua)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_ = stdin.Close()

	err = cmd.Wait()

	if err != nil {
		_, isExit := err.(*exec.ExitError)
		if isExit && stderr.Len() > 0 {
			err = errors.Errorf(stderr.String())
		} else if isExit && stdout.Len() > 0 {
			err = errors.Errorf(stdout.String())
		}

		tmp, _ := ioutil.TempFile(os.TempDir(), "luamin")
		tmp.Write(lua)

		return nil, errors.Wrapf(err, "Failed to luamin: (dumped to %s)", tmp.Name())
	}

	return stdout.Bytes(), nil
}
