package causaltest_test

import (
	"context"
	"net"
	"testing"

	"github.com/prelift/replicateddata/causal"
	"github.com/prelift/replicateddata/causal/causaltest.go"
	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theerr"
	"github.com/szabba/assert/v3/assertions/theval"
)

func TestScriptedPeer(t *testing.T) {

	t.Run("RecvReportsNetClosedWithoutScript", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		// when
		msg, err := peer.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Zero(msg)).
			That(theerr.Is(err, net.ErrClosed))
	})

	t.Run("RecvReportsScriptedMessage", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		want := causal.SyncMsg[int]{}
		want.DoYouKnow.IDs = make([]causal.EventID, 1)
		want.DoYouKnow.Versions = make([]causal.Version, 1)

		peer.ScriptSend(want, nil)

		// when
		msg, err := peer.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.DeepEqual(msg, want)).
			That(theerr.IsNil(err))
	})

	t.Run("RecvReportsScriptedError", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		zero := causal.SyncMsg[int]{}
		want := net.ErrWriteToConnected

		peer.ScriptRecv(zero, want)

		// when
		msg, err := peer.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Zero(msg)).
			That(theerr.Is(err, want))
	})

	t.Run("RecvReportsScriptMismatchWhenASendWasScripted", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		expect := causal.SyncMsg[int]{}
		report := net.ErrWriteToConnected

		peer.ScriptSend(expect, report)

		// when
		msg, err := peer.Recv(context.Background())

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Zero(msg)).
			That(theerr.Is(err, causaltest.ErrScriptMismatch()))
	})

	t.Run("SendReportsNetClosedWithoutScript", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]
		msg := causal.SyncMsg[int]{}

		// when
		err := peer.Send(context.Background(), msg)

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.Is(err, net.ErrClosed))
	})

	t.Run("SendReportsScriptMismatchWhenTheSentMessageDoesNotMatch", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		expect := causal.SyncMsg[int]{}
		expect.IDontKnow = make([]causal.EventID, 1)

		peer.ScriptSend(expect, nil)

		msg := causal.SyncMsg[int]{}
		msg.IDontKnow = make([]causal.EventID, 2)

		// when
		err := peer.Send(context.Background(), msg)

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.Is(err, causaltest.ErrScriptMismatch()))
	})

	t.Run("SendReportsScriptedErrorWhenMessageMatchesScript", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		expect := causal.SyncMsg[int]{}
		expect.IDontKnow = make([]causal.EventID, 1)

		peer.ScriptSend(expect, net.ErrWriteToConnected)

		// when
		err := peer.Send(context.Background(), expect)

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.Is(err, net.ErrWriteToConnected))
	})

	t.Run("SendReportsScriptMismatchWhenARecvWasScripted", func(t *testing.T) {
		// given
		var peer causaltest.ScriptedPeer[int]

		give := causal.SyncMsg[int]{}
		give.IDontKnow = make([]causal.EventID, 1)

		peer.ScriptRecv(give, nil)

		// when
		err := peer.Send(context.Background(), causal.SyncMsg[int]{})

		// then
		assert.UsingFmt(t.Errorf).
			That(theerr.Is(err, causaltest.ErrScriptMismatch()))
	})

}
