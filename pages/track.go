package pages

import (
	"encoding/json"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	uniswap "github.com/hirokimoto/uniswap-api"
	unitrade "github.com/hirokimoto/uniswap-api/swap"
	unitrades "github.com/hirokimoto/uniswap-api/swaps"
	"github.com/leekchan/accounting"
	"github.com/uniswap-auto-gui/data"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	var selected uniswap.Swaps
	var activePair string
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	transactions := map[string]string{}

	ai := 0.1
	aidata := binding.BindFloat(&ai)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(aidata, "Price change percent (*100): %f"))
	entry := widget.NewEntryWithData(binding.FloatToString(aidata))
	alertInterval := container.NewBorder(nil, nil, label, entry)

	min := 0.3
	mindata := binding.BindFloat(&min)
	minLabel := widget.NewLabel("Minimum")
	minEntry := widget.NewEntryWithData(binding.FloatToString(mindata))

	max := 0.4
	maxdata := binding.BindFloat(&max)
	maxLabel := widget.NewLabel("Maximum")
	maxEntry := widget.NewEntryWithData(binding.FloatToString(maxdata))

	alerts := widget.NewRadioGroup([]string{"Alert any changes!", "Alert special changes!"}, func(s string) {
		fmt.Println(s)
	})
	alerts.Horizontal = true
	alerts.SetSelected("Alert any changes!")

	pairs := data.ReadTrackPairs()
	pairsdata := binding.BindStringList(&pairs)

	name := widget.NewEntry()
	name.SetPlaceHolder("0x385769E84B650C070964398929DB67250B7ff72C")
	append := widget.NewButton("Append", func() {
		if name.Text != "" {
			pairs = append(pairs, name.Text)
			pairsdata.Reload()
			data.SaveTrackPairs(pairs)
		}
	})

	control := container.NewVBox(name, append)

	trades := widget.NewList(
		func() int {
			return len(selected.Data.Swaps)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("target"), widget.NewLabel("price"), widget.NewLabel("amount"), widget.NewLabel("amount1"), widget.NewLabel("amount2"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			price, target, amount, amount1, amount2 := unitrade.Trade(selected.Data.Swaps[id])

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(target)
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("$%f", price))
			item.(*fyne.Container).Objects[3].(*widget.Label).SetText(amount)
			item.(*fyne.Container).Objects[4].(*widget.Label).SetText(amount1)
			item.(*fyne.Container).Objects[5].(*widget.Label).SetText(amount2)
		},
	)

	list := widget.NewListWithData(pairsdata,
		func() fyne.CanvasObject {
			left := container.NewHBox(widget.NewHyperlink("DEX", parseURL("https://github.com/hirokimoto")), widget.NewLabel("token"), widget.NewLabel("price"), widget.NewLabel("change"), widget.NewLabel("duration"))
			return container.NewBorder(nil, nil, left, widget.NewButton("-", nil))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s := item.(binding.String)

			left := obj.(*fyne.Container).Objects[0].(*fyne.Container)
			dex := left.Objects[0].(*widget.Hyperlink)
			pair, _ := s.Get()
			url := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)
			dex.SetURL(parseURL(url))

			label := left.Objects[1].(*widget.Label)
			price := left.Objects[2].(*widget.Label)
			change := left.Objects[3].(*widget.Label)
			duration := left.Objects[4].(*widget.Label)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {}

			var swaps uniswap.Swaps
			cc := make(chan string, 1)

			go func() {
				for {
					go uniswap.SwapsByCounts(cc, 100, pair)

					msg := <-cc
					json.Unmarshal([]byte(msg), &swaps)

					if len(swaps.Data.Swaps) == 0 {
						time.Sleep(time.Second * 5)
						continue
					}
					if transactions[pair] == swaps.Data.Swaps[0].Id {
						continue
					}

					n := unitrade.Name(swaps.Data.Swaps[0])
					p, _ := unitrade.Price(swaps.Data.Swaps[0])
					_, c := unitrades.WholePriceChanges(swaps)
					_, _, d := unitrades.Duration(swaps)

					if label.Text != n {
						label.SetText(n)
					}
					_price := money.FormatMoney(p)
					if price.Text != _price {
						price.SetText(_price)
					}
					_change := money.FormatMoney(c)
					if change.Text != _change {
						change.SetText(_change)
					}
					_duration := fmt.Sprintf("%.2f hours", d)
					if duration.Text != _duration {
						duration.SetText(_duration)
					}

					if activePair == pair {
						selected = swaps
						trades.Refresh()
					}
					transactions[pair] = swaps.Data.Swaps[0].Id
					time.Sleep(time.Second * 5)
				}
			}()
		})

	list.OnSelected = func(id widget.ListItemID) {
		activePair, _ = pairsdata.GetValue(id)
	}

	minPanel := container.NewHBox(minLabel, minEntry)
	maxPanel := container.NewHBox(maxLabel, maxEntry)
	alertSpecialPanel := container.NewBorder(nil, nil, minPanel, maxPanel)
	settings := container.NewVBox(alertInterval, alerts, alertSpecialPanel)
	listPanel := container.NewBorder(settings, control, nil, nil, list)
	return container.NewHSplit(listPanel, trades)
}
