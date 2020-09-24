package dustructs

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/rubensayshi/duconverter/src/utils"
	"github.com/stretchr/testify/require"
)

func TestParseJson(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.json"))
	assert.NoError(err)

	export := &ScriptExport{}
	err = json.Unmarshal(f, export)
	assert.NoError(err)

	assert.Equal(1, len(export.Slots))
	assert.Equal("unit", export.Slots[SLOT_IDX_UNIT].Name)
	assert.Equal("unit.exit() -- ontick Live\n", export.Handlers[1].Code)
}
