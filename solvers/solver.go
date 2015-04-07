package main

import (
	"flag"
	"fmt"
	"math"
	"time"
)

var (
	target float64
)

func init() {
	flag.Float64Var(&target, "target", 1000, "Eueler Project problem #1 target limit")
	flag.Parse()
}

func main() {
	processStart := time.Now()
	var sum float64

	for i := float64(1); i < target; i++ {
		if math.Mod(i, 3) == 0 || math.Mod(i, 5) == 0 {
			sum += i
		}
	}

	fmt.Printf("Found sum %d in %s\n", int(sum), time.Since(processStart))
}