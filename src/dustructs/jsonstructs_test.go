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

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "src/dustructs", "testvector1.json"))
	assert.NoError(err)

	export := &ScriptExport{}
	err = json.Unmarshal(f, export)
	assert.NoError(err)

	assert.Equal(13, len(export.Slots))
	assert.Equal("library", export.Slots[SLOT_IDX_LIBRARY].Name)
	assert.Equal("hematite", export.Slots[3].Name)
	assert.Equal("unit.setTimer(\"Live\",1)\nswitch.activate()\n", export.Handlers[1].Code)
}
