package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	args := []string{
		"../../testdata/",
		"./study.toml",
	}
	err := run(args)
	assert.NoError(t, err)
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(nil)
	}
}
