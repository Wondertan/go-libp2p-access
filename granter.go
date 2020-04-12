package access

import (
	"context"
	"errors"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
)

var ErrNotGranted = errors.New("access: not granted")

// Token is a string which represents a namespace within which access is given.
type Token string

// Granter controls accesses between peers and tracks the errors occurred within the token.
// Useful when you need not only to manage access in many components but also to track the errors produced by them.
type Granter interface {
	// Grant gives access for specified peers.
	// Returns chan to track errors provided by the Granted.
	// NOTE: Grant renewal on the peer with the same token removes previous grant.
	Grant(context.Context, Token, ...peer.ID) <-chan error

	// Granted checks whenever access have been given for peer within the token.
	// On success returns chan to send errors to.
	Granted(Token, peer.ID) (chan<- error, error)
}

// granter implements Granter.
type granter struct {
	l      sync.Mutex
	grants map[Token]map[peer.ID]chan error
}

// NewGranter creates new Granter.
func NewGranter() Granter {
	return &granter{grants: make(map[Token]map[peer.ID]chan error)}
}

// Grant implements Granter.Grant.
func (g *granter) Grant(ctx context.Context, t Token, peers ...peer.ID) <-chan error {
	g.l.Lock()
	defer g.l.Unlock()

	tg, ok := g.grants[t]
	if !ok {
		g.grants[t] = make(map[peer.ID]chan error, 1)
		tg = g.grants[t]
	}

	out := make(chan error)
	wg := new(sync.WaitGroup)
	for _, p := range peers {
		in := make(chan error)
		wg.Add(1)
		go func(in chan error, p peer.ID) {
			defer wg.Done()
			select {
			case err := <-in:
				if err != nil {
					select {
					case out <- NewError(p, t, err): // notify client with error and peer the error happened to.
					case <-ctx.Done():
					}
				}
			case <-ctx.Done():
				select {
				case out <- NewError(p, t, ctx.Err()): // this allows checking exact peers that not finished exchange on context cancel.
				case <-ctx.Done():
				}
			}
		}(in, p)
		tg[p] = in
	}

	go func() {
		defer close(out) // closes out in case all peers are done with no errors.
		wg.Wait()
	}()

	return out
}

// Granted implements Granter.Granted.
func (g *granter) Granted(t Token, p peer.ID) (chan<- error, error) {
	g.l.Lock()
	defer g.l.Unlock()

	tg, ok := g.grants[t]
	if !ok {
		return nil, NewError(p, t, ErrNotGranted)
	}

	ch, ok := tg[p]
	if !ok {
		return nil, NewError(p, t, ErrNotGranted)
	}

	return ch, nil
}

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
