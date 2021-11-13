package services

import (
	"strconv"

	"github.com/uniswap-auto-gui/utils"
)

func SwapsInfo(eth utils.Crypto, swaps utils.Swaps) (name string, price float64, change float64) {
	name = tokenName(swaps)
	price, change = priceChanges(eth, swaps)
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

func priceChanges(eth utils.Crypto, swaps utils.Swaps) (price float64, change float64) {
	if eth.Data.Bundles != nil && swaps.Data.Swaps != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		var firstAmount float64
		var lastAmount float64
		last := swaps.Data.Swaps[len(swaps.Data.Swaps)-1]
		if swaps.Data.Swaps[0].Pair.Token0.Symbol == "WETH" {
			firstAmount, _ = strconv.ParseFloat(swaps.Data.Swaps[0].Pair.Token1.DerivedETH, 32)
			lastAmount, _ = strconv.ParseFloat(last.Pair.Token1.DerivedETH, 32)
		} else {
			firstAmount, _ = strconv.ParseFloat(swaps.Data.Swaps[0].Pair.Token0.DerivedETH, 32)
			lastAmount, _ = strconv.ParseFloat(last.Pair.Token0.DerivedETH, 32)
		}
		price = unit * firstAmount
		change = unit * (firstAmount - lastAmount)
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
