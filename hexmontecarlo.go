package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	prng "github.com/ericlagergren/go-prng/xorshift"
)

type trialParams struct {
	size, x, y int
}

type trial struct {
	b1, b2 board
}

func newTrial(b board, x, y int) *trial {
	t := trial{
		b1: b,
		b2: b.cp(),
	}

	t.b1[y+1][x+1] = 1
	t.b2[y+1][x+1] = 2

	return &t
}

func (t *trial) run() bool {
	v, err := t.b1.gale()
	if err != nil {
		panic(fmt.Sprintf("error:board=1: %v\n", err))
	}

	w, err := t.b2.gale()
	if err != nil {
		panic(fmt.Sprintf("error:board=2: %v\n", err))
	}

	return v != w
}

type result struct {
	x, y int
	inc  bool
}

func showProgress(ch <-chan struct{}, numTrials, totalRuns int) {
	start := time.Now()
	count := 0
	for _ = range ch {
		count++
		if count%numTrials == 0 {
			elapsed := time.Since(start)
			remaining := time.Duration((float64(elapsed) / float64(count)) * float64(totalRuns-count))
			percentDone := (float64(count) / float64(totalRuns)) * 100

			fmt.Fprintf(os.Stderr, "\rProgress: %d/%d (%0.0f%%) %20s", count, totalRuns, percentDone, remaining)
		}
	}

	fmt.Fprintf(os.Stderr, "\rProgress: 100%% %50s\nSimulation finished in %s; generating output\n", "", time.Since(start))
}

func main() {
	var (
		numTrials     int
		halfBoardSize int
	)
	flag.IntVar(&numTrials, "n", 1000, "Numer of trials to run")
	flag.IntVar(&halfBoardSize, "s", 3, "Half the size of the board")
	flag.Parse()

	boardSize := (2 * halfBoardSize) + 1
	fmt.Fprintf(os.Stderr, "Simulating %dx%d with N=%d trials\n", boardSize, boardSize, numTrials)

	totalRuns := numTrials * boardSize * boardSize
	var wg sync.WaitGroup
	wg.Add(totalRuns)

	// Start a goroutine for updating the progress meter.
	progress := make(chan struct{})
	go showProgress(progress, numTrials, totalRuns)

	// Start a goroutine for updating the critical count table, based upon
	// the result of each trial.
	criticalCount := make([][]int, boardSize)
	for i := 0; i < boardSize; i++ {
		criticalCount[i] = make([]int, boardSize)
	}
	results := make(chan result)
	go func() {
		for r := range results {
			if r.inc {
				criticalCount[r.y][r.x] += 1
			}
		}
	}()

	// Start goroutines that will be responsible for cranking through
	// each trial.
	trials := make(chan trialParams, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			rng := new(prng.Shift128Plus)
			rng.Seed()
			for params := range trials {
				b := newRandomBernoulliBoard(rng, params.size)
				t := newTrial(b, params.x, params.y)
				results <- result{
					inc: t.run(),
					x:   params.x,
					y:   params.y,
				}
				progress <- struct{}{}
				wg.Done()
			}
		}()
	}

	// Generate a trial for each X/Y coordinate.
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			for i := 1; i <= numTrials; i++ {
				trials <- trialParams{size: boardSize, x: x, y: y}
			}
		}
	}

	wg.Wait()
	close(trials)
	close(results)
	close(progress)

	// Display the critical count.
	printBoard(os.Stdout, criticalCount)
}

func printBoard(w io.Writer, b board) {
	tw := tabwriter.NewWriter(w, 0, 80, 1, ' ', 0)
	defer tw.Flush()

	for y := 0; y < len(b); y++ {
		row := make([]string, len(b[y]))
		for x := 0; x < len(b[y]); x++ {
			row[x] = strconv.Itoa(b[y][x])
		}
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
}

type board [][]int

// Produces a filled board, by tossing a "fair" (1,2] coin for each cell.
//
// The board returned by this function also marks the edges of the board with
// stones denoting who needs to connect (for use in Gale's Algorithm).
func newRandomBernoulliBoard(rng *prng.Shift128Plus, size int) board {
	b := make([][]int, size+2)

	// First row of stones for Gale's Algorithm.
	b[0] = make([]int, size+2)
	b[0][0] = 2
	for i := 1; i < size+2; i++ {
		b[0][i] = 1
	}

	for y := 1; y < size+1; y++ {
		b[y] = make([]int, size+2)
		b[y][0] = 2

		for x := 1; x < size+1; x++ {
			b[y][x] = int((rng.Next() % 2) + 1)
		}

		b[y][size+1] = 2
	}

	// Last row of stones for Gale's Algorithm.
	b[size+1] = make([]int, size+2)
	for i := 0; i < size+1; i++ {
		b[size+1][i] = 1
	}
	b[size+1][size+1] = 2

	/*
		fmt.Fprintln(os.Stdout, "-----")
		printBoard(os.Stdout, b)
		fmt.Fprintln(os.Stdout, "-----")
	*/

	return b
}

func (b board) equal(other board) bool {
	if len(other) != len(b) {
		return false
	} else if len(other) >= 1 && len(b) >= 1 && len(other[0]) != len(b[0]) {
		return false
	}
	for y := 0; y < len(b); y++ {
		for x := 0; x < len(b[y]); x++ {
			if other[y][x] != b[y][x] {
				return false
			}
		}
	}
	return true
}

func (b board) cp() board {
	bb := make([][]int, len(b))
	for y, row := range b {
		bb[y] = make([]int, len(row))
		copy(bb[y], b[y])
	}
	return bb
}

// Given a filled board, return the winner.
func (b board) gale() (int, error) {
	v := [][]int{
		{-1, 1},
		{0, 0},
		{0, 1},
		{1, 0},
	}

	n := len(b[0]) - 2
	for (-1 <= v[3][0] && v[3][0] <= n+1) && (0 <= v[3][1] && v[3][1] <= n+1) {
		w := [][]int{
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
		}
		if b[v[3][0]][v[3][1]] == b[v[2][0]][v[2][1]] {
			// Go left.
			b.galeCopy(w[0], v[2])
			b.galeCopy(w[1], v[1])
			b.galeCopy(w[2], v[3])
			vec := []int{
				v[1][0] + (v[1][0] - v[0][0]),
				v[1][1] + (v[1][1] - v[0][1]),
			}
			b.galeCopy(w[3], vec)
			w[3][0], w[3][1] = v[1][0]+(v[1][0]-v[0][0]), v[1][1]+(v[1][1]-v[0][1])
		} else if b[v[3][0]][v[3][1]] == b[v[1][0]][v[1][1]] {
			// Go right.
			b.galeCopy(w[0], v[1])
			b.galeCopy(w[1], v[3])
			b.galeCopy(w[2], v[2])
			vec := []int{
				v[2][0] + (v[2][0] - v[0][0]),
				v[2][1] + (v[2][1] - v[0][1]),
			}
			b.galeCopy(w[3], vec)
		}
		b.galeCopy(v[0], w[0])
		b.galeCopy(v[1], w[1])
		b.galeCopy(v[2], w[2])
		b.galeCopy(v[3], w[3])
	}

	if v[3][1] == -1 {
		return 1, nil
	} else if v[3][0] == 0 {
		return 2, nil
	}

	return 0, errors.New("broken algorithm: not [* -1] or [0 *]")
}

func (b board) galeCopy(dst, src []int) {
	if len(dst) != len(src) {
		panic("galeCopy: slices are of different lengths")
	}
	for i := 0; i < len(dst); i++ {
		dst[i] = src[i]
	}
	return
}
