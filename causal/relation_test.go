package causal_test

import (
	"testing"

	"github.com/szabba/assert/v3"
	"github.com/szabba/assert/v3/assertions/theval"

	"github.com/prelift/replicateddata-go/causal"
)

func TestRelation(t *testing.T) {

	t.Run("ZeroValue", func(t *testing.T) {
		// given

		// when
		rel := causal.UnknownOrder()

		// then
		assert.
			UsingFmt(t.Errorf).
			That(theval.Zero(rel))
	})

	t.Run("String", func(t *testing.T) {
		// given
		cases := map[string]struct {
			Relation causal.Relation
			Want     string
		}{
			"UnknownOrder": {causal.UnknownOrder(), "UnknownOrder"},
			"Before":       {causal.Before(), "Before"},
			"Equal":        {causal.Equal(), "Equal"},
			"After":        {causal.After(), "After"},
			"Concurrent":   {causal.Concurrent(), "Concurrent"},
		}

		for name, tt := range cases {
			t.Run(name, func(t *testing.T) {

				// when
				got := tt.Relation.String()

				// then
				assert.
					UsingFmt(t.Errorf).
					That(theval.Equal(got, tt.Want))

			})
		}
	})

}
