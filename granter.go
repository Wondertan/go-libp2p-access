package access

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
)

// granter implements Granter.
type granter struct {
	l      sync.Mutex
	grants map[Token]*grant
}

// NewGranter creates new Granter.
func NewGranter() Granter {
	return &granter{grants: make(map[Token]*grant)}
}

// Grant implements Granter.Grant
func (g *granter) Grant(ctx context.Context, tkn Token, ps ...peer.ID) <-chan error {
	g.l.Lock()
	defer g.l.Unlock()

	gnt, ok := g.grants[tkn]
	if !ok {
		gnt = &grant{
			tkn:   tkn,
			out:   make(chan error, len(ps)),
			acts:  make(chan action, 8),
			done:  make(chan struct{}),
			peers: make(map[peer.ID]*timedErrorHandler, len(ps)),
		}

		go func() {
			gnt.handle()
			g.l.Lock()
			delete(g.grants, tkn)
			g.l.Unlock()
		}()

		g.grants[tkn] = gnt
	}

	gnt.GrantPeers(ctx, ps...)
	return gnt.out
}

// Granted implements Granter.Granted.
func (g *granter) Granted(t Token, p peer.ID) (chan<- error, error) {
	g.l.Lock()
	defer g.l.Unlock()

	gnt, ok := g.grants[t]
	if !ok {
		return nil, NewError(p, t, ErrNotGranted)
	}

	ch := gnt.PeerGranted(p)
	if ch != nil {
		return ch, nil
	}

	return nil, NewError(p, t, ErrNotGranted)
}

type grant struct {
	tkn Token
	out chan error

	acts  chan action
	done  chan struct{}
	peers map[peer.ID]*timedErrorHandler
}

func (g *grant) GrantPeers(ctx context.Context, ps ...peer.ID) {
	select {
	case g.acts <- &add{ctx: ctx, ps: ps}:
	case <-g.done:
	}
}

func (g *grant) PeerGranted(p peer.ID) chan<- error {
	res := make(chan chan<- error)
	select {
	case g.acts <- &get{p: p, res: res}:
	case <-g.done:
	}

	select {
	case ch := <-res:
		return ch
	case <-g.done:
		return nil
	}
}

func (g *grant) handle() {
	defer func() {
		close(g.done)
		close(g.out)
	}()
	for {
		select {
		case act := <-g.acts:
			if act.handle(g) {
				return
			}
		}
	}
}

type action interface {
	handle(*grant) bool
}

type add struct {
	ctx context.Context
	ps  []peer.ID
}

func (act *add) handle(g *grant) bool {
	for _, p := range act.ps {
		if eh, ok := g.peers[p]; ok {
			eh.Reset(act.ctx)
			continue
		}

		p := p
		eh := newTimedErrorHandler(act.ctx, func(err error) {
			if err == eoh {
				g.acts <- &rm{p: p}
				return
			}

			g.out <- NewError(p, g.tkn, err)
		}, Timeout)
		g.peers[p] = eh
	}

	return false
}

type get struct {
	p   peer.ID
	res chan chan<- error
}

func (act *get) handle(g *grant) bool {
	if eh, ok := g.peers[act.p]; ok {
		act.res <- eh.In()
	}

	close(act.res)
	return false
}

type rm struct {
	p peer.ID
}

func (act *rm) handle(g *grant) bool {
	delete(g.peers, act.p)
	return len(g.peers) == 0
}
