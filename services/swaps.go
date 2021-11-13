package services

import (
	"strconv"

	"github.com/uniswap-auto-gui/utils"
)

func SwapsInfo(swaps utils.Swaps) (name string, price float64, change float64) {
	name = tokenName(swaps)
	price, change = priceChanges(swaps)
	return name, price, change
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
		if swap.Amount0In == "0" && swap.Amount1Out == "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
			target = "BUY"
		} else if swap.Amount0Out == "0" && swap.Amount1In == "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount1Out, 32)
			target = "SELL"
		} else if swap.Amount0In != "0" && swap.Amount0Out != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
			target = "BUY"
		}
	} else {

	}

	price = amountUSD / amountToken
	return price, target
}
