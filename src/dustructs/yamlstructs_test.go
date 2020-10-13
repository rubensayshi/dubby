package dustructs

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/rubensayshi/dubby/src/utils"
	"github.com/stretchr/testify/require"
)

func TestMarshalYaml(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.json"))
	assert.NoError(err)

	export := &ScriptExport{}
	err = json.Unmarshal(f, export)
	assert.NoError(err)

	y, err := yaml.Marshal(export)
	assert.NoError(err)

	fexpected, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.conf"))
	assert.NoError(err)

	assert.Equal(string(fexpected), string(y))
}

func TestUnmarshalYaml(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.conf"))
	assert.NoError(err)

	export := &ScriptExport{}
	err = yaml.Unmarshal(f, export)
	assert.NoError(err)

	assert.Equal(0, len(export.Slots))
	assert.Equal(2, len(export.Handlers))
	assert.Equal("yeeehaaaa(\"tick\")", export.Handlers[1].Code)
}
