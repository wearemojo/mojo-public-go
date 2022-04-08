package cryptorand

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
)

func New() *mathrand.Rand {
	return mathrand.New(NewSource()) // nolint:gosec
}

type source struct{}

func NewSource() mathrand.Source {
	return source{}
}

func (source) Seed(_ int64) {}

func (source) Int63() int64 {
	var data [8]byte
	if _, err := rand.Read(data[:]); err != nil {
		panic(err)
	}

	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(data[:]) & (1<<63 - 1))
}
