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
	gosxnotifier "github.com/deckarep/gosx-notifier"
	uniswap "github.com/hirokimoto/uniswap-api"
	unitrade "github.com/hirokimoto/uniswap-api/swap"
	unitrades "github.com/hirokimoto/uniswap-api/swaps"
	"github.com/leekchan/accounting"
	"github.com/uniswap-auto-gui/data"
	"github.com/uniswap-auto-gui/services"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	var selected uniswap.Swaps
	var activePair string
	var oldPrices = map[string]float64{}

	money := accounting.Accounting{Symbol: "$", Precision: 6}
	trades := map[string]uniswap.Swaps{}

	pairs := data.ReadTrackPairs()
	pairsdata := binding.BindStringList(&pairs)

	name := widget.NewEntry()
	name.SetPlaceHolder("0x385769E84B650C070964398929DB67250B7ff72C")
	append := widget.NewButton("Append", func() {
		if name.Text != "" {
			isExisted := false
			for _, item := range pairs {
				if item == name.Text {
					isExisted = true
				}
			}
			if !isExisted {
				pairs = append(pairs, name.Text)
				pairsdata.Reload()
				data.SaveTrackPairs(pairs)
			}
		}
	})

	control := container.NewVBox(name, append)

	rightList := widget.NewList(
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

	leftList := widget.NewListWithData(pairsdata,
		func() fyne.CanvasObject {
			left := container.NewHBox(widget.NewHyperlink("DEX", parseURL("https://github.com/hirokimoto")), widget.NewLabel("token"), widget.NewLabel("price"), widget.NewLabel("change"), widget.NewLabel("duration"))
			right := container.NewHBox(widget.NewButton("#", nil), widget.NewButton("-", nil))
			return container.NewBorder(nil, nil, left, right)
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

			right := obj.(*fyne.Container).Objects[1].(*fyne.Container)
			btnSettings := right.Objects[0].(*widget.Button)
			btnSettings.OnTapped = func() {
				w := fyne.CurrentApp().NewWindow("Settings")

				min, max := data.ReadMinMax(pair)
				mindata := binding.BindFloat(&min)
				minLabel := widget.NewLabel("Minimum")
				minEntry := widget.NewEntryWithData(binding.FloatToString(mindata))
				minPanel := container.NewGridWithColumns(2, minLabel, minEntry)

				maxdata := binding.BindFloat(&max)
				maxLabel := widget.NewLabel("Maximum")
				maxEntry := widget.NewEntryWithData(binding.FloatToString(maxdata))
				maxPanel := container.NewGridWithColumns(2, maxLabel, maxEntry)

				btnSave := widget.NewButton("Save", func() {
					data.SaveTrackSettings(pair, min, max)
				})

				settingsPanel := container.NewVBox(minPanel, maxPanel, btnSave)
				w.SetContent(settingsPanel)

				w.Resize(fyne.NewSize(340, 180))
				w.SetFixedSize(true)
				w.Show()
			}
			btnRemove := right.Objects[1].(*widget.Button)
			btnRemove.OnTapped = func() {}

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
					if len(trades[pair].Data.Swaps) > 0 &&
						(trades[pair].Data.Swaps[0].Id == swaps.Data.Swaps[0].Id ||
							trades[pair].Data.Swaps[0].Timestamp > swaps.Data.Swaps[0].Timestamp) {
						continue
					}

					n := unitrade.Name(swaps.Data.Swaps[0])
					p, _ := unitrade.Price(swaps.Data.Swaps[0])
					_, c := unitrades.WholePriceChanges(swaps)
					_, _, d := unitrades.Duration(swaps)

					label.SetText(n)
					_price := money.FormatMoney(p)
					price.SetText(_price)
					_change := money.FormatMoney(c)
					change.SetText(_change)
					_duration := fmt.Sprintf("%.2f hours", d)
					duration.SetText(_duration)

					if activePair == pair &&
						selected.Data.Swaps[0].Id != swaps.Data.Swaps[0].Id {
						selected = swaps
						rightList.Refresh()
					}
					trades[pair] = swaps

					if p != oldPrices[pair] {
						t := time.Now()
						message := fmt.Sprintf("%s: %f %f %f", n, p, c, d)
						title := "Priced Up!"
						if c < 0 {
							title = "Priced Down!"
						}
						link := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)

						min, max := data.ReadMinMax(pair)

						if p < min {
							title = fmt.Sprintf("Warning Low! Watch %s", n)
						}
						if p > max {
							title = fmt.Sprintf("Warning High! Watch %s", n)
						}

						services.Alert(title, message, link, gosxnotifier.Morse)

						fmt.Println(".")
						fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, p, c, d)
						fmt.Println(".")
					}
					oldPrices[pair] = p

					time.Sleep(time.Second * 1)
				}
			}()
		})

	leftList.OnSelected = func(id widget.ListItemID) {
		activePair, _ = pairsdata.GetValue(id)
		selected = trades[activePair]
		rightList.Refresh()
	}

	listPanel := container.NewBorder(nil, control, nil, nil, leftList)
	return container.NewHSplit(listPanel, rightList)
}
