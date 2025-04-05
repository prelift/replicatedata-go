package causaltest

import (
	"context"

	"github.com/prelift/replicateddata-go/causal"
)

// A Pipe is a pair of remote peers exchanging messages in-memory.
//
// A single end of the pipe is not safe to use from multiple goroutines concurrently.
//
// The two ends cannot be used from the same goroutine.
// Send and Recv block, forcing the two goroutines to work in lock-step.
func Pipe[Operation any]() (_, _ causal.RemotePeer[Operation]) {

	inLeft, inRight := make(chan causal.SyncMsg[Operation]), make(chan causal.SyncMsg[Operation])

	left := &_PipeEnd[Operation]{inLeft, inRight}
	right := &_PipeEnd[Operation]{inRight, inLeft}

	return left, right
}

type _PipeEnd[Operation any] struct {
	in  <-chan causal.SyncMsg[Operation]
	out chan<- causal.SyncMsg[Operation]
}

type _PipeMsg[Operation any] struct {
	msg causal.SyncMsg[Operation]
	err error
}

func (peer *_PipeEnd[Operation]) Recv(ctx context.Context) (causal.SyncMsg[Operation], error) {
	select {

	case <-ctx.Done():
		return causal.SyncMsg[Operation]{}, ctx.Err()

	case msg, ok := <-peer.in:

		if !ok {
			return causal.SyncMsg[Operation]{}, causal.ErrCommClosed()
		}

		return msg, nil
	}
}

func (peer *_PipeEnd[Operation]) Send(ctx context.Context, msg causal.SyncMsg[Operation]) error {
	select {

	case <-ctx.Done():
		return ctx.Err()

	case peer.out <- msg:
		return nil
	}
}

func (peer *_PipeEnd[Operation]) Close(_ context.Context) error {
	close(peer.out)
	return nil
}
