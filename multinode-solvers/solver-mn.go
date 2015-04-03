package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	goroutines int
	listen     string
)

func init() {
	flag.IntVar(&goroutines, "goroutines", 1, "Number of Goroutines")
	flag.StringVar(&listen, "listen", "127.0.0.1:9000", "Listen addr:port")
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Received")

	r.ParseForm()

	reqRange := strings.Split(r.Form["range"][0], ":")

	// Set target value and range increment by
	// dividing our target value by the number of workers.
	se := &startEnd{}
	se.end, _ = strconv.ParseFloat(reqRange[0], 64)
	se.max, _ = strconv.ParseFloat(reqRange[1], 64)
	se.incr = math.Floor(float64((int(se.max - se.end)) / goroutines) + .5)
	rangeSize := int(se.max - se.end)

	result := make(map[string]interface{})

	processStart := time.Now()

	// Load numbers channel with []float64 pairs of
	// start and end range values.
	numbers := make(chan []float64, goroutines+1)
	for se.end < se.max {
		start, end := se.next()
		r := []float64{start, end}
		numbers <- r
	}

	result["poptime"] = time.Since(processStart)
	close(numbers)
	fmt.Printf("Populated %d numbers in %s\n",
		rangeSize, result["poptime"])

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
	result["findtime"] = time.Since(findStart)
	fmt.Printf("Found values in %s\n", result["findtime"])
	close(results)

	// Sum results.
	sumStart := time.Now()
	var sum float64
	for n := range results {
		sum += n
	}

	result["sumtime"] = time.Since(sumStart)
	result["totaltime"] = time.Since(processStart)
	result["sum"] = int(sum)
	fmt.Printf("Calculated sum %d in %s\n", result["sum"], result["sumtime"])
	fmt.Printf("Total run time %s\n", result["totaltime"])

	response, _ := json.Marshal(result)
	fmt.Fprintf(w, string(response))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(listen, nil)
}
