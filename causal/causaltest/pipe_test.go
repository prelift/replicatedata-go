package causaltest_test

import (
	"context"
	"testing"

	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theerr"
	"github.com/szabba/assert/v3/assertions/theval"

	"github.com/prelift/replicateddata-go/causal"
	"github.com/prelift/replicateddata-go/causal/causaltest"
	"github.com/prelift/replicateddata-go/internal/catch"
)

func TestPipe(t *testing.T) {

	t.Run("RecvCanTimeOutWhenTheOtherPeerCannotProceed", func(t *testing.T) {
		// given

		// Since we don't retain the right end, receiving on the left must either block or time out.
		left, _ := causaltest.Pipe[int]()

		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		// when
		msg, err := left.Recv(ctx)

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Zero(msg)).
			That(theerr.Is(err, context.DeadlineExceeded))
	})

	t.Run("RecvReturnsTheMessageSentAtTheOtherEnd", func(t *testing.T) {
		// given
		left, right := causaltest.Pipe[int]()

		sent := causal.SyncMsg[int]{}
		sent.DoYouKnow = []causal.RawEventID{"ev-1"}

		go right.Send(context.Background(), sent)

		// when
		got, err := left.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.DeepEqual(got, sent)).
			That(theerr.IsNil(err))
	})

	t.Run("RecvFailsWhenTheOtherEndClosedThePipe", func(t *testing.T) {
		// given

		left, right := causaltest.Pipe[int]()

		right.Close(context.Background())

		// when
		got, err := left.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Zero(got)).
			That(theerr.Is(err, causal.ErrCommClosed()))
	})

	t.Run("SendCanTimeOutWhenTheOtherPeerCannotProceed", func(t *testing.T) {
		// given

		// Since we don't retain the right end, receiving on the left must either block or time out.
		left, _ := causaltest.Pipe[int]()

		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		msg := causal.SyncMsg[int]{}
		msg.DoYouKnow = []causal.RawEventID{"ev-1"}

		// when
		err := left.Send(ctx, msg)

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.Is(err, context.DeadlineExceeded))
	})

	t.Run("SendSucceedsWhenTheOtherEndIsTryingToReceive", func(t *testing.T) {
		// given
		left, right := causaltest.Pipe[int]()

		go right.Recv(context.Background())

		msg := causal.SyncMsg[int]{}
		msg.DoYouKnow = []causal.RawEventID{"ev-1"}

		// when
		err := left.Send(context.Background(), msg)

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.IsNil(err))
	})

	t.Run("SendPanicsWhenTheSendingEndIsAlreadyClosed", func(t *testing.T) {
		// given
		left, _ := causaltest.Pipe[int]()

		left.Close(context.Background())

		msg := causal.SyncMsg[int]{}
		msg.DoYouKnow = []causal.RawEventID{"ev-1"}

		// when
		caught := catch.Panic(func() { left.Send(context.Background(), msg) })

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.NotZero(caught))
	})

}
