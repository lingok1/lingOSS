package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}

}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length

		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()

		}
		current += next
		length -= next

	}
	return len(p), nil

}

// 文件分片
func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return

	}
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i := range shards {
		//调用http patch方式
		e.writers[i].Write(shards[i])

	}
	e.cache = []byte{}

}
