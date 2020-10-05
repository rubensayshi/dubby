package luamin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLuaMin(t *testing.T) {
	assert := require.New(t)

	out, err := LuaMin([]byte("function hiThere()\n--comment\nprint()end"))
	assert.NoError(err)
	assert.Equal("function hiThere()print()end\n", string(out))
}

func TestLuaMinErr(t *testing.T) {
	assert := require.New(t)

	_, err := LuaMin([]byte("14.blabla"))
	assert.Error(err)
	assert.Contains(err.Error(), "unexpected number '14' near 'blabla'")
}
