package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	target   float64
	nodeList string
	nodes    []string
)

type Result struct {
	FindTime  time.Duration `json:"findtime"`
	PopTime   time.Duration `json:"poptime"`
	SumTime   time.Duration `json:"sumtime"`
	TotalTime time.Duration `json:"totaltime"`
	Sum       int           `json:"sum"`
}

func init() {
	flag.Float64Var(&target, "target", 1000, "Eueler Project problem #1 target limit")
	flag.StringVar(&nodeList, "nodes", "127.0.0.1:9000", "Comma delimited list of nodes")
	flag.Parse()

	nodes = strings.Split(nodeList, ",")
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

func requester(wg *sync.WaitGroup, n string, numbers chan []float64, results chan *Result) {

	r := <-numbers

	rangeReq := fmt.Sprintf("%d:%d", int(r[0]), int(r[1]))

	fmt.Printf("Sending request for range %s to %s\n", rangeReq, n)

	resp, err := http.PostForm("http://"+n,
		url.Values{"range": []string{rangeReq}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Result := &Result{}
	json.Unmarshal(body, Result)
	results <- Result

	defer wg.Done()
}

func updateGlobalResult(gr, r *Result) {
	gr.FindTime += r.FindTime
	gr.PopTime += r.PopTime
	gr.SumTime += r.SumTime
	gr.TotalTime += r.TotalTime
	gr.Sum += r.Sum
}

func main() {

	processStart := time.Now()

	nNodes := len(nodes)

	// Set target value, range and increment.
	se := &startEnd{max: target}
	se.incr = math.Floor(se.max / float64(nNodes) + .5)

	wg := &sync.WaitGroup{}
	wg.Add(nNodes)

	// Load numbers channel with []float64 pairs of
	// start and end range values.
	numbers := make(chan []float64, nNodes+1)
	for se.end < se.max {
		start, end := se.next()
		r := []float64{start, end}
		numbers <- r
	}

	results := make(chan *Result, nNodes+1)

	for _, node := range nodes {
		go requester(wg, node, numbers, results)
	}

	close(numbers)
	wg.Wait()
	close(results)

	globalResult := &Result{}
	for i := range results {
		updateGlobalResult(globalResult, i)
	}

	fmt.Printf("\nFound answer %d in %s\n\n", 
		globalResult.Sum, time.Since(processStart))
	fmt.Println("Aggregate compute times:")
	fmt.Printf("\tFound multiples in %s\n", globalResult.FindTime)
	fmt.Printf("\tSummed multiples in %s\n", globalResult.SumTime)
}
