package causal

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"unique"
)

// LogOptions are configuration values for [NewLog].
type LogOptions[Operation any] struct {
	funcs []func(log *Log[Operation]) error
}

func (opts LogOptions[Operation]) apply(log *Log[Operation]) error {
	for _, f := range opts.funcs {
		err := f(log)
		if err != nil {
			return err
		}
	}
	return nil
}

// JoinOptions combines many log options into one set.
// In case of conflicts, the latest of the confliction options wins.
func JoinOptions[Operation any](opts ...LogOptions[Operation]) LogOptions[Operation] {
	size := 0
	for _, o := range opts {
		size += len(o.funcs)
	}

	funcs := make([]func(*Log[Operation]) error, 0, size)
	for _, o := range opts {
		funcs = append(funcs, o.funcs...)
	}

	return LogOptions[Operation]{funcs}
}

func WithRandomPeerID[Operation any]() LogOptions[Operation] {
	return opt(func(log *Log[Operation]) error {

		var buf [sha512.Size]byte
		_, _ = rand.Read(buf[:])
		asHex := hex.EncodeToString(buf[:])

		log.here = PeerID{unique.Make(asHex)}

		return nil
	})
}

func WithPeerID[Operation any](id RawPeerID) LogOptions[Operation] {
	return opt(func(log *Log[Operation]) error {
		if id == "" {
			return ErrEmptyPeerID()
		}

		log.here = PeerID{unique.Make(string(id))}
		return nil
	})
}

// opt is a helper function for returning log options.
func opt[Operation any](f func(log *Log[Operation]) error) LogOptions[Operation] {
	opts := LogOptions[Operation]{}
	opts.funcs = append(opts.funcs, f)
	return opts
}
