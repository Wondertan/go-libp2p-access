package access

import (
	"context"
	"errors"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

const Timeout = time.Minute * 15

// ErrNotGranted tells that access was not given to specified Topic.
var ErrNotGranted = errors.New("access: not granted")

// Token represents namespace/topic of some arbitrary process.
type Token string

// Granter tracks access permissions to arbitrary processes for network peers
// and provides convenient error handling.
type Granter interface {
	// Grant adds new access permission to some process defined with Token for specified set of peers.
	// Access for any peer is active until timeout is reached or context is closed.
	// All the following Grant calls over the same peers reset timeout for them.
	// Returns PeerError chan that is closed when active Token accesses for all peers ends.
	Grant(context.Context, Token, ...peer.ID) <-chan error

	// Granted checks if a peer has access to some process defines with Token.
	// On success returns chan to send errors happened within the process.
	// It is allowed to send multiple errors.
	// Closing returned chan ends access for peer earlier.
	Granted(Token, peer.ID) (chan<- error, error)
}
