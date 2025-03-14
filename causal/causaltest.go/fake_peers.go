package causaltest

import (
	"context"
	"fmt"

	"github.com/prelift/replicateddata/causal"
)

// A SendVCR is a remote peer that records all the messages it gets from the local one.
type SendVCR[Operation any] struct {
	// All the messages the local peer sent (or attempted to send).
	Sends []causal.SyncMsg[Operation]

	// An inner remote peer to provide actual behaviour.
	RemotePeer causal.RemotePeer[Operation]
}

func (peer *SendVCR[Operation]) Send(ctx context.Context, msg causal.SyncMsg[Operation]) error {
	err := peer.RemotePeer.Send(ctx, msg)
	peer.Sends = append(peer.Sends, msg)
	return err
}

func (peer *SendVCR[Operation]) Recv(ctx context.Context) (causal.SyncMsg[Operation], error) {
	return peer.RemotePeer.Recv(ctx)
}

// Student is a remote peer that:
//   - starts knowing about it's Known events
//   - accurately reports about what events it already knows
//   - does not attempt to inform the local peer of any events.
type Student[Operation any] struct {
	Known map[causal.EventID]bool

	nextResp, prevResp causal.SyncMsg[Operation]
}

func (peer *Student[Operation]) Send(_ context.Context, msg causal.SyncMsg[Operation]) error {
	peer.init()

	// Store now-known

	if len(msg.LetMeIntroduce.Ats) != len(peer.prevResp.IDontKnow) {
		msg := fmt.Sprintf(
			"trying to introduce %d events when %d are not known",
			len(msg.LetMeIntroduce.Ats),
			len(peer.prevResp.IDontKnow))
		panic(msg)
	}

	for _, id := range peer.prevResp.IDontKnow {
		peer.Known[id] = true
	}

	// Generate next response

	resp := causal.SyncMsg[Operation]{}
	resp.IDontKnow = make([]causal.EventID, 0, len(msg.DoYouKnow.IDs))

	for _, id := range msg.DoYouKnow.IDs {
		if !peer.Known[id] {
			continue
		}
		resp.IDontKnow = append(resp.IDontKnow, id)
	}

	peer.prevResp = peer.nextResp
	peer.nextResp = resp

	return nil
}

func (peer *Student[Operation]) Recv(_ context.Context) (causal.SyncMsg[Operation], error) {

	resp := peer.nextResp
	peer.nextResp = causal.SyncMsg[Operation]{}
	peer.prevResp = resp

	return resp, nil
}

func (peer *Student[Operation]) init() {
	if peer.Known != nil {
		return
	}
	peer.Known = map[causal.EventID]bool{}
}
