package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	err := run(nil)
	assert.NoError(t, err)
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(nil)
	}
}
