package pages

import (
	"encoding/json"
	"fmt"
	"sync"
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
	"github.com/skratchdot/open-golang/open"
	"github.com/uniswap-auto-gui/data"
	"github.com/uniswap-auto-gui/services"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	var selected uniswap.Swaps

	pairs := data.ReadTrackPairs()
	records, _ := data.ReadTrackSettings()

	oldNames := make([]string, 0)
	oldPrices := make([]float64, 0)
	oldChanges := make([]float64, 0)
	oldDurations := make([]float64, 0)

	for _, _ = range pairs {
		oldNames = append(oldNames, "")
		oldPrices = append(oldPrices, 0.0)
		oldChanges = append(oldChanges, 0.0)
		oldDurations = append(oldDurations, 0.0)
	}

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
				oldNames = append(oldNames, "")
				oldPrices = append(oldPrices, 0.0)
				oldChanges = append(oldChanges, 0.0)
				oldDurations = append(oldDurations, 0.0)
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
			if len(selected.Data.Swaps) > 1 {
				price, target, amount, amount1, amount2 := unitrade.Trade(selected.Data.Swaps[id])

				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(target)
				item.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("$%f", price))
				item.(*fyne.Container).Objects[3].(*widget.Label).SetText(amount)
				item.(*fyne.Container).Objects[4].(*widget.Label).SetText(amount1)
				item.(*fyne.Container).Objects[5].(*widget.Label).SetText(amount2)
			}
		},
	)

	table := widget.NewTable(
		func() (int, int) { return len(pairs), 5 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row+1))
			case 1:
				if len(oldNames[id.Row]) > 20 {
					label.SetText(oldNames[id.Row][0:20] + "...")
				} else {
					label.SetText(oldNames[id.Row])
				}
			case 2:
				label.SetText(fmt.Sprintf("%f", oldPrices[id.Row]))
			case 3:
				label.SetText(fmt.Sprintf("%f", oldChanges[id.Row]))
			case 4:
				label.SetText(fmt.Sprintf("%f", oldDurations[id.Row]))
			default:
			}
		})
	table.SetColumnWidth(0, 60)
	table.SetColumnWidth(1, 202)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)
	table.OnSelected = func(id widget.TableCellID) {
		pair := pairs[id.Row]
		if id.Col == 0 {
			open.Run(fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair))
		}
		if id.Col == 1 {
			go func() {
				for {
					var swaps uniswap.Swaps
					cc := make(chan string, 1)
					go uniswap.SwapsByCounts(cc, 20, pair)

					msg := <-cc
					json.Unmarshal([]byte(msg), &swaps)
					selected = swaps
					time.Sleep(time.Second * 1)
				}
			}()
		}
		if id.Col == 2 {
			records, _ = data.ReadTrackSettings()
			w := fyne.CurrentApp().NewWindow("Settings")

			min, max := data.ReadMinMax(records, pair)
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
	}

	listPanel := container.NewBorder(nil, control, nil, nil, table)

	go func() {
		for {
			for index, pair := range pairs {
				var wg sync.WaitGroup
				wg.Add(1)
				fmt.Print(".")

				var swaps uniswap.Swaps
				cc := make(chan string, 1)
				go uniswap.SwapsByCounts(cc, 2, pair)

				msg := <-cc
				json.Unmarshal([]byte(msg), &swaps)

				if len(swaps.Data.Swaps) == 0 || swaps.Data.Swaps == nil {
					time.Sleep(time.Second * 1)
					continue
				}

				n := unitrade.Name(swaps.Data.Swaps[0])
				p, _ := unitrade.Price(swaps.Data.Swaps[0])
				_, c := unitrades.WholePriceChanges(swaps)
				_, _, d := unitrades.Duration(swaps)

				if oldPrices[index] != p {
					go alert(records, pair, n, p, c, d)
					table.Refresh()
					oldPrices[index] = p
				}

				oldNames[index] = n
				oldChanges[index] = c
				oldDurations[index] = d

				wg.Done()
			}
			time.Sleep(time.Second * 1)
		}
	}()

	return container.NewHSplit(listPanel, rightList)
}

func alert(records [][]string, pair string, n string, p float64, c float64, d float64) {
	message := fmt.Sprintf("%s: %f %f %f", n, p, c, d)
	title := "Priced Up!"
	if c < 0 {
		title = "Priced Down!"
	}
	link := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair)

	min, max := data.ReadMinMax(records, pair)
	sound := gosxnotifier.Morse

	if p < min {
		title = fmt.Sprintf("Warning Low! Watch %s", n)
		sound = gosxnotifier.Default
	}
	if p > max {
		title = fmt.Sprintf("Warning High! Watch %s", n)
		sound = gosxnotifier.Default
	}

	services.Alert(title, message, link, sound)
}
