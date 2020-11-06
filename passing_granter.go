package access

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

// passingGranter implements Granter which automatically allows access to everything.
type passingGranter struct{}

// NewPassingGranter builds new Granter which automatically allows access to everything.
func NewPassingGranter() Granter {
	return &passingGranter{}
}

// Granted implements Granter.Granted.
func (p *passingGranter) Grant(context.Context, Token, ...peer.ID) <-chan error {
	ch := make(chan error)
	close(ch)
	return ch
}

// Granted implements Granter.Granted.
func (p *passingGranter) Granted(Token, peer.ID) (chan<- error, error) {
	ch := make(chan error, 1)
	return ch, nil
}
