package fl

import (
	"github.com/ssrebrev/pq"
	"sync"
	"testing"
)

const (
	bufSize    = 1024
	numWorkers = 32
)

func Benchmark_MpFreeList(b *testing.B) {
	benchMpFL(b, bufSize)
}

func Benchmark_SyncPool(b *testing.B) {
	benchSyncPool(b, bufSize)
}

func Benchmark_NoFreeList(b *testing.B) {
	benchMemAlloc(b, bufSize)
}

func benchMpFL(b *testing.B, bufSizeToAllocate int) {
	fl := SliceMpFreeList[byte]{}

	queues := make([]pq.Queue[[]byte], numWorkers)
	for i := 0; i < numWorkers; i++ {
		queues[i] = pq.NewQueue[[]byte]()
	}

	worker := func(reader *pq.Reader[[]byte], wg *sync.WaitGroup) {
		for {
			buf := reader.Dequeue()
			if buf == nil {
				break
			}
			// simulate buffer access
			//buf[len(buf)-1] = 6
			fl.Put(buf)
		}
		wg.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(queues[i].Reader, wg)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := queues[i%numWorkers].Writer
		buf := fl.GetOrCreate(bufSizeToAllocate)

		// simulate buffer access
		//buf[len(buf)-1] = 6
		writer.Enqueue(buf)
	}

	for i := 0; i < numWorkers; i++ {
		queues[i%numWorkers].Writer.Enqueue(nil)
	}
	wg.Wait()
}

func benchSyncPool(b *testing.B, bufSize int) {
	pool := sync.Pool{New: func() any { return make([]byte, bufSize) }}

	queues := make([]pq.Queue[[]byte], numWorkers)
	for i := 0; i < numWorkers; i++ {
		queues[i] = pq.NewQueue[[]byte]()
	}

	worker := func(reader *pq.Reader[[]byte], wg *sync.WaitGroup) {
		for {
			buf := reader.Dequeue()
			if buf == nil {
				break
			}
			// simulate buffer access
			//buf[len(buf)-1] = 6
			pool.Put(buf)
		}
		wg.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(queues[i].Reader, wg)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := queues[i%numWorkers].Writer
		buf := pool.Get().([]byte)

		// simulate buffer access
		//buf[len(buf)-1] = 6
		writer.Enqueue(buf)
	}

	for i := 0; i < numWorkers; i++ {
		queues[i%numWorkers].Writer.Enqueue(nil)
	}
	wg.Wait()
}

func benchMemAlloc(b *testing.B, bufSize int) {
	queues := make([]pq.Queue[[]byte], numWorkers)
	for i := 0; i < numWorkers; i++ {
		queues[i] = pq.NewQueue[[]byte]()
	}

	worker := func(reader *pq.Reader[[]byte], wg *sync.WaitGroup) {
		for {
			buf := reader.Dequeue()
			if buf == nil {
				break
			}
			// simulate buffer access
			//buf[len(buf)-1] = 6
		}
		wg.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(queues[i].Reader, wg)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := queues[i%numWorkers].Writer
		buf := make([]byte, bufSize)

		// simulate buffer access
		//buf[len(buf)-1] = 6
		writer.Enqueue(buf)
	}

	for i := 0; i < numWorkers; i++ {
		queues[i%numWorkers].Writer.Enqueue(nil)
	}
	wg.Wait()
}
