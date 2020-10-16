package srcwriter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rubensayshi/dubby/src/srcutils"

	"github.com/rubensayshi/dubby/src/utils"
	"github.com/stretchr/testify/require"
)

func TestSrcWriter_JsonWriteTo1(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors/testvector1", "input.json"))
	assert.NoError(err)

	export := &srcutils.ScriptExport{}
	err = json.Unmarshal(f, export)
	assert.NoError(err)

	dir, err := ioutil.TempDir(path.Join(utils.ROOT, "tmp"), "test")
	assert.NoError(err)
	defer os.RemoveAll(dir) // always cleanup the mess

	w := NewSrcWriter(export)
	err = w.WriteTo(dir)
	assert.NoError(err)

	actualDir := dir
	expectedDir := path.Join(utils.ROOT, "testvectors/testvector1", "output")

	checkActualDir(assert, actualDir, expectedDir)
	checkExpectedDir(assert, actualDir, expectedDir, "autoconf.yml")
}

func TestSrcWriter_YamlWriteTo1(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors/testvector1", "input.conf"))
	assert.NoError(err)

	export, err := srcutils.UnmarshalAutoConf(f)
	assert.NoError(err)

	dir, err := ioutil.TempDir(path.Join(utils.ROOT, "tmp"), "test")
	assert.NoError(err)
	defer os.RemoveAll(dir) // always cleanup the mess

	w := NewSrcWriter(export)
	err = w.WriteTo(dir)
	assert.NoError(err)

	actualDir := dir
	expectedDir := path.Join(utils.ROOT, "testvectors/testvector1", "output")

	checkActualDir(assert, actualDir, expectedDir)
	checkExpectedDir(assert, actualDir, expectedDir)
}

func checkActualDir(assert *require.Assertions, actualDir string, expectedDir string) {
	actualFiles, err := ioutil.ReadDir(actualDir)
	assert.NoError(err)

	for _, actualFile := range actualFiles {
		actualPath := path.Join(actualDir, actualFile.Name())
		expectedPath := path.Join(expectedDir, actualFile.Name())

		if !fileExists(expectedPath) {
			assert.FailNowf("file not expected", actualPath)
		} else {
			if actualFile.IsDir() {
				checkActualDir(assert, actualPath, expectedPath)
			} else {
				actualContent, err := ioutil.ReadFile(actualPath)
				assert.NoError(err)
				expectedContent, err := ioutil.ReadFile(expectedPath)
				assert.NoError(err)

				assert.Equal(string(expectedContent), string(actualContent))
			}
		}
	}
}

func checkExpectedDir(assert *require.Assertions, actualDir string, expectedDir string, ignoreFiles ...string) {
	expectedFiles, err := ioutil.ReadDir(expectedDir)
	assert.NoError(err)

	for _, expectedFile := range expectedFiles {
		ignore := false
		for _, ignoreFile := range ignoreFiles {
			if expectedFile.Name() == ignoreFile {
				ignore = true
				continue
			}
		}
		if ignore {
			continue
		}

		actualPath := path.Join(actualDir, expectedFile.Name())
		expectedPath := path.Join(expectedDir, expectedFile.Name())

		if !fileExists(actualPath) {
			assert.Failf("file expected", expectedPath)
		} else {
			if expectedFile.IsDir() {
				checkExpectedDir(assert, actualPath, expectedPath)
			} else {
				actualContent, err := ioutil.ReadFile(actualPath)
				assert.NoError(err)
				expectedContent, err := ioutil.ReadFile(expectedPath)
				assert.NoError(err)

				assert.Equal(string(expectedContent), string(actualContent))
			}
		}
	}
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
