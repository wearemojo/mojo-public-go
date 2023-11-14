package slicecontains_test

import (
	"testing"

	"github.com/wearemojo/mojo-public-go/lib/slicecontains"
)

func TestEmpty(t *testing.T) {
	var a []string
	var b []string

	if !slicecontains.SameStrings(a, b) {
		t.Error("slices should be equal")
	}

	a = []string{}
	b = []string{}

	if !slicecontains.SameStrings(a, b) {
		t.Error("slices should be equal")
	}
}

func TestSameStringsMatch(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"a", "c"}

	if slicecontains.SameStrings(a, b) {
		t.Error("slices should not be equal")
	}
}

func TestSameStringsSubset(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"a"}

	if slicecontains.SameStrings(a, b) {
		t.Error("slices should not be equal")
	}
}
