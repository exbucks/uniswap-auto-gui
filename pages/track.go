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
	"github.com/uniswap-auto-gui/services"
	"github.com/uniswap-auto-gui/utils"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	var selected utils.Swaps
	var oldPrice float64
	var activePair string
	alertType := "Alert any changes!"

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
		alertType = s
	})
	alerts.Horizontal = true
	alerts.SetSelected("Alert any changes!")

	pairs := services.ReadPairs()
	pairsdata := binding.BindStringList(&pairs)

	name := widget.NewEntry()
	name.SetPlaceHolder("0x7a99822968410431edd1ee75dab78866e31caf39")
	append := widget.NewButton("Append", func() {
		if name.Text != "" {
			pairsdata.Append(name.Text)
			services.WritePairs(pairs)
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
			price, target, amount, amount1, amount2 := services.SwapInfo(selected.Data.Swaps[id])
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(target)
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("%f", price))
			item.(*fyne.Container).Objects[3].(*widget.Label).SetText(amount)
			item.(*fyne.Container).Objects[4].(*widget.Label).SetText(amount1)
			item.(*fyne.Container).Objects[5].(*widget.Label).SetText(amount2)
		},
	)

	list := widget.NewListWithData(pairsdata,
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewHyperlink("DEX", parseURL("https://github.com/hirokimoto")), widget.NewLabel("token"), widget.NewLabel("price"), widget.NewLabel("change"), widget.NewLabel("duration"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s := item.(binding.String)

			dex := obj.(*fyne.Container).Objects[0].(*widget.Hyperlink)
			pair, _ := s.Get()
			url := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)
			dex.SetURL(parseURL(url))

			label := obj.(*fyne.Container).Objects[1].(*widget.Label)
			price := obj.(*fyne.Container).Objects[2].(*widget.Label)
			change := obj.(*fyne.Container).Objects[3].(*widget.Label)
			duration := obj.(*fyne.Container).Objects[4].(*widget.Label)

			var swaps utils.Swaps
			cc := make(chan string, 1)

			go func() {
				for {
					utils.Post(cc, "swaps", pair)

					msg := <-cc
					json.Unmarshal([]byte(msg), &swaps)
					n, p, c, d, a := services.SwapsInfo(swaps, ai)
					label.SetText(n)
					price.SetText(fmt.Sprintf("%f", p))
					change.SetText(fmt.Sprintf("%f", c))
					duration.SetText(fmt.Sprintf("%.2f hours", d))

					if activePair == pair {
						selected = swaps
						oldPrice, _, _, _, _ = services.SwapInfo(swaps.Data.Swaps[0])
						trades.Refresh()
					}
					if a {
						services.Notify("Alert large changes!", n, url)
					}
					if alertType == "Alert any changes!" {
						if pair == activePair && oldPrice != p {
							services.Notify("Price any changes!", fmt.Sprintf("%s $%f  $%f", n, p, c), url)
							oldPrice = p
						}
					} else {
						if pair == activePair && oldPrice != p {
							oldPrice = p
							services.Notify("Price special changes!", fmt.Sprintf("%s $%f  $%f", n, p, c), url)
						}
					}
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
