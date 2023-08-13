package fl

import "testing"

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

func BenchmarkSliceFreeList(b *testing.B) {
	fl := SliceFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v := fl.GetOrCreate(10)
		fl.Put(v)
	}
}

func BenchmarkSliceMpFreeList(b *testing.B) {
	fl := SliceMpFreeList[byte]{}
	for i := 0; i < b.N; i++ {
		v := fl.GetOrCreate(10)
		fl.Put(v)
	}
}

func BenchmarkLogBase2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nextLogBase2(uint32(i))
	}
}
