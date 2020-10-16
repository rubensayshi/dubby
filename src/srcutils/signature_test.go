package srcutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeFilterCallFromFn(t *testing.T) {
	for _, tt := range []struct {
		fn           string
		args         []string
		expectedCall string
	}{
		{"tick", []string{}, "tick()"},
		{"tick", []string{"redraw"}, "tick([redraw])"},
		{"tick", []string{"Live", "and let Die"}, "tick([Live, and let Die])"},
		{"tick", []string{"Live", "and, let Die"}, "tick([Live, 'and, let Die'])"},
	} {
		tt := tt
		t.Run(tt.expectedCall, func(t *testing.T) {
			assert := require.New(t)

			call, err := MakeFilterCallFromFn(tt.fn, tt.args)
			assert.NoError(err)
			assert.Equal(tt.expectedCall, call)
		})
	}
}

func TestMakeFilterCallFromSignature(t *testing.T) {
	for _, tt := range []struct {
		signature    string
		args         []string
		expectedCall string
	}{
		{"tick(timerId)", []string{"redraw"}, "tick([redraw])"},
		{"tick(timerId, smtsmt)", []string{"Live", "and let Die"}, "tick([Live, and let Die])"},
		{"tick(timerId, smtsmt)", []string{"Live", "and, let Die"}, "tick([Live, 'and, let Die'])"},
	} {
		tt := tt
		t.Run(tt.expectedCall, func(t *testing.T) {
			assert := require.New(t)

			call, err := MakeFilterCallFromSignature(tt.signature, tt.args)
			assert.NoError(err)
			assert.Equal(tt.expectedCall, call)
		})
	}
}

func TestMakeFilterCallFromSignatureErrs(t *testing.T) {
	for _, tt := range []struct {
		signature string
		args      []string
	}{
		{"tick(timerId)", []string{"Live", "and let Die"}},
		{"tick(timerId, smtsmt)", []string{"Live"}},
	} {
		tt := tt
		t.Run(tt.signature, func(t *testing.T) {
			assert := require.New(t)

			_, err := MakeFilterCallFromSignature(tt.signature, tt.args)
			assert.Error(err)
		})
	}
}

func TestParseFilterCall(t *testing.T) {
	for _, tt := range []struct {
		call         string
		expectedFn   string
		expectedArgs []string
	}{
		{"tick()", "tick", nil},
		{"tick(redraw)", "tick", []string{"redraw"}},
		{"tick(\"redraw\")", "tick", []string{"redraw"}},
		{"tick([redraw])", "tick", []string{"redraw"}},
		{"tick([\"redraw\"])", "tick", []string{"redraw"}},
		{"tick(Live, and let Die)", "tick", []string{"Live", "and let Die"}},
		{"tick([Live, and let Die])", "tick", []string{"Live", "and let Die"}},
		{"tick([Live, \"and, let Die\"])", "tick", []string{"Live", "and, let Die"}},
		{"tick([Live, 'and, let Die'])", "tick", []string{"Live", "and, let Die"}},
	} {
		tt := tt
		t.Run(tt.call, func(t *testing.T) {
			assert := require.New(t)
			fn, args, err := ParseFilterCall(tt.call)
			assert.NoError(err)
			assert.Equal(tt.expectedFn, fn)
			assert.EqualValues(tt.expectedArgs, args)
		})
	}
}
