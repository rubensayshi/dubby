package srcreader

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/rubensayshi/dubby/src/srcutils"
	"github.com/rubensayshi/dubby/src/utils"
	"github.com/stretchr/testify/require"
)

func TestSrcReader_Read1(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors/testvector1", "input.json"))
	assert.NoError(err)

	expected := &srcutils.ScriptExport{}
	err = json.Unmarshal(f, expected)
	assert.NoError(err)

	actual, err := Read(path.Join(utils.ROOT, "testvectors/testvector1", "output"))
	assert.NoError(err)

	assert.Equal(expected, actual)
}
