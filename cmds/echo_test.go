package cmds

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeEcho(t *testing.T) {
	fakeMessage := gofakeit.HipsterSentence(5)

	var cmd Echo
	var cmdI hsmlib.Command
	cmd = MakeEcho(fakeMessage)

	assert.Equal(t, fakeMessage, cmd.Message)
	cmdI = cmd

	require.NotNil(t, cmdI)
	assert.Implements(t, (*hsmlib.Command)(nil), cmdI)

	assert.Equal(t, []byte("B2"), cmdI.Code())

	var expectedData bytes.Buffer
	fmt.Fprintf(&expectedData, "%04X", len(fakeMessage))
	expectedData.WriteString(fakeMessage)
	assert.Equal(t, expectedData.Bytes(), cmdI.Data())
}
