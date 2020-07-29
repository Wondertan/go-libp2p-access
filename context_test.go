package access

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Token(t *testing.T) {
	ctx := context.Background()

	tkn, err := GetToken(ctx)
	assert.Zero(t, tkn)
	assert.Equal(t, ErrNoToken, err)

	ctx = WithToken(ctx, "test")
	tkn, err = GetToken(ctx)
	assert.Nil(t, err, err)
	assert.NotZero(t, tkn)
}
