package hsmlib

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthPrefix2B(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		in := []byte("")
		out := LengthPrefix2B(in)
		assert.Len(t, out, 2)
		assert.Equal(t, []byte{0, 0}, out)
	})
	t.Run("B2 Command", func(t *testing.T) {
		in := []byte("9999B20004text")
		out := LengthPrefix2B(in)
		assert.Len(t, out, len(in)+2)

		expected := append([]byte{0, byte(len(in))}, in...)
		assert.Equal(t, expected, out)
	})
	t.Run("Various lengths", func(t *testing.T) {
		for _, length := range []int{255, 256, 66534} {
			t.Run(fmt.Sprint(length), func(t *testing.T) {
				in := bytes.Repeat([]byte{0x99}, length)
				out := LengthPrefix2B(in)
				assert.Len(t, out, len(in)+2)

				expected := append([]byte{byte(len(in) / 256), byte(len(in))}, in...)
				assert.Equal(t, expected, out)
			})
		}
	})
}
