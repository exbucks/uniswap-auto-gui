package services

import (
	"encoding/json"
	"fmt"
	"sync"

	uniswap "github.com/hirokimoto/uniswap-api"
)

func GetAllPairs(target chan int) {
	skip := 0

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(1)
			cc := make(chan string, 1)
			go uniswap.RequestPairs(cc, 1000, 1000*skip)
			msg := <-cc
			var pairs uniswap.Pairs
			json.Unmarshal([]byte(msg), &pairs)
			counts := len(pairs.Data.Pairs)
			fmt.Println(skip, ": ", counts)
			if counts == 0 {
				target <- 111
				return
			}

			// SaveTradePairs(&pairs)
			skip += 1
			target <- skip
			defer wg.Done()
		}
	}()
	target <- 111
}
