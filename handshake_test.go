package access

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandshake(t *testing.T) {
	tkn := Token("test")
	ctx := WithToken(context.Background(), tkn)

	var buf bytes.Buffer
	WriteToken(&buf, tkn)

	errs, err := TakeHand(NewPassingGranter(), &buf, "peer")
	require.NoError(t, err)
	require.NotNil(t, errs)

	err = GiveHand(ctx, &buf)
	require.NoError(t, err)
}
