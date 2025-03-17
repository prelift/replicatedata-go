package causal

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/szabba/assert/v3"
)

// ErrEmptyPeerID indicates a peer ID was empty.
func ErrEmptyPeerID() error { return errEmptyPeerID }

var errEmptyPeerID = errors.New("empty peer ID")

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

// Guarantees that snapshots created by local transactions on the latest snapshot come before ones due to remote sync.
func (lob *Log[Operation]) Snapshots() iter.Seq[Snapshot[Operation]] {
	panic("TODO")
}

func (log *Log[Operation]) SyncTo(ctx context.Context, p RemotePeer[Operation]) error {
	return nil
}

func (log *Log[Operation]) SyncFrom(ctx context.Context, p RemotePeer[Operation]) error {
	return nil
}
