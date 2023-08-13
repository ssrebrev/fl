package internal

import (
	_ "unsafe"
)

//go:linkname RegisterGCStart sync.runtime_registerPoolCleanup
func RegisterGCStart(f func())
