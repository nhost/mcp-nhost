package pkg

import "sync/atomic"

type AtomicBool struct {
	b atomic.Value
}

func NewAtomicBool() *AtomicBool {
	b := &AtomicBool{}
	b.b.Store(false)
	return b
}

func (b *AtomicBool) Store(value bool) {
	b.b.Store(value)
}

func (b *AtomicBool) Load() bool {
	return b.b.Load().(bool)
}
