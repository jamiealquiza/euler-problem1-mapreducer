package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	target     float64
	goroutines int
	profile    bool
)

func init() {
	flag.Float64Var(&target, "target", 1000, "Eueler Project problem #1 target limit")
	flag.IntVar(&goroutines, "goroutines", 1, "Number of Goroutines")
	flag.BoolVar(&profile, "profile", false, "Run pprof")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type startEnd struct {
	start float64
	end   float64
	incr  float64
	max   float64
}

// startEnd.next() returns the start and end
// range values, incrementing to the next range values each call.
func (se *startEnd) next() (float64, float64) {
	se.start, se.end = se.end, se.end+se.incr
	if se.end > se.max {
		se.end = se.max
		return se.start, se.end
	}
	return se.start, se.end
}

// modSumTask receives a range (startEnd.start to startEnd.end) and finds
// all numbers within this range that are multiples of 3 and 5, then sums and
// returns the values to the results channel.
func modSumTask(wg *sync.WaitGroup, nrange chan []float64, results chan float64) {
	r := <-nrange
	nums := make([]float64, int(r[1]-r[0]))

	for n := r[0]; n < r[1]; n++ {
		if math.Mod(n, 3) == 0 || math.Mod(n, 5) == 0 {
			nums = append(nums, n)
		}
	}

	var sum float64
	for _, val := range nums {
		sum += val
	}

	results <- sum
	defer wg.Done()
}

func main() {
	flag.Parse()
	if profile {
		f, err := os.Create(os.Args[0] + ".prof")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	processStart := time.Now()

	// Set target value and range increment by
	// dividing our target value by the number of workers.
	se := &startEnd{max: target}
	se.incr = math.Floor(se.max/float64(goroutines) + 1.5)

	// Load numbers channel with []float64 pairs of
	// start and end range values.
	numbers := make(chan []float64, goroutines+1)
	for se.end < se.max {
		start, end := se.next()
		r := []float64{start, end}
		numbers <- r
	}

	close(numbers)
	fmt.Printf("Populated %d numbers in %s\n",
		int(se.max), time.Since(processStart))

	// Start n workers Goroutines.
	fmt.Printf("Starting %d Goroutines\n", goroutines)

	results := make(chan float64, goroutines)
	wg := &sync.WaitGroup{}
	wg.Add(goroutines)

	findStart := time.Now()
	for i := 0; i < goroutines; i++ {
		go modSumTask(wg, numbers, results)
	}

	// Wait for modSumTask jobs to crunch numbers and return results.
	wg.Wait()
	fmt.Printf("Found values in %s\n", time.Since(findStart))
	close(results)

	// Sum results.
	sumStart := time.Now()
	var sum float64
	for n := range results {
		sum += n
	}

	fmt.Printf("Calculated sum %d in %s\n", int(sum), time.Since(sumStart))
	fmt.Printf("Total run time %s\n", time.Since(processStart))
}
