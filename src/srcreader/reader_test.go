package srcreader

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/rubensayshi/dubby/src/srcutils"
	"github.com/rubensayshi/dubby/src/utils"
	"github.com/stretchr/testify/require"
)

func TestSrcReader_Read1_Json(t *testing.T) {
	assert := require.New(t)

	actual, err := Read(path.Join(utils.ROOT, "testvectors/testvector1", "output"))
	assert.NoError(err)

	actualJson, err := json.MarshalIndent(actual, "", "  ")
	assert.NoError(err)

	expected, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors/testvector1", "input.json"))
	assert.NoError(err)
	expectedJson := strings.TrimSuffix(string(expected), "\n")

	assert.Equal(expectedJson, string(actualJson))
}

func TestSrcReader_Read1_Yaml(t *testing.T) {
	assert := require.New(t)

	actual, err := Read(path.Join(utils.ROOT, "testvectors/testvector1", "output"))
	assert.NoError(err)

	actualYaml, err := srcutils.MarshalAutoConf(actual)
	assert.NoError(err)

	expected, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors/testvector1", "input.conf"))
	assert.NoError(err)

	assert.Equal(string(expected), string(actualYaml))
}
