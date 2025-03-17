package causal_test

import (
	"testing"

	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theerr"
	"github.com/szabba/assert/v3/assertions/theval"

	"github.com/prelift/replicateddata-go/causal"
	"github.com/prelift/replicateddata-go/internal/catch"
)

func TestSnapshot(t *testing.T) {

	t.Run("Zero", func(t *testing.T) {
		// given
		for name, f := range noninitSnapshotCases() {
			t.Run(name, func(t *testing.T) {

				var zero causal.Snapshot[int]

				// when
				caught := catch.Panic(func() { f(zero) })

				// then
				assert.UsingFmt(t.Errorf).
					That(theval.Equal(caught, "snapshot was not initialized"))
			})
		}
	})

	t.Run("Here", func(t *testing.T) {
		// given
		log, err := causal.NewLog[int]()
		assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

		snapshot := log.Snapshot()

		want := log.Here()

		// when
		got := snapshot.Here()

		// then
		assert.UsingFmt(t.Errorf).
			That(theval.Equal(got, want))
	})

	t.Run("Version", func(t *testing.T) {

		t.Run("EmptyLog", func(t *testing.T) {
			// given
			log, err := causal.NewLog[int]()
			assert.UsingFmt(t.Fatalf).That(theerr.IsNil(err))

			snapshot := log.Snapshot()

			var want causal.Version

			// when
			got := snapshot.Version()

			// then
			assert.UsingFmt(t.Errorf).
				That(theval.DeepEqual(got, want))
		})

	})
}

func noninitSnapshotCases() map[string]func(s causal.Snapshot[int]) {
	return map[string]func(causal.Snapshot[int]){
		"Here":    func(s causal.Snapshot[int]) { s.Here() },
		"Version": func(s causal.Snapshot[int]) { s.Version() },
	}
}
