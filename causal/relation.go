package causal

import "fmt"

// A Relation describes the relative position of two versions in the (partial) causal order.
type Relation struct{ raw uint8 }

// UnknownOrder means we do not know the causal relation between the two versions in question.
//
// TODO: it's possible we won't need this in practice.
func UnknownOrder() Relation { return Relation{} }

// Before means that the left versions is earliesr than the right one.
func Before() Relation { return Relation{1} }

// Concurrent means that the versions are equal and neither preceedes the other.
func Concurrent() Relation { return Relation{2} }

// Equal means the two versions are equal.
func Equal() Relation { return Relation{3} }

// After means that the left version is later than the right one.
func After() Relation { return Relation{4} }

func (cr Relation) String() string {
	switch cr {
	case UnknownOrder():
		return "UnknownOrder"
	case Before():
		return "Before"
	case Concurrent():
		return "Concurrent"
	case Equal():
		return "Equal"
	case After():
		return "After"
	default:
		err := fmt.Errorf("unsupported causal relation %#v", cr)
		panic(err)
	}
}
