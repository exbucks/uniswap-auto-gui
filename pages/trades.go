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
	pairsList := binding.BindStringList(&[]string{})
	trades := map[string]Trade{}

	infProgress := widget.NewProgressBarInfinite()
	infProgress.Stop()
	// command := make(chan string)
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
				money := accounting.Accounting{Symbol: "$", Precision: 6}
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
		pc := make(chan []uniswap.Pair, 1)

		go services.UniswapMarkketPairs(pc)
		msg := <-pc

		go func() {
			for _, v := range msg {
				var wg sync.WaitGroup
				wg.Add(1)

				var t Trade
				var s uniswap.Swaps

				sc := make(chan string, 1)
				go uniswap.SwapsByCounts(sc, 2, v.Id)
				msg := <-sc
				json.Unmarshal([]byte(msg), &s)

				if len(s.Data.Swaps) > 0 {
					t.pair = v
					t.swaps = s
					trades[v.Id] = t
					fmt.Println(unitrade.Name(s.Data.Swaps[0]))
				}

				defer wg.Done()
			}

			for _, v := range trades {
				pairsList.Append(v.pair.Id)
			}
		}()
	}()

	controls := container.NewVBox(find, infProgress)
	return container.NewBorder(controls, nil, nil, nil, list)
}
