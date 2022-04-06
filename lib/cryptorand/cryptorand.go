package cryptorand

import (
	"crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
)

func New() *mathrand.Rand {
	return mathrand.New(NewSource())
}

type source struct{}

func NewSource() mathrand.Source {
	return source{}
}

func (_ source) Seed(_ int64) {}

func (_ source) Int63() int64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}

	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}
