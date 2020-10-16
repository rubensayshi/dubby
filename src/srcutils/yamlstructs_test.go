package srcutils

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/rubensayshi/dubby/src/utils"
	"github.com/stretchr/testify/require"
)

func TestMarshalAutoConf(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.conf"))
	assert.NoError(err)

	export, err := UnmarshalAutoConf(f)
	assert.NoError(err)

	y, err := MarshalAutoConf(export)
	assert.NoError(err)

	fexpected, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.conf"))
	assert.NoError(err)

	assert.Equal(string(fexpected), string(y))
}

func TestUnmarshalAutoConf(t *testing.T) {
	assert := require.New(t)

	f, err := ioutil.ReadFile(path.Join(utils.ROOT, "testvectors", "testvector1", "input.conf"))
	assert.NoError(err)

	export, err := UnmarshalAutoConf(f)
	assert.NoError(err)

	assert.Equal(5, len(export.Slots))
	assert.Equal(4, len(export.Handlers))
	assert.Equal("yeeehaaaa(\"tick\")", export.Handlers[3].Code)
}
