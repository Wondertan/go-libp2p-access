package access

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerError is used to link a peer and error within the Token.
type PeerError struct {
	Peer  peer.ID
	Token Token
	Inner error
}

func NewError(peer peer.ID, token Token, err error) error {
	return &PeerError{Peer: peer, Token: token, Inner: err}
}

func (e *PeerError) Error() string {
	return fmt.Errorf("for %s with token(%s): %w", e.Peer.ShortString(), e.Token, e.Inner).Error()
}
