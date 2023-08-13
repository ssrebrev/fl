package fl

import (
	"math/rand"
	"testing"
)

const (
	bufSizeAlloc = 128
	minRandBuf   = 66
	randBuf      = 1024
)

func BenchmarkFreeList(b *testing.B) {
	fl := FreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v, _ := fl.Get()
		fl.Put(v)
	}
}

func BenchmarkMpFreeList(b *testing.B) {
	fl := MpFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v, _ := fl.Get()
		fl.Put(v)
	}
}

func BenchmarkMemAlloc(b *testing.B) {
	bufAlloc := bufSizeAlloc
	for i := 0; i < b.N; i++ {
		v := make([]byte, bufAlloc)
		_ = v
	}
}

func BenchmarkSliceFreeList(b *testing.B) {
	fl := SliceFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v := fl.GetOrCreate(bufSizeAlloc)
		fl.Put(v)
	}
}

func BenchmarkSliceMpFreeList(b *testing.B) {
	fl := SliceMpFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v := fl.GetOrCreate(bufSizeAlloc)
		fl.Put(v)
	}
}

func BenchmarkRandMemAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bufAlloc := minRandBuf + rand.Intn(randBuf)
		v := make([]byte, bufAlloc)
		_ = v
	}
}

func BenchmarkRandSliceFreeList(b *testing.B) {
	fl := SliceFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		bufAlloc := minRandBuf + rand.Intn(randBuf)
		v := fl.GetOrCreate(bufAlloc)
		fl.Put(v)
	}
}

func BenchmarkRandSliceMpFreeList(b *testing.B) {
	fl := SliceMpFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		bufAlloc := minRandBuf + rand.Intn(randBuf)
		v := fl.GetOrCreate(bufAlloc)
		fl.Put(v)
	}
}

func BenchmarkLogBase2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nextLogBase2(uint32(i))
	}
}
