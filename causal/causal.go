package causal

import (
	"context"
	"iter"
	"slices"
	"unique"
)

type RemotePeer[Operation any] interface {
	Send(ctx context.Context, msg SyncMsg[Operation]) error
	Recv(ctx context.Context) (SyncMsg[Operation], error)
}

type SyncMsg[Operation any] struct {
	DoYouKnow struct {
		IDs      []EventID
		Versions []Version
	}

	IDontKnow []EventID

	LetMeIntroduce struct {
		Ats        []PeerID
		Operations [][]Operation
	}
}

// A Log is a causally-ordered log of events.
type Log[Operation any] struct{}

type Snapshot[Operation any] struct{}

type Transaction[Operation any] struct{}

func (log *Log[Operation]) Sync(ctx context.Context, p RemotePeer[Operation]) error {
	return nil
}

func (log *Log[Operation]) Here() PeerID {
	panic("TODO")
}

func (log *Log[Operation]) Snapshot() Snapshot[Operation] {
	panic("TODO")
}

// Guarantees that snapshots created by local transactions on the latest snapshot come before ones due to remote sync.
func (lob *Log[Operation]) Snapshots() iter.Seq[Snapshot[Operation]] {
	panic("TODO")
}

func (s Snapshot[Operation]) Here() PeerID {
	panic("TODO")
}

func (s Snapshot[Operation]) Version() Version {
	panic("TODO")
}

func (s Snapshot[Operation]) Event(id EventID) (Event[Operation], bool) {
	panic("TODO")
}

func (s Snapshot[Operation]) Changes() iter.Seq[Event[Operation]] {
	panic("TODO")
}

func (s Snapshot[Operation]) Begin() *Transaction[Operation] {
	panic("TODO")
}

func (t *Transaction[Operation]) WIP() Snapshot[Operation] {
	panic("TODO")
}

func (t *Transaction[Operation]) Add(ops ...Operation) error {
	panic("TODO")
}

func (t *Transaction[Operation]) Commit(ctx context.Context) error {
	panic("TODO")
}
func (t *Transaction[Operation]) Abort(ctx context.Context) error {
	panic("TODO")
}

type Event[Operation any] struct {
	id         EventID
	peer       PeerID
	parents    Version
	operations []Operation
}

func (evt Event[Operation]) ID() EventID { return evt.id }

func (evt Event[Operation]) Peer() PeerID { return evt.peer }

func (evt Event[Operation]) Parents() Version { return evt.parents }

// Logically, all operations in an event happen at the same time, at the same peer.
// We can't force random-iteration order.
// We need a consistent order so that the event ID can be a hash, preserving integrity.
func (evt Event[Operation]) Operations() iter.Seq[Operation] {
	return slices.Values(evt.operations)
}

type Version struct{ frontier []EventID }

func (v Version) LatestEvents() iter.Seq[EventID] {
	return slices.Values(v.frontier)
}

type PeerID struct{ h unique.Handle[string] }

type EventID struct{ h unique.Handle[string] }
