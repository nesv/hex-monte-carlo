# hex-monte-carlo

A Monte Carlo simulator, using Gale's Algorithm to determine the winner of a
Hex game, based on the difference of the stone played on a single cell.

There are two implementations in this repository. The first os written in
Python, by [pgadey](https://github.com/pgadey). The second is written in Go,
by [myself](https://github.com/nesv).

## Building the Go version

There is a `GNUmakefile` provided in this repository to ease the building of
the Go version.

> **NOTE**
>
> If this is the first time you are going to be building the Go version of the
> simulator, you must run `make deps` before running `make`.

Simply running

	$ make

will build the `hmc` binary. You can run

	$ ./hmc -h

to see all of the supported flags.
