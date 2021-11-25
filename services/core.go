package services

import (
	"encoding/json"
	"sync"

	uniswap "github.com/hirokimoto/uniswap-api"
)

func UniswapMarkketPairs(target chan<- []uniswap.Pair, progress chan<- int) {
	var pairs []uniswap.Pair
	skip := 0

	for {
		var wg sync.WaitGroup
		wg.Add(1)

		cc := make(chan string, 1)
		go uniswap.RequestPairs(cc, 1000, 1000*skip)
		msg := <-cc
		var data uniswap.Pairs
		json.Unmarshal([]byte(msg), &data)
		counts := len(data.Data.Pairs)
		pairs = append(pairs, data.Data.Pairs...)
		if counts == 0 {
			progress <- -1
			target <- pairs
			return
		}
		skip += 1
		progress <- skip

		defer wg.Done()
	}
}
