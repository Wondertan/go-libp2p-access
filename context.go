package access

import (
	"context"
	"errors"
)

var ErrNoToken = errors.New("access: no token")

type tokenKey struct{}

func WithToken(ctx context.Context, token Token) context.Context {
	if len(token) == 0 {
		return ctx
	}

	return context.WithValue(ctx, tokenKey{}, token)
}

func GetToken(ctx context.Context) (Token, error) {
	token, ok := ctx.Value(tokenKey{}).(Token)
	if !ok {
		return "", ErrNoToken
	}

	return token, nil
}
