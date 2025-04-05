package causal

import (
	"context"
)

// A RemotePeer can that can be sent and send back sync messages.
type RemotePeer[Operation any] interface {
	// Send attempts to send a sync message to the remote peer.
	Send(ctx context.Context, msg SyncMsg[Operation]) error

	// Recv attempts to receive a sync message from the remote peer.
	Recv(ctx context.Context) (SyncMsg[Operation], error)

	// Close attempts to inform the remote peer about intentional communication closure.
	Close(ctx context.Context) error
}

// A SyncMsg is a message exchanged between peers as part of log syncing.
//
// A sync message is not correct by construction.
// The zero message is correct at some points in the sync protocol.
type SyncMsg[Operation any] struct {

	// DoYouKnow is a query about which of the specified events the peer already knows.
	//
	// The IDs must be topologically sorted by the happens-before relation.
	// If some of the events are not known to the receiving peer, this cannot be verified immediately.
	DoYouKnow []RawEventID

	// IDontKnow is a response to a DoYouKnow query from the previous message.
	// It contains the subsequence of IDs the sending peer does not know yet.
	// It must preserve the order that was used in DoYouKnow.
	//
	// In a proper message, it is empty when either
	//
	//     - there were no previous messages in the current run of the protocol,
	//     - DoYouKnow in the previous message was empty.
	//
	IDontKnow []RawEventID

	// LetMeIntroduce is a response to the IDontKnow part of the previous message.
	// It contains information about events the receiving peer claimed not to know.
	//
	// Details of an event have the same index as it's ID has in the matching IDontKnow.
	LetMeIntroduce LetMeIntroduce[Operation]
}

// A LetMeIntroduce is a response to an IDontKnow from the previous [SyncMsg].
type LetMeIntroduce[Operation any] struct {
	Peers      []RawPeerID
	Versions   []RawVersion
	Operations [][]Operation
}

// A RawPeerID is a raw, potentially incorrect, form of a peer ID.
// Used in protocol messages.
type RawPeerID string

// A RawEvent ID is a raw, potentially incorrect, form of an event ID.
// Used in sync protocol messages.
type RawEventID string

// A RawVersion is a raw, potentially incorrect, form of a snapshot version.
// Used in sync protocol messages.
type RawVersion []RawEventID
