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

	var num = make([]float64, 0, int(target))
	for i := float64(1); i < target; i++ {
		if math.Mod(i, 3) == 0 || math.Mod(i, 5) == 0 {
			num = append(num, i)
		}
	}
	fmt.Printf("Found values in %s\n", time.Since(processStart))

	sumStart := time.Now()
	var n float64
	for _, val := range num {
		n += val
	}
	fmt.Printf("Calculated sum %d in %s\n", int(n), time.Since(sumStart))
	fmt.Printf("Total run time %s\n", time.Since(processStart))
}
