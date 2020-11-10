package access

import (
	"context"
	"errors"
	"time"
)

var eoh = errors.New("end of handling")

type timedErrorHandler struct {
	in    chan error
	reset chan context.Context
	tout  time.Duration
}

func newTimedErrorHandler(ctx context.Context, h func(error), tout time.Duration) *timedErrorHandler {
	teh := &timedErrorHandler{
		in:    make(chan error, 1),
		reset: make(chan context.Context, 1),
		tout:  tout,
	}

	go teh.run(ctx, h, tout)
	return teh
}

func (gnt *timedErrorHandler) In() chan<- error {
	return gnt.in
}

func (gnt *timedErrorHandler) Reset(ctx context.Context) {
	gnt.reset <- ctx
}

func (gnt *timedErrorHandler) run(ctx context.Context, h func(error), tout time.Duration) {
	defer h(eoh)
	t := time.NewTimer(tout)
	for {
		select {
		case ctx = <-gnt.reset:
			if !t.Stop() {
				<-t.C
			}
			t.Reset(tout)
		case <-t.C:
			return
		case err, ok := <-gnt.in:
			if !ok {
				return
			} else if err == nil {
				continue
			}

			h(err)
		case <-ctx.Done():
			h(ctx.Err())
			return
		}
	}
}
