package srcutils

import (
	"testing"

	"github.com/rubensayshi/duconverter/src/dustructs"
	"github.com/stretchr/testify/require"
)

func TestSignatureWithArgs(t *testing.T) {
	assert := require.New(t)

	{
		res, err := SignatureWithArgs("tick(timerId)", []dustructs.Arg{{"Live"}})
		assert.NoError(err)
		assert.Equal("tick(\"Live\")", res)
	}

	{
		res, err := SignatureWithArgs("tick(timerId, cookie)", []dustructs.Arg{{"Live"}, {"and Let Die"}})
		assert.NoError(err)
		assert.Equal("tick(\"Live\", \"and Let Die\")", res)
	}
}

func TestArgsFromFileHeader(t *testing.T) {
	assert := require.New(t)

	{
		fn, args, err := ArgsFromFileHeader("tick(\"redraw\")")
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(1, len(args))
		assert.Equal("redraw", args[0].Value)
	}

	{
		fn, args, err := ArgsFromFileHeader("tick(\"Live\")")
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(1, len(args))
		assert.Equal("Live", args[0].Value)
	}

	{
		fn, args, err := ArgsFromFileHeader("tick(\"Live\", \"and Let Die\")")
		assert.NoError(err)
		assert.Equal("tick", fn)
		assert.Equal(2, len(args))
		assert.Equal("Live", args[0].Value)
		assert.Equal("and Let Die", args[1].Value)
	}
}
