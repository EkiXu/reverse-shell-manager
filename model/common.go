package model

import (
	"bytes"
	"sync"
)

type SafeBuffer struct {
	b  bytes.Buffer
	rw sync.RWMutex
}

func (b *SafeBuffer) Read(p []byte) (n int, err error) {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.b.Read(p)
}
func (b *SafeBuffer) Write(p []byte) (n int, err error) {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.b.Write(p)
}
