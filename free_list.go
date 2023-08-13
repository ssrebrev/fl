// Copyright (c) 2022, Svetoslav Srebrev <srebrev dot svetoslav at gmail dot com>. All rights reserved.
// Use of this source code is governed by a 3-Clause license that can be found in the LICENSE file.

package fl

import (
	"github.com/ssrebrev/pq"
	"runtime"
	"sync/atomic"
)

const (
	finalizerStopped  = 0
	finalizerPending  = 1
	finalizerProgress = 3
)

type twoRoundCleaner interface {
	roundOneClean()
	roundTwoClean()
}

var (
	flQWriter   *pq.Writer[twoRoundCleaner]
	flQReader   *pq.Reader[twoRoundCleaner]
	flQRepeater *pq.ReadRepeater[twoRoundCleaner]

	finalizerState atomic.Int32
)

func init() {
	q := pq.NewQueue[twoRoundCleaner]()
	flQWriter = q.Writer
	flQReader = q.Reader
	flQRepeater = q.Reader.KeepDataAndCreateRepeater()
}

type dummy struct {
	v *uint64
}

func registerCleaner(cleaner twoRoundCleaner) {
	flQWriter.Enqueue(cleaner)

	// 'finalizer approach' specific code
	for {
		state := finalizerState.Load()
		if state == finalizerStopped {
			if finalizerState.CompareAndSwap(finalizerStopped, finalizerPending) {
				d := new(dummy)
				runtime.SetFinalizer(d, dummyFinalizer)
			}
			return
		} else if state == finalizerPending {
			return
		} else { // state == finalizerProgress
			if finalizerState.CompareAndSwap(finalizerProgress, finalizerPending) {
				return
			}
			// state == finalizerStopped; try again
		}
	}
}

func dummyFinalizer(d *dummy) {
	go finalizerProc(d)
}

func finalizerProc(d *dummy) {
	finalizerState.Store(finalizerProgress)
	more := finalizerClean()
	if !more && finalizerState.CompareAndSwap(finalizerProgress, finalizerStopped) {
		return
	}

	// Run dummyFinalizer again. Either we have more items to clean OR state == finalizerPending
	runtime.SetFinalizer(d, dummyFinalizer)
}

func finalizerClean() bool {
	more := false
	for _, v, ok := flQRepeater.NextAvailable(); ok; _, v, ok = flQRepeater.NextAvailable() {
		v.roundTwoClean()
	}

	for _, v, ok := flQReader.TryDequeue(); ok; _, v, ok = flQReader.TryDequeue() {
		v.roundOneClean()
		more = true
	}
	return more
}
