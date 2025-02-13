package ksuid

import (
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for range b.N {
		Generate(b.Context(), "user")
	}
}
