// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	args := []string{
		"./testdata/study.toml",
		"./testdata/out/",
	}
	err := run(args)
	assert.NoError(t, err)
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert.NoError(b, run(nil))
	}
}
