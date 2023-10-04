package hsmlib

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceiveResponse_MustExpectHeader(t *testing.T) {
	inputBytes, err := hex.DecodeString("0008F1F2F3F4C1C2E1E2")
	require.NoError(t, err)

	buf := bytes.NewReader(inputBytes)

	respWithHeader, err := ReceiveResponse(buf)
	require.NoError(t, err)

	assert.Equal(t, []byte{0xF1, 0xF2, 0xF3, 0xF4}, respWithHeader.Header)
}
