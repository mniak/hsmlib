package cmds

import (
	"testing"

	"github.com/mniak/hsmlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeHealthcheck(t *testing.T) {
	var cmd Healthcheck
	var cmdI hsmlib.Command
	cmd = MakeHealthcheck()
	cmdI = cmd

	require.NotNil(t, cmdI)
	assert.Implements(t, (*hsmlib.Command)(nil), cmdI)

	assert.Equal(t, []byte("JK"), cmdI.Code())
	assert.Equal(t, []byte{}, cmdI.Data())
}
