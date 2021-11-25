package pages

import (
	"encoding/json"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/hirokimoto/crypto-auto/services"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

func tradesScreen(_ fyne.Window) fyne.CanvasObject {
	money := accounting.Accounting{Symbol: "$", Precision: 6}

	dataList := binding.BindStringList(&[]string{})

	infProgress := widget.NewProgressBarInfinite()
	infProgress.Stop()

	find := widget.NewButton("Find Trading Coins", func() {
		infProgress.Start()
		command := make(chan string)
		progress := make(chan int)
		tt := &services.Tokens{}
		go func() {
			for {
				go services.AnalyzePairs(command, progress, tt)
				time.Sleep(time.Minute * 20)
			}
		}()
	})

	list := widget.NewListWithData(dataList,
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
					var swaps utils.Swaps
					c1 := make(chan string, 1)
					utils.SwapsByCounts(c1, 2, pair)
					msg := <-c1
					json.Unmarshal([]byte(msg), &swaps)
					n, p, c, d, _, _ := services.SwapsInfo(swaps, 0.1)
					label.SetText(n)
					price.SetText(money.FormatMoney(p))
					change.SetText(money.FormatMoney(c))
					duration.SetText(fmt.Sprintf("%.2f hours", d))

					url := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)
					dex.SetURL(parseURL(url))
				}
			}()
		})

	controls := container.NewVBox(find, infProgress)
	return container.NewBorder(controls, nil, nil, nil, list)
}
