package causal_test

import (
	"testing"

	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theerr"
	"github.com/szabba/assert/v3/assertions/theval"

	"github.com/prelift/replicateddata-go/causal"
	"github.com/prelift/replicateddata-go/internal/catch"
)

func TestLog(t *testing.T) {

	t.Run("Zero", func(t *testing.T) {
		// given
		for name, f := range noninitLogCases() {
			t.Run(name, func(t *testing.T) {

				var zero *causal.Log[int]

				// when
				caught := catch.Panic(func() { f(zero) })

				// then
				assert.UsingFmt(t.Errorf).
					That(theval.Equal(caught, "log is nil"))
			})
		}
	})

	t.Run("PtrToZero", func(t *testing.T) {
		// given
		for name, f := range noninitLogCases() {
			t.Run(name, func(t *testing.T) {

				var zero causal.Log[int]

				// when
				caught := catch.Panic(func() { f(&zero) })

				// then
				assert.UsingFmt(t.Errorf).
					That(theval.Equal(caught, "log is nil"))
			})
		}
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

func noninitLogCases() map[string]func(log *causal.Log[int]) {
	return map[string]func(*causal.Log[int]){
		"Here":     func(l *causal.Log[int]) { l.Here() },
		"Snapshot": func(l *causal.Log[int]) { l.Snapshot() },
	}
}
