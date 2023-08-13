package fl

import (
	"github.com/ssrebrev/pq"
	"sync/atomic"
)

// MpFreeList is multi producer single consumer free list.
// Put is safe for concurrent access and can be accessed from a multiple goroutines.
// Get should be accessed from a single goroutine.
type MpFreeList[V any] struct {
	live   atomic.Pointer[pq.DynQueue[V]]
	victim atomic.Pointer[pq.DynQueue[V]]
}

// Put is safe for concurrent access and can be accessed from a multiple goroutines.
func (fl *MpFreeList[V]) Put(v V) {
	l := fl.liveItems()
	l.Writer.Enqueue(v)
}

// Get should be accessed from a single goroutine.
func (fl *MpFreeList[V]) Get() (V, bool) {
	l := fl.live.Load()
	if l != nil {
		if _, val, ok := l.Reader.TryDequeue(); ok {
			return val, true
		}
	}

	v := fl.victim.Load()
	if v != nil {
		if _, val, ok := v.Reader.TryDequeue(); ok {
			return val, true
		}
	}

	var zero V
	return zero, false
}

// GetOrCreate should be accessed from a single goroutine.
func (fl *MpFreeList[V]) GetOrCreate(factory func() V) V {
	if v, ok := fl.Get(); ok {
		return v
	}
	return factory()
}

func (fl *MpFreeList[V]) roundOneClean() {
	l := fl.live.Swap(nil)
	fl.victim.Store(l)
}

func (fl *MpFreeList[V]) roundTwoClean() {
	fl.victim.Store(nil)
}

func (fl *MpFreeList[V]) liveItems() *pq.DynQueue[V] {
	l := fl.live.Load()
	if l != nil {
		return l
	}
	q := pq.NewDynQueue[V]()
	fl.live.Store(&q)
	registerCleaner(fl)
	return &q
}
