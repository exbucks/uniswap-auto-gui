package pages

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	uniswap "github.com/hirokimoto/uniswap-api"
	unitrade "github.com/hirokimoto/uniswap-api/swap"
	unitrades "github.com/hirokimoto/uniswap-api/swaps"
	"github.com/leekchan/accounting"
	"github.com/uniswap-auto-gui/services"
)

type Trade struct {
	pair  uniswap.Pair
	swaps uniswap.Swaps
}

func tradesScreen(_ fyne.Window) fyne.CanvasObject {
	money := accounting.Accounting{Symbol: "$", Precision: 6}

	pairsList := binding.BindStringList(&[]string{})
	trades := map[string]Trade{}

	infProgress := widget.NewProgressBarInfinite()
	// command := make(chan string)
	infProgress.Stop()
	// command <- "Pause", "Stop"

	find := widget.NewButton("Find Trading Pairs", func() {
		infProgress.Start()
	})

	list := widget.NewListWithData(pairsList,
		func() fyne.CanvasObject {
			leftPane := container.NewHBox(widget.NewHyperlink("DEX", parseURL("https://fyne.io/")), widget.NewLabel("token"), widget.NewLabel("price"), widget.NewLabel("change"), widget.NewLabel("duration"))
			return container.NewBorder(nil, nil, leftPane, widget.NewButton("+", nil))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			lc := obj.(*fyne.Container).Objects[0].(*fyne.Container)

			dex := lc.Objects[0].(*widget.Hyperlink)

			f := item.(binding.String)
			pair, _ := f.Get()

			label := lc.Objects[1].(*widget.Label)
			price := lc.Objects[2].(*widget.Label)
			change := lc.Objects[3].(*widget.Label)
			duration := lc.Objects[4].(*widget.Label)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				btn.Refresh()
			}

			go func() {
				for {
					swaps := trades[pair].swaps

					n := unitrade.Name(swaps.Data.Swaps[0])
					p, c, d := unitrades.LastPriceChanges(swaps)
					fmt.Println(n)

					label.SetText(n)
					price.SetText(money.FormatMoney(p))
					change.SetText(money.FormatMoney(c))
					duration.SetText(fmt.Sprintf("%.2f hours", d))

					url := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)
					dex.SetURL(parseURL(url))
				}
			}()
		})

	go func() {
		pairs := make(chan []uniswap.Pair, 1)

		go services.UniswapMarkketPairs(pairs)
		msg := <-pairs

		go func() {
			for _, v := range msg {
				var wg sync.WaitGroup
				wg.Add(1)

				var t Trade
				var swaps uniswap.Swaps

				cc := make(chan string, 1)
				go uniswap.SwapsByCounts(cc, 2, v.Id)
				msg := <-cc
				json.Unmarshal([]byte(msg), &swaps)

				t.pair = v
				t.swaps = swaps
				trades[v.Id] = t
				pairsList.Append(v.Id)
				fmt.Println(unitrade.Name(swaps.Data.Swaps[0]))

				defer wg.Done()
			}
		}()
	}()

	controls := container.NewVBox(find, infProgress)
	return container.NewBorder(controls, nil, nil, nil, list)
}
