package main

import (
	"testing"

	prng "github.com/ericlagergren/go-prng/xorshift"
)

func BenchmarkNewRandomBernoulliBoard(b *testing.B) {
	rng := new(prng.Shift128Plus)
	rng.Seed()

	for i := 0; i < b.N; i++ {
		newRandomBernoulliBoard(rng, 101)
	}
}

func BenchmarkBoardGale(b *testing.B) {
	rng := new(prng.Shift128Plus)
	rng.Seed()

	brd := newRandomBernoulliBoard(rng, 101)

	for i := 0; i < b.N; i++ {
		brd.gale()
	}
}
