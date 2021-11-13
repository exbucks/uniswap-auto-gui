package services

import (
	"strconv"

	"github.com/uniswap-auto-gui/utils"
)

func LastPriceFromSwaps(eth utils.Crypto, swaps utils.Swaps) (price float64) {
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

func PriceUpDownFromSwaps(eth utils.Crypto, swaps utils.Swaps, counts int64) (state bool) {
	if eth.Data.Bundles != nil && swaps.Data.Swaps != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		var amount1 float64
		var amount2 float64
		if swaps.Data.Swaps[0].Pair.Token0.Symbol == "WETH" {
			amount1, _ = strconv.ParseFloat(swaps.Data.Swaps[0].Pair.Token1.DerivedETH, 32)
		} else {
			amount1, _ = strconv.ParseFloat(swaps.Data.Swaps[0].Pair.Token0.DerivedETH, 32)
		}
		if swaps.Data.Swaps[1].Pair.Token0.Symbol == "WETH" {
			amount2, _ = strconv.ParseFloat(swaps.Data.Swaps[1].Pair.Token1.DerivedETH, 32)
		} else {
			amount2, _ = strconv.ParseFloat(swaps.Data.Swaps[1].Pair.Token0.DerivedETH, 32)
		}
		price1 := unit * amount1
		price2 := unit * amount2
		state = price1 > price2
	}
	return state
}
