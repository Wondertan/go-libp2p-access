package access

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimedErrorHandler_In(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rtrn := make(chan error, 1)
	eh := newTimedErrorHandler(ctx, func(err error) {
		rtrn <- err
	}, time.Millisecond * 5)

	for range make([]bool, 5){
		eh.In() <- fmt.Errorf("test")
		<-rtrn
	}

	require.Equal(t, eoh, <-rtrn)
}

func TestTimedErrorHandler_Reset(t *testing.T) {
	ctx, newCncl := context.WithCancel(context.Background())

	rtrn := make(chan error, 1)
	eh := newTimedErrorHandler(ctx, func(err error) {
		rtrn <- err
	}, time.Millisecond * 15)


	for range make([]bool, 5){
		oldCncl := newCncl
		select {
		case <-time.After(time.Millisecond * 10):
			ctx, newCncl = context.WithCancel(context.Background())
			eh.Reset(ctx)
			oldCncl()
		case <-rtrn:
			t.Fatal("timer was not reset")
		}
	}

	newCncl()
	require.Equal(t, context.Canceled, <-rtrn)
	require.Equal(t, eoh, <-rtrn)
}
