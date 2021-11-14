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
	activeIndex := 0

	pc := 0.1
	data := binding.BindFloat(&pc)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "Price change alert percent (*100): %f"))
	entry := widget.NewEntryWithData(binding.FloatToString(data))
	floats := container.NewGridWithColumns(2, label, entry)

	dataList := binding.BindStringList(&[]string{"0x9d9681d71142049594020bd863d34d9f48d9df58", "0x7a99822968410431edd1ee75dab78866e31caf39"})
	var selectedSwaps utils.Swaps

	name := widget.NewEntry()
	name.SetPlaceHolder("0x7a99822968410431edd1ee75dab78866e31caf39")
	append := widget.NewButton("Append", func() {
		if name.Text != "" {
			dataList.Append(name.Text)
		} else {
			services.Notify("test", "test", "https://www.github.com/hirokimoto")
		}
	})

	control := container.NewVBox(name, append)

	list := widget.NewListWithData(dataList,
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

			var eth utils.Crypto
			var swaps utils.Swaps
			c1 := make(chan string)
			c2 := make(chan string)

			go func() {
				for {
					pair, _ := s.Get()
					utils.Post(c1, "bundles", "")
					utils.Post(c2, "swaps", pair)
					time.Sleep(time.Second * 2)
				}
			}()
			go func() {
				for {
					select {
					case msg1 := <-c1:
						json.Unmarshal([]byte(msg1), &eth)
					case msg2 := <-c2:
						json.Unmarshal([]byte(msg2), &swaps)
						n, p, c, d, a := services.SwapsInfo(swaps, pc)
						selectedSwaps = swaps
						label.SetText(n)
						price.SetText(fmt.Sprintf("%f", p))
						change.SetText(fmt.Sprintf("%f", c))
						duration.SetText(fmt.Sprintf("%f hours", d))
						if a {
							services.Notify("Price Change Alert", n, url)
						}
						_pair, _ := dataList.GetValue(activeIndex)
						if pair == _pair {
							selectedSwaps = swaps
						}
					}
				}
			}()
		})

	list.OnSelected = func(id widget.ListItemID) {
		activeIndex = id
	}
	list.OnUnselected = func(id widget.ListItemID) {
		fmt.Println(id)
	}

	trades := widget.NewList(
		func() int {
			return len(selectedSwaps.Data.Swaps)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Token 0"), widget.NewLabel("Token 1"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(selectedSwaps.Data.Swaps[id].Pair.Token0.Name)
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(selectedSwaps.Data.Swaps[id].Pair.Token1.Name)
		},
	)

	listPanel := container.NewBorder(floats, control, nil, nil, list)
	return container.NewGridWithColumns(2, listPanel, trades)
}
