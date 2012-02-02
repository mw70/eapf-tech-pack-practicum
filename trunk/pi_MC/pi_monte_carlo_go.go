package main

import (
	"flag"
	"fmt"
	"runtime"
)

var (
	nThrow = flag.Int("n", 1e6, "number of throws")
	nCPU   = flag.Int("cpu", 1, "number of CPUs to use")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*nCPU) // Set number of OS threads to use.

	// Split all throws to that number of parallel tasks.
	// Constant factor of 64 is to ensure sufficient dynamic load balancing.
	nTasks := *nCPU * 64
	// Channel to collect partial results.
	parts := make(chan int, nTasks)
	for i := 0; i < nTasks; i++ {
		// Kick off a parallel task.
		go func(me int) {
			hits := 0
			// Create task-local PRNG to avoid synchronization issues.
			r := MakeLeapFrogRand(nTasks, me)
			// Calculate number of throws for this task.
			n := *nThrow / nTasks
			if me < (*nThrow % nTasks) {
				n++
			}
			// Do the throws.
			for i := 0; i < n; i++ {
				x := r.Float64()
				y := r.Float64()
				if x*x+y*y < 1 {
					hits++
				}
			}
			parts <- hits // Send the result back.
		}(i)
	}

	// Aggregate partial results.
	hits := 0
	for i := 0; i < nTasks; i++ {
		hits += <-parts
	}
	pi := 4 * float64(hits) / float64(*nThrow)
	fmt.Printf("PI = %g\n", pi)
}

// LeapFrogRand ensures that all parallel tasks
// receive non-overlapping regions of random space.
// That is to ensure absence of over sampled regions.
type LeapFrogRand struct {
	mult, next uint64
	pad        [64]uint8
}

const (
	LeapFrogMult = 764261123
	LeapFrogPmod = 2147483647
)

func MakeLeapFrogRand(total, me int) LeapFrogRand {
	var r LeapFrogRand
	mult := uint64(LeapFrogMult)
	seed := uint64(1)
	for i := 0; i < total; i++ {
		mult = (mult * LeapFrogMult) % LeapFrogPmod
		seed = (seed * LeapFrogMult) % LeapFrogPmod
		if i == me {
			r.next = seed
		}
	}
	r.mult = mult
	return r
}

func (r *LeapFrogRand) Float64() float64 {
	r.next = (r.next * r.mult) % LeapFrogPmod
	return float64(r.next) / float64(LeapFrogPmod)
}
