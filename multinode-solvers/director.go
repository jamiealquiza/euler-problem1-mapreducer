package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	"strings"
	"encoding/json"
	"sync"
)

var (
	target     float64
	nodeList	string
	nodes []string
)

type response struct {
	FindTime time.Duration `json:"findtime"`
	PopTime time.Duration `json:"poptime"`
	SumTime time.Duration `json:"sumtime"`
	TotalTime time.Duration `json:"totaltime"`
	Sum int `json:"sum"`
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

func requester(wg *sync.WaitGroup, n string, numbers chan[]float64, results chan *response) {

	r := <- numbers

	rangeReq := fmt.Sprintf("%f:%f", r[0], r[1])

	resp, err := http.PostForm("http://"+n,
		url.Values{"range": []string{rangeReq}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := &response{}
	json.Unmarshal(body, result)
	results <- result

	defer wg.Done()
}

func main() {
	nNodes := len(nodes)

	// Set target value, range and increment.
	se := &startEnd{max: target}
	se.incr = se.max / float64(nNodes)

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

	results := make(chan *response, nNodes)

	for _, node := range nodes {
		go requester(wg, node, numbers, results)
	}
	close(numbers)

	wg.Wait()

	for i := range results {
		fmt.Println(i)
	}
}