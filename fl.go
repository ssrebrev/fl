package fl

import "sync/atomic"

// FreeList provides object free list for a goroutine. It is not safe for concurrent access.
type FreeList[V any] struct {
	live   atomic.Pointer[flItems[V]]
	victim atomic.Pointer[flItems[V]]
}

func (fl *FreeList[V]) Put(v V) {
	l := fl.liveItems()
	l.items = append(l.items, v)
}

func (fl *FreeList[V]) Get() (V, bool) {
	l := fl.live.Load()
	if l != nil {
		if val, ok := l.tryGet(); ok {
			return val, true
		}
	}

	v := fl.victim.Load()
	if v != nil {
		if val, ok := v.tryGet(); ok {
			return val, true
		}
	}

	var zero V
	return zero, false
}

func (fl *FreeList[V]) GetOrCreate(factory func() V) V {
	if v, ok := fl.Get(); ok {
		return v
	}
	return factory()
}

type flItems[V any] struct {
	items []V
}

func (fl *FreeList[V]) roundOneClean() {
	l := fl.live.Swap(nil)
	fl.victim.Store(l)
}

func (fl *FreeList[V]) roundTwoClean() {
	fl.victim.Store(nil)
}

func (fl *FreeList[V]) liveItems() *flItems[V] {
	l := fl.live.Load()
	if l != nil {
		return l
	}
	l = &flItems[V]{}
	fl.live.Store(l)
	registerCleaner(fl)
	return l
}

func (fli *flItems[V]) tryGet() (V, bool) {
	var zero V

	length := len(fli.items)
	if length > 0 {
		v := fli.items[length-1]
		fli.items[length-1] = zero
		fli.items = fli.items[:length-1]
		return v, true
	}

	return zero, false
}
