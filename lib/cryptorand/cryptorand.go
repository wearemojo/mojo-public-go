package cryptorand

import (
	"crypto/rand"
	"encoding/binary"

	mathrand "math/rand/v2"
)

// TODO: could this be removed entirely now that we have math/rand/v2?
// the docs do still say "it should not be used for security-sensitive work"

func New() *mathrand.Rand {
	return mathrand.New(NewSource())
}

type source struct{}

func NewSource() mathrand.Source {
	return source{}
}

func (source) Seed(_ int64) {}

func (source) Uint64() uint64 {
	var data [8]byte
	if _, err := rand.Read(data[:]); err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(data[:])
}
