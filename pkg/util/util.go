package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

// Fatalf prints the message and exits the program.
func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Println("")
	os.Exit(1)
}

// Fatal prints the message and exits the program.
func Fatal(args ...interface{}) {
	fmt.Fprint(os.Stderr, args...)
	fmt.Println("")
	os.Exit(1)
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandBytes fills the bytes with alphabetic characters randomly
func RandBytes(r *rand.Rand, b []byte) {
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
}

// BufPool is a bytes.Buffer pool
type BufPool struct {
	p *sync.Pool
}

// NewBufPool creates a buffer pool.
func NewBufPool() *BufPool {
	p := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	return &BufPool{
		p: p,
	}
}

// Get gets a buffer.
func (b *BufPool) Get() *bytes.Buffer {
	buf := b.p.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer.
func (b *BufPool) Put(buf *bytes.Buffer) {
	b.p.Put(buf)
}
