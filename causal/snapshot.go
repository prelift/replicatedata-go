package causal

import (
	"iter"

	"github.com/szabba/assert/v3"
)

type Snapshot[Operation any] struct {
	log *Log[Operation]
}

func (s Snapshot[Operation]) wasInited() {
	assert.UsingPanic().True(s.log != nil, "snapshot was not initialized")
}

// PeerID is the ID of the peer at which the snapshot was taken.
func (s Snapshot[Operation]) Here() PeerID {
	s.wasInited()
	return s.log.here
}

// Guarantees that snapshots created by local transactions on the latest snapshot come before ones due to remote sync.
func (lob *Log[Operation]) Snapshots() iter.Seq[Snapshot[Operation]] {
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
