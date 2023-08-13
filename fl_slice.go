package fl

import (
	"math"
	"math/bits"
)

const (
	maxBucketLen   = math.MaxInt32
	maxBucketItems = 32
)

type SliceFreeList[V any] struct {
	slots [maxBucketItems]FreeList[[]V]
}

func (fl *SliceFreeList[V]) Put(v []V) {
	capacity := cap(v)
	if capacity < 1 || capacity > maxBucketLen {
		return
	}
	idx := prevLogBase2(uint32(capacity))
	fl.slots[idx].Put(v)
}

func (fl *SliceFreeList[V]) Get(length int) []V {
	if length < 1 || length > maxBucketLen {
		return make([]V, length)
	}

	idx := nextLogBase2(uint32(length))
	if v, _ := fl.slots[idx].Get(); v != nil {
		return v[:length]
	}
	return nil
}

func (fl *SliceFreeList[V]) GetOrCreate(length int) []V {
	if length < 1 || length > maxBucketLen {
		return make([]V, length)
	}

	idx := nextLogBase2(uint32(length))
	if v, _ := fl.slots[idx].Get(); v != nil {
		return v[:length]
	}
	return make([]V, 1<<idx)[:length]
}

type SliceMpFreeList[V any] struct {
	slots [maxBucketItems]MpFreeList[[]V]
}

func (fl *SliceMpFreeList[V]) Put(v []V) {
	capacity := cap(v)
	if capacity < 1 || capacity > maxBucketLen {
		return
	}
	idx := prevLogBase2(uint32(capacity))
	fl.slots[idx].Put(v)
}

func (fl *SliceMpFreeList[V]) Get(length int) []V {
	if length < 1 || length > maxBucketLen {
		return make([]V, length)
	}

	idx := nextLogBase2(uint32(length))
	if v, _ := fl.slots[idx].Get(); v != nil {
		return v[:length]
	}
	return nil
}

func (fl *SliceMpFreeList[V]) GetOrCreate(length int) []V {
	if length < 1 || length > maxBucketLen {
		return make([]V, length)
	}

	idx := nextLogBase2(uint32(length))
	if v, _ := fl.slots[idx].Get(); v != nil {
		return v[:length]
	}
	return make([]V, 1<<idx)[:length]
}

func nextLogBase2(v uint32) int {
	return bits.Len32(v - 1)
}

func prevLogBase2(v uint32) int {
	next := nextLogBase2(v)
	if v == (1 << next) {
		return next
	}
	return next - 1
}
