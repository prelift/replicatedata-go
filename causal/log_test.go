package causal_test

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

func TestLog(t *testing.T) {

	t.Run("Zero", func(t *testing.T) {
		// given
		for name, f := range noninitLogCases() {
			t.Run(name, func(t *testing.T) {

				var zero *causal.Log[int]

				// when
				caught := catch.Panic(func() { f(zero) })

				// then
				assert.UsingFmt(t.Errorf).
					That(theval.Equal(caught, "log is nil"))
			})
		}
	})

	t.Run("PtrToZero", func(t *testing.T) {
		// given
		for name, f := range noninitLogCases() {
			t.Run(name, func(t *testing.T) {

				var zero causal.Log[int]

				// when
				caught := catch.Panic(func() { f(&zero) })

				// then
				assert.UsingFmt(t.Errorf).
					That(theval.Equal(caught, "log was not initialized"))
			})
		}
	})

	t.Run("NewLog", func(t *testing.T) {

		t.Run("HasANonEmptyLocalPeerIDByDefault", func(t *testing.T) {
			// given
			log, err := causal.NewLog[int]()

			assert.UsingFmt(t.Fatalf).
				That(theval.NotZero(log)).
				That(theerr.IsNil(err))

			// when
			here := log.Here()

			// then
			assert.UsingFmt(t.Errorf).
				That(theval.NotEqual(here.String(), ""))
		})

		t.Run("FailsGivenAnEmptyPeerID", func(t *testing.T) {
			// given
			opt := causal.WithPeerID[int]("")

			// when
			log, err := causal.NewLog(opt)

			// then
			assert.UsingFmt(t.Fatalf).
				That(theval.Zero(log)).
				That(theerr.Is(err, causal.ErrEmptyPeerID()))
		})

		t.Run("RetainsANonEmptyPeerID", func(t *testing.T) {
			// given
			opt := causal.WithPeerID[int]("a-peer")

			// when
			log, err := causal.NewLog(opt)

			// then
			assert.UsingFmt(t.Fatalf).
				That(theval.NotZero(log)).
				That(theerr.IsNil(err))
		})

	})

	t.Run("AcceptSync", func(t *testing.T) {

		t.Run("FailsWhenThePeerDoesNotProduceTheFirstMessage", func(t *testing.T) {
			// given
			var peer causaltest.ScriptedPeer[int]
			peer.ScriptRecv(causal.SyncMsg[int]{}, causal.ErrGotNoMessage())
			peer.ScriptClose(nil)

			log, err := causal.NewLog[int]()
			assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

			// when
			err = log.AcceptSync(context.Background(), &peer)

			// then
			assert.UsingFmt(t.Errorf).
				That(peer.PlayedOut()).
				That(theerr.Is(err, causal.ErrGotNoMessage()))
		})

		t.Run("EndsWhenInformedOfCommunicationClosureOnFirstMessageReceiveAttempt", func(t *testing.T) {
			// given
			var peer causaltest.ScriptedPeer[int]
			peer.ScriptRecv(causal.SyncMsg[int]{}, causal.ErrCommClosed())
			peer.ScriptClose(nil)

			log, err := causal.NewLog[int]()
			assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

			// when
			err = log.AcceptSync(context.Background(), &peer)

			// then
			assert.UsingFmt(t.Errorf).
				That(peer.PlayedOut()).
				That(theerr.IsNil(err))
		})

		t.Run("FailsWhenSendingTheFirstMessageToThePeerDoes", func(t *testing.T) {
			// given
			var peer causaltest.ScriptedPeer[int]
			peer.ScriptRecv(causal.SyncMsg[int]{}, nil)
			peer.ScriptSend(causal.SyncMsg[int]{}, causal.ErrSendFailed())
			peer.ScriptClose(nil)

			log, err := causal.NewLog[int]()
			assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

			// when
			err = log.AcceptSync(context.Background(), &peer)

			// then
			assert.UsingFmt(t.Errorf).
				That(peer.PlayedOut()).
				That(theerr.Is(err, causal.ErrSendFailed()))
		})

		t.Run("EndsWhenNotifiedOfCommunicationClosureOnFirstSendAttempt", func(t *testing.T) {
			// given
			var peer causaltest.ScriptedPeer[int]
			peer.ScriptRecv(causal.SyncMsg[int]{}, nil)
			peer.ScriptSend(causal.SyncMsg[int]{}, causal.ErrCommClosed())
			peer.ScriptClose(nil)

			log, err := causal.NewLog[int]()
			assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

			// when
			err = log.AcceptSync(context.Background(), &peer)

			// then
			assert.UsingFmt(t.Errorf).
				That(peer.PlayedOut()).
				That(theerr.IsNil(err))
		})

	})

}

func noninitLogCases() map[string]func(log *causal.Log[int]) {
	return map[string]func(*causal.Log[int]){
		"Here":     func(l *causal.Log[int]) { l.Here() },
		"Snapshot": func(l *causal.Log[int]) { l.Snapshot() },
		"AcceptSync": func(l *causal.Log[int]) {
			l.AcceptSync(context.Background(), &causaltest.ScriptedPeer[int]{})
		},
		"OfferSync": func(l *causal.Log[int]) {
			l.OfferSync(context.Background(), &causaltest.ScriptedPeer[int]{})
		},
	}
}
