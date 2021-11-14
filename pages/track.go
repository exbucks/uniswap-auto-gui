package pages

import (
	"encoding/json"
	"fmt"

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

	ai := 0.1
	aidata := binding.BindFloat(&ai)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(aidata, "Price change percent (*100): %f"))
	entry := widget.NewEntryWithData(binding.FloatToString(aidata))
	alertInterval := container.NewBorder(nil, nil, label, entry)

	ac := false
	acdata := binding.BindBool(&ac)
	alertChange := widget.NewCheckWithData("Alert changes!", acdata)

	pairs := binding.BindStringList(&[]string{"0x9d9681d71142049594020bd863d34d9f48d9df58", "0x7a99822968410431edd1ee75dab78866e31caf39"})

	name := widget.NewEntry()
	name.SetPlaceHolder("0x7a99822968410431edd1ee75dab78866e31caf39")
	append := widget.NewButton("Append", func() {
		if name.Text != "" {
			pairs.Append(name.Text)
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

	list := widget.NewListWithData(pairs,
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
					pair, _ := s.Get()
					utils.Post(cc, "swaps", pair)

					msg := <-cc
					json.Unmarshal([]byte(msg), &swaps)
					n, p, c, d, a := services.SwapsInfo(swaps, ai)
					label.SetText(n)
					price.SetText(fmt.Sprintf("%f", p))
					change.SetText(fmt.Sprintf("%f", c))
					duration.SetText(fmt.Sprintf("%.2f hours", d))
					if a {
						services.Notify("Price Change Alert", n, url)
					}
					alert, _ := acdata.Get()
					if alert {
						services.Notify("Price changed!", fmt.Sprintf("%s %f", n, p), url)
					}
				}
			}()
		})

	list.OnSelected = func(id widget.ListItemID) {
		go func() {
			cc := make(chan string, 1)
			pair, _ := pairs.GetValue(id)
			utils.Post(cc, "swaps", pair)
			msg := <-cc
			json.Unmarshal([]byte(msg), &selected)
			trades.Refresh()
		}()
	}

	settings := container.NewVBox(alertInterval, alertChange)
	listPanel := container.NewBorder(settings, control, nil, nil, list)
	return container.NewHSplit(listPanel, trades)
}
