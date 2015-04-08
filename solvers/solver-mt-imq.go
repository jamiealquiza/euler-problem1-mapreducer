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

func isMultiple(n float64) bool {
	if math.Mod(n, 3) == 0 || math.Mod(n, 5) == 0 {
		return true
	}

	return false
}

func moder(wg *sync.WaitGroup, portion int, numbers chan float64, results chan float64) {
	var nums = make([]float64, 0, portion)

	for n := range numbers {
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

	numbersA := make(chan float64, int(target)/2)
	numbersB := make(chan float64, int(target)/2)

	count := 1
	for i := float64(0); i < target; i++ {
		if count%2 == 0 {
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
		if i%2 == 0 {
			go moder(wg, int(target)/goroutines, numbersA, results)
		} else {
			go moder(wg, int(target)/goroutines, numbersB, results)
		}
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
