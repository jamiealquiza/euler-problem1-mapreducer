package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
	"runtime/pprof"
	"os"
)

var (
	target     float64
	goroutines int
)

func init() {
	flag.Float64Var(&target, "target", 1000, "Eueler Project problem #1 target limit")
	flag.IntVar(&goroutines, "goroutines", 1, "Number of Goroutines")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func isMultiple(n float64) bool {
	if math.Mod(n, 3) == 0 || math.Mod(n, 5) == 0 {
			return true
	}

	return false
}

func moder(wg *sync.WaitGroup, portion int, a, b chan float64, results chan float64) {
	var nums = make([]float64, 0, portion)

	
	for {
		select {
		case n := <- a:
			if isMultiple(n) {
				nums = append(nums, n)
			}
		case n := <- b:
			if isMultiple(n) {
				nums = append(nums, n)
			}
		}
	}
	
	/*count := 1
	for i := float64(0); i < target/2; i++ {
		if count % 2 == 0 {
			numbersA <- i
		} else {
			numbersB <- i
		}
		count++
	}*/

	var sum float64
	for _, val := range nums {
		sum += val
	}

	results <- sum
	defer wg.Done()
}

func main() {
	f, _ := os.Create("out")
	pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

	processStart := time.Now()

	numbersA := make(chan float64, int(target))
	numbersB := make(chan float64, int(target))

	count := 1
	for i := float64(0); i < target/2; i++ {
		if count % 2 == 0 {
			numbersA <- i
		} else {
			numbersB <- i
		}
		count++
	}

	close(numbersA)
	close(numbersB)
	
	fmt.Printf("Populated %d numbers in %s\n",
		int(target), time.Since(processStart))

	fmt.Printf("\tQueue A length: %d\n\tQueue B length: %d\n", 
		len(numbersA), len(numbersB))

	fmt.Printf("Starting %d Goroutines\n", goroutines)

	results := make(chan float64, int(target))
	wg := &sync.WaitGroup{}
	wg.Add(goroutines)

	findStart := time.Now()
	for i := 0; i < goroutines; i++ {
		go moder(wg, int(target)/goroutines, numbersA, numbersB, results)
	}

	wg.Wait()
	fmt.Printf("Found values in %s\n", time.Since(findStart))
	close(results)

	sumStart := time.Now()
	var sum float64
	for n := range results {
		sum += n
	}

	fmt.Printf("Calculated sum %d in %s\n", int(sum), time.Since(sumStart))
	fmt.Printf("Total run time %s\n", time.Since(processStart))
	
}
