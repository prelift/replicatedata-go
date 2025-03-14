package causaltest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/prelift/replicateddata-go/causal"
)

func ErrScriptMismatch() error { return errScriptMismatch }

var errScriptMismatch = errors.New("script mismatch")

type ScriptedPeer[Operation any] struct {
	Lines []causal.RemotePeer[Operation]
}

func (peer *ScriptedPeer[Operation]) Recv(ctx context.Context) (causal.SyncMsg[Operation], error) {
	if len(peer.Lines) == 0 {
		return causal.SyncMsg[Operation]{}, net.ErrClosed
	}

	next, rest := peer.Lines[0], peer.Lines[1:]
	peer.Lines = rest

	return next.Recv(ctx)
}

func (peer *ScriptedPeer[Operation]) Send(ctx context.Context, msg causal.SyncMsg[Operation]) error {
	if len(peer.Lines) == 0 {
		return net.ErrClosed
	}

	next, rest := peer.Lines[0], peer.Lines[1:]
	peer.Lines = rest

	return next.Send(ctx, msg)
}

func (peer *ScriptedPeer[Operation]) ScriptSend(want causal.SyncMsg[Operation], err error) {
	peer.Lines = append(peer.Lines, _Send[Operation]{want, err})
}

func (peer *ScriptedPeer[Operation]) ScriptRecv(give causal.SyncMsg[Operation], err error) {
	peer.Lines = append(peer.Lines, _Recv[Operation]{give, err})
}

type _Send[Operation any] struct {
	want causal.SyncMsg[Operation]
	err  error
}

func (peer _Send[Operation]) Recv(_ context.Context) (causal.SyncMsg[Operation], error) {
	err := fmt.Errorf("got Recv, want Send %#v: %w", peer.want, ErrScriptMismatch())
	return causal.SyncMsg[Operation]{}, err
}

func (peer _Send[Operation]) Send(_ context.Context, msg causal.SyncMsg[Operation]) error {
	if !reflect.DeepEqual(msg, peer.want) {
		return fmt.Errorf("got Send %#v, not %#v: %w", msg, peer.want, ErrScriptMismatch())
	}
	return peer.err
}

type _Recv[Operation any] struct {
	give causal.SyncMsg[Operation]
	err  error
}

func (peer _Recv[Operation]) Recv(_ context.Context) (causal.SyncMsg[Operation], error) {
	return peer.give, peer.err
}

func (peer _Recv[Operation]) Send(_ context.Context, msg causal.SyncMsg[Operation]) error {
	return fmt.Errorf("got Send %#v, not Recv: %w", msg, ErrScriptMismatch())
}
