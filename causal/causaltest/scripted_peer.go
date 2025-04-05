package causaltest

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/prelift/replicateddata-go/causal"
)

func ErrScriptMismatch() error { return errScriptMismatch }

var errScriptMismatch = errors.New("script mismatch")

type ScriptedPeer[Operation any] struct {
	lines []causal.RemotePeer[Operation]
}

func (peer *ScriptedPeer[Operation]) PlayedOut() error {
	if len(peer.lines) > 0 {
		callErrs := make([]error, len(peer.lines))
		for i, l := range peer.lines {
			callErrs[i] = fmt.Errorf("pending call %#v", l)
		}
		return errors.Join(callErrs...)
	}
	return nil
}

func (peer *ScriptedPeer[Operation]) Recv(ctx context.Context) (causal.SyncMsg[Operation], error) {
	if len(peer.lines) == 0 {
		err := fmt.Errorf("got Recv in played-out script: %w", ErrScriptMismatch())
		return causal.SyncMsg[Operation]{}, err
	}

	next, rest := peer.lines[0], peer.lines[1:]
	peer.lines = rest

	return next.Recv(ctx)
}

func (peer *ScriptedPeer[Operation]) Send(ctx context.Context, msg causal.SyncMsg[Operation]) error {
	if len(peer.lines) == 0 {
		return fmt.Errorf("got Send %#v in played-out script: %w", msg, ErrScriptMismatch())
	}

	next, rest := peer.lines[0], peer.lines[1:]
	peer.lines = rest

	return next.Send(ctx, msg)
}

func (peer *ScriptedPeer[Operation]) Close(ctx context.Context) error {
	if len(peer.lines) == 0 {
		return fmt.Errorf("got Close in played-out script: %w", ErrScriptMismatch())
	}

	next, rest := peer.lines[0], peer.lines[1:]
	peer.lines = rest

	return next.Close(ctx)
}

func (peer *ScriptedPeer[Operation]) ScriptSend(want causal.SyncMsg[Operation], err error) {
	peer.lines = append(peer.lines, _Send[Operation]{want, err})
}

func (peer *ScriptedPeer[Operation]) ScriptRecv(give causal.SyncMsg[Operation], err error) {
	peer.lines = append(peer.lines, _Recv[Operation]{give, err})
}

func (peer *ScriptedPeer[Operation]) ScriptClose(err error) {
	peer.lines = append(peer.lines, _Close[Operation]{err})
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

func (peer _Send[Operation]) Close(_ context.Context) error {
	return fmt.Errorf("got Close, want Send %#v: %w", peer.want, ErrScriptMismatch())
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

func (peer _Recv[Operation]) Close(_ context.Context) error {
	return fmt.Errorf("got Close, not Recv: %w", ErrScriptMismatch())
}

type _Close[Operation any] struct {
	err error
}

func (peer _Close[Operation]) Recv(_ context.Context) (causal.SyncMsg[Operation], error) {
	err := fmt.Errorf("got Recv, not Close: %w", ErrScriptMismatch())
	return causal.SyncMsg[Operation]{}, err
}

func (peer _Close[Operation]) Send(_ context.Context, msg causal.SyncMsg[Operation]) error {
	return fmt.Errorf("got Send %#v, not Close: %w", msg, ErrScriptMismatch())
}

func (peer _Close[Operation]) Close(_ context.Context) error {
	return peer.err
}
