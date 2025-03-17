package causal

import (
	"context"
	"iter"
	"slices"
	"unique"
)

type Snapshot[Operation any] struct{}

type Transaction[Operation any] struct{}

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

// A PeerID uniquely identifies a peer.
type PeerID struct{ h unique.Handle[string] }

func (id PeerID) String() string { return id.h.Value() }

// RawPeerID is the raw form of a peer ID.
// Raw IDs are used in sync messages.
func (id PeerID) Raw() RawPeerID { return RawPeerID(id.h.Value()) }

type EventID struct{ h unique.Handle[string] }
