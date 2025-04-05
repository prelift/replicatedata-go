package causal

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"

	"github.com/szabba/assert/v3"
)

// ErrEmptyPeerID indicates a peer ID was empty.
func ErrEmptyPeerID() error { return errEmptyPeerID }

// ErrGotNoMessage indicates no message from a peer when one was expected.
func ErrGotNoMessage() error { return errGotNoMessage }

// ErrSendFailed indicates sending a message to a peer failed.
func ErrSendFailed() error { return errSendFailed }

// ErrCommClosed indicates a remote peer decided to close communication.
func ErrCommClosed() error { return errCommClosed }

var (
	errEmptyPeerID  = errors.New("empty peer ID")
	errGotNoMessage = errors.New("got no message")
	errSendFailed   = errors.New("send failed")
	errCommClosed   = errors.New("communication closed")
)

// A Log is a causally-ordered log of events.
type Log[Operation any] struct {
	inited bool
	here   PeerID
}

// NewLog creates a new log.
//
// The zero value is not useful.
// Clients must call NewLog to get a log that's good for something.
func NewLog[Operation any](opts ...LogOptions[Operation]) (_ *Log[Operation], err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("causal.NewLog: %w", err)
		}
	}()

	log := new(Log[Operation])

	err = defaultOpts[Operation]().apply(log)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		err = o.apply(log)
		if err != nil {
			return nil, err
		}
	}

	return log, nil
}

func defaultOpts[Operation any]() LogOptions[Operation] {
	return JoinOptions(
		opt(func(log *Log[Operation]) error {
			log.inited = true
			return nil
		}),
		WithRandomPeerID[Operation](),
	)
}

func (log *Log[Operation]) wasInited() {
	assert.UsingPanic().
		True(log != nil, "log is nil").
		True(log.inited, "log was not initialized")
}

// Here is the local peer ID.
func (log *Log[Operation]) Here() PeerID {
	log.wasInited()
	return log.here
}

// Snapshot takes a snapshot of the current state of the log.
func (log *Log[Operation]) Snapshot() Snapshot[Operation] {
	log.wasInited()
	return Snapshot[Operation]{log: log}
}

// AcceptSync runs a sync with a remote peer that initated communication.
func (log *Log[Operation]) AcceptSync(ctx context.Context, p RemotePeer[Operation]) (err error) {
	log.wasInited()

	defer func() {
		if err != nil {
			err = fmt.Errorf("causal.%T.AcceptSync: %w", log, err)
		}
		cErr := p.Close(ctx)
		if cErr != nil {
			err = errors.Join(err, cErr)
		}
	}()

	for {
		err := log.acceptOne(ctx, p)
		if err != nil {
			return log.tolerateClose(err)
		}

		err = log.sendOne(ctx, p)
		if err != nil {
			return log.tolerateClose(err)
		}
	}
}

// OfferSync initiates communication and runs a sync with a remote peer.
func (log *Log[Operation]) OfferSync(ctx context.Context, p RemotePeer[Operation]) (err error) {
	log.wasInited()

	defer func() {
		if err != nil {
			err = fmt.Errorf("causal.%T.OfferSync: %w", log, err)
		}
		cErr := p.Close(ctx)
		if cErr != nil {
			err = errors.Join(err, cErr)
		}
	}()

	for {
		err = log.sendOne(ctx, p)
		if err != nil {
			return log.tolerateClose(err)
		}

		err := log.acceptOne(ctx, p)
		if err != nil {
			return log.tolerateClose(err)
		}
	}
}

func (log *Log[Operation]) acceptOne(ctx context.Context, p RemotePeer[Operation]) error {
	_, err := p.Recv(ctx)
	if err != nil {
		return err
	}
	slog.WarnContext(ctx, "TODO: process accepted message")
	return nil
}

func (log *Log[Operation]) sendOne(ctx context.Context, p RemotePeer[Operation]) error {
	slog.WarnContext(ctx, "TODO: prepare message to send")
	return p.Send(ctx, SyncMsg[Operation]{})
}

func (*Log[Operation]) tolerateClose(err error) error {
	if errors.Is(err, ErrCommClosed()) {
		return nil
	}
	return err
}

// Guarantees that snapshots created by local transactions on the latest snapshot come before ones due to remote sync.
func (lob *Log[Operation]) Snapshots() iter.Seq[Snapshot[Operation]] {
	panic("TODO")
}
