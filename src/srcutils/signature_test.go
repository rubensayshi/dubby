package srcutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeHeader(t *testing.T) {
	assert := require.New(t)

	{
		res, err := MakeFilterCallFromSignature("tick(timerId)", []Arg{{"Live"}})
		assert.NoError(err)
		assert.Equal("tick([Live])", res)
	}

	{
		res, err := MakeFilterCallFromSignature("tick(timerId, cookie)", []Arg{{"Live"}, {"and Let Die"}})
		assert.NoError(err)
		assert.Equal("tick([Live, and Let Die])", res) // @TODO: how is this sane?
	}
}

func TestParseHeader(t *testing.T) {
	assert := require.New(t)

	{
		fn, args, err := ParseFilterCall("tick(\"redraw\")")
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(1, len(args))
		assert.Equal("redraw", args[0].Value)
	}

	for _, header := range []string{"tick(\"Live\")", "tick([Live])", "tick([\"Live\"])"} {
		fn, args, err := ParseFilterCall(header)
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(1, len(args))
		assert.Equal("Live", args[0].Value)
	}

	for _, header := range []string{"tick(\"Live\", \"and let Die\")", "tick([\"Live\", \"and let Die\"])"} {
		fn, args, err := ParseFilterCall(header)
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(2, len(args))
		assert.Equal("Live", args[0].Value)
		assert.Equal("and let Die", args[1].Value)
	}

	{
		fn, args, err := ParseFilterCall("tick([Live, and Let Die])")
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(2, len(args))
		assert.Equal("Live", args[0].Value)
		assert.Equal("and Let Die", args[1].Value)
	}

	{
		// @TODO: not sure about this case, how should these args really be parsed?
		_, _, err := ParseFilterCall("tick([\"Live, and Let Die\"])")
		assert.NoError(err)
		//assert.Equal("tick", fn)
		//assert.Equal(2, len(args))
		//assert.Equal("Live", args[0].Value)
		//assert.Equal("and Let Die", args[1].Value)
	}
}
