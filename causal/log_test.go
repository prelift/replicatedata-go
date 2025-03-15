package causal_test

import (
	"testing"

	"github.com/prelift/replicateddata-go/causal"
	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theerr"
	"github.com/szabba/assert/v3/assertions/theval"
)

func TestLog(t *testing.T) {

	t.Run("Zero", func(t *testing.T) {
		// TODO: write
	})

	t.Run("NewLog", func(t *testing.T) {

		t.Run("HasANonEmptyLocalPeerIDByDefault", func(t *testing.T) {
			// given
			log, err := causal.NewLog[int]()

			assert.UsingFmt(t.Fatalf).
				That(theval.NotZero(log)).
				That(theerr.IsNil(err))

			// when
			here := log.Here()

			// then
			assert.UsingFmt(t.Errorf).
				That(theval.NotEqual(here.String(), ""))
		})

		t.Run("FailsGivenAnEmptyPeerID", func(t *testing.T) {
			// given
			opt := causal.WithPeerID[int]("")

			// when
			log, err := causal.NewLog(opt)

			// then
			assert.UsingFmt(t.Fatalf).
				That(theval.Zero(log)).
				That(theerr.Is(err, causal.ErrEmptyPeerID()))
		})

		t.Run("RetainsANonEmptyPeerID", func(t *testing.T) {
			// given
			opt := causal.WithPeerID[int]("a-peer")

			// when
			log, err := causal.NewLog(opt)

			// then
			assert.UsingFmt(t.Fatalf).
				That(theval.NotZero(log)).
				That(theerr.IsNil(err))
		})

	})

}
