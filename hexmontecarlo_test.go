package main

import (
	"strconv"
	"testing"

	prng "github.com/ericlagergren/go-prng/xorshift"
)

var testBoardSizes = []int{3, 5, 7, 25, 50}

func BenchmarkNewRandomBernoulliBoard(b *testing.B) {
	rng := new(prng.Shift128Plus)
	rng.Seed()

	for _, size := range testBoardSizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newRandomBernoulliBoard(rng, 101)
			}
		})
	}
}

func BenchmarkBoardGale(b *testing.B) {
	rng := new(prng.Shift128Plus)
	rng.Seed()

	for _, size := range testBoardSizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			brd := newRandomBernoulliBoard(rng, size)
			for i := 0; i < b.N; i++ {
				brd.gale()
			}
		})
	}

}
