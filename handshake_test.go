package access

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandshake(t *testing.T) {
	in := Token("test")
	ctx := WithToken(context.Background(), in)

	var buf bytes.Buffer
	WriteToken(&buf, in)

	out, errs, err := TakeHand(NewPassingGranter(), &buf, "peer")
	require.NoError(t, err)
	require.NotNil(t, errs)
	require.Equal(t, in, out)

	err = GiveHand(ctx, &buf)
	require.NoError(t, err)
}
