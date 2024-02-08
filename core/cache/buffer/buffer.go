package buffer

import (
	"bytes"
	"sync"
)

type SyncWriter struct {
	wr *bytes.Buffer
	m  sync.Mutex
}

func NewWriterSize(n int) *SyncWriter {
	b := make([]byte, n)
	return &SyncWriter{wr: bytes.NewBuffer(b)}
}

func (sw *SyncWriter) Reset() {
	sw.wr.Reset()
}

func (sw *SyncWriter) Write(data []byte) (n int, err error) {
	sw.m.Lock()
	n, err = sw.wr.Write(data)
	sw.m.Unlock()
	return
}

func (sw *SyncWriter) String() string {
	sw.m.Lock()
	defer sw.m.Unlock()
	return sw.wr.String()
}
