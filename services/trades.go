package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2/data/binding"
	"github.com/uniswap-auto-gui/utils"
)

func StableTokens(wg *sync.WaitGroup, pairs utils.Pairs, list binding.ExternalStringList) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		c := make(chan string, 1)
		fmt.Println(item.Id)
		go utils.Post(c, "swaps", item.Id)
		stableToken(c, item.Id, list)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs, list binding.ExternalStringList) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		c := make(chan string, 1)
		fmt.Println(item.Id)
		go utils.Post(c, "swaps", item.Id)
		tradableToken(c, item.Id, list)
	}
}

func stableToken(pings <-chan string, id string, list binding.ExternalStringList) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		min, max, _, _, _, _ := minMax(swaps)
		last, _ := priceOfSwap(swaps.Data.Swaps[0])
		_, _, period := periodOfSwaps(swaps)
		howold := howMuchOld(swaps)

		if (max-min)/last < 0.1 && period > 24 && howold < 24 {
			list.Append(id)
		}
	}
}

func tradableToken(pings <-chan string, id string, list binding.ExternalStringList) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		min, max, _, _, _, _ := minMax(swaps)
		last, _ := priceOfSwap(swaps.Data.Swaps[0])
		_, _, period := periodOfSwaps(swaps)
		howold := howMuchOld(swaps)

		if (max-min)/last > 0.1 && period < 6 && howold < 24 {
			list.Append(id)
		}
	}
}
