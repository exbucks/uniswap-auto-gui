package services

import (
	"strconv"
	"time"

	"github.com/uniswap-auto-gui/utils"
)

func SwapsInfo(swaps utils.Swaps) (name string, price float64, change float64, duration float64, alert bool) {
	name = tokenName(swaps)
	price, change = priceChanges(swaps)
	_, _, duration = periodOfSwaps(swaps)
	alert = priceAlert(swaps, 0.1)
	return name, price, change, duration, alert
}

func tokenName(swaps utils.Swaps) (name string) {
	if swaps.Data.Swaps != nil {
		if swaps.Data.Swaps[0].Pair.Token0.Symbol == "WETH" {
			name = swaps.Data.Swaps[0].Pair.Token1.Name
		} else {
			name = swaps.Data.Swaps[0].Pair.Token0.Name
		}

	}
	return name
}

func priceChanges(swaps utils.Swaps) (price float64, change float64) {
	if swaps.Data.Swaps != nil {
		price, _ = priceOfSwap(swaps.Data.Swaps[0])
		last, _ := priceOfSwap(swaps.Data.Swaps[len(swaps.Data.Swaps)-1])
		change = price - last
	}
	return price, change
}

func priceOfSwap(swap utils.Swap) (price float64, target string) {
	amountUSD, _ := strconv.ParseFloat(swap.AmountUSD, 32)
	amountToken, _ := strconv.ParseFloat(swap.Amount0Out, 32)

	if swap.Pair.Token0.Symbol == "WETH" {
		if swap.Amount0In != "0" && swap.Amount1Out != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount1Out, 32)
			target = "BUY"
		} else if swap.Amount0Out != "0" && swap.Amount1In != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount1In, 32)
			target = "SELL"
		}
	} else {
		if swap.Amount0Out != "0" && swap.Amount1In != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
			target = "BUY"
		} else if swap.Amount0In != "0" && swap.Amount1Out != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0In, 32)
			target = "SELL"
		}
	}

	price = amountUSD / amountToken
	return price, target
}

func periodOfSwaps(swaps utils.Swaps) (first time.Time, last time.Time, period float64) {
	var duration float64
	if swaps.Data.Swaps != nil {
		_first, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
		_last, _ := strconv.ParseInt(swaps.Data.Swaps[len(swaps.Data.Swaps)-1].Timestamp, 10, 64)
		first = time.Unix(_first, 0)
		last = time.Unix(_last, 0)
		_period := first.Sub(last)
		duration = _period.Hours()
	}
	return first, last, duration
}

func priceAlert(swaps utils.Swaps, change float64) (state bool) {
	if swaps.Data.Swaps != nil {
		first, _ := priceOfSwap(swaps.Data.Swaps[0])
		second, _ := priceOfSwap(swaps.Data.Swaps[1])
		state = (first-second)/second > change
	}
	return state
}
