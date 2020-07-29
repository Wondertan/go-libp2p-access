package access

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadWriteToken(t *testing.T) {
	var buf bytes.Buffer
	in := Token("test")

	n1, err := WriteToken(&buf, in)
	require.NoError(t, err)
	require.NotEqual(t, len(in), n1)

	out, n2, err := ReadToken(&buf)
	require.NoError(t, err)
	require.NotEqual(t, len(in), n2)

	assert.Equal(t, n1, n2)
	assert.Equal(t, in, out)
}
