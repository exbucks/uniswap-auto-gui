package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"

	"github.com/uniswap-auto-gui/utils"
)

func Notify(title string, message string) {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = gosxnotifier.Default
	note.Push()
}

func printInfo(id string, status string, min float64, max float64, minTarget string, maxTarget string, minTime time.Time, maxTime time.Time, period time.Duration) {
	fmt.Println()
	fmt.Println("Token ID: ", id)
	fmt.Println("Status: ", status)
	fmt.Println("Min price: ", min, minTarget, minTime)
	fmt.Println("Max price: ", max, maxTarget, maxTime)
	fmt.Println("Time Frame: ", period)
}

func findToken(pings <-chan string, id string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		min, max, minTarget, maxTarget, minTime, maxTime := MinAndMax(swaps)
		last := LastPrice(swaps)
		_, _, period := PeriodOfSwaps(swaps)
		gap := PeriodOfGap(swaps)

		if gap < time.Duration(3*time.Hour) {
			if (max-min)/last > 0.1 && period < time.Duration(6*time.Hour) {
				printInfo(id, "Tradable", min, max, minTarget, maxTarget, minTime, maxTime, period)
			} else if (max-min)/last < 0.1 && period > time.Duration(24*time.Hour) {
				printInfo(id, "Stable", min, max, minTarget, maxTarget, minTime, maxTime, period)
			} else {
				fmt.Print(".")
			}
		}
	}
}

func Price(eth utils.Crypto, tokens utils.Tokens) (price float64) {
	if eth.Data.Bundles != nil && tokens.Data.Tokens != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		amount, _ := strconv.ParseFloat(tokens.Data.Tokens[0].DerivedETH, 32)
		price = unit * amount
	}
	return price
}

func PriceFromSwaps(eth utils.Crypto, swaps utils.Swaps) (price float64) {
	if eth.Data.Bundles != nil && swaps.Data.Swaps != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		var amount float64
		if swaps.Data.Swaps[0].Pair.Token0.Symbol == "WETH" {
			amount, _ = strconv.ParseFloat(swaps.Data.Swaps[1].Pair.Token1.DerivedETH, 32)
		} else {
			amount, _ = strconv.ParseFloat(swaps.Data.Swaps[0].Pair.Token0.DerivedETH, 32)
		}
		price = unit * amount
	}
	return price
}

func LastPrice(swaps utils.Swaps) (price float64) {
	item := swaps.Data.Swaps[0]
	price, _ = priceOfSwap(item)
	return price
}

func MinAndMax(swaps utils.Swaps) (
	min float64,
	max float64,
	minTarget string,
	maxTarget string,
	minTime time.Time,
	maxTime time.Time,
) {
	min = 0
	max = 0
	var _min int64
	var _max int64
	for _, item := range swaps.Data.Swaps {
		price, target := priceOfSwap(item)
		minTarget = target
		maxTarget = target
		if min == 0 || max == 0 {
			min = price
			max = price
		}
		if price < min {
			min = price
			_min, _ = strconv.ParseInt(item.Timestamp, 10, 64)
		}
		if price > max {
			max = price
			_max, _ = strconv.ParseInt(item.Timestamp, 10, 64)
		}
	}
	minTime = time.Unix(_min, 0)
	maxTime = time.Unix(_max, 0)
	return min, max, minTarget, maxTarget, minTime, maxTime
}

func PeriodOfSwaps(swaps utils.Swaps) (time.Time, time.Time, time.Duration) {
	first, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
	last, _ := strconv.ParseInt(swaps.Data.Swaps[len(swaps.Data.Swaps)-1].Timestamp, 10, 64)
	tf := time.Unix(first, 0)
	tl := time.Unix(last, 0)
	period := tf.Sub(tl)
	return tl, tf, period
}

func PeriodOfGap(swaps utils.Swaps) time.Duration {
	latest, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
	end := time.Unix(latest, 0)
	now := time.Now()
	period := now.Sub(end)
	return period
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		c := make(chan string, 1)
		go utils.Post(c, "swaps", item.Id)
		findToken(c, item.Id)
	}
}
