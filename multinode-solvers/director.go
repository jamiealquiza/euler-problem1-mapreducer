package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	//"math"
	"net/http"
	"net/url"
	"os"
	"time"
	"strconv"
	"encoding/json"
)

var (
	target     float64
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
	flag.Parse()
}

func main() {
	reqMax := strconv.FormatFloat(target, 'f', -1, 64)
	rangeReq := fmt.Sprintf("0:%s", reqMax)

	resp, err := http.PostForm("http://localhost:9000",
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
	fmt.Println(result)
}