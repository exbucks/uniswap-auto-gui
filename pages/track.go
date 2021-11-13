package pages

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/uniswap-auto-gui/services"
	"github.com/uniswap-auto-gui/utils"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	dataList := binding.BindStringList(&[]string{"0x9d9681d71142049594020bd863d34d9f48d9df58", "0x7a99822968410431edd1ee75dab78866e31caf39"})

	append := widget.NewButton("Append", func() {
		dataList.Append("0x7a99822968410431edd1ee75dab78866e31caf39")
	})

	list := widget.NewListWithData(dataList,
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Address: "), widget.NewLabel("address"), widget.NewLabel("Price x"), widget.NewButton("Track", nil))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s := item.(binding.String)
			address := obj.(*fyne.Container).Objects[1].(*widget.Label)
			address.Bind(s)

			f := binding.NewFloat()
			f.Set(0.1)
			price := obj.(*fyne.Container).Objects[2].(*widget.Label)
			price.Bind(binding.FloatToStringWithFormat(f, "Price %f"))

			btn := obj.(*fyne.Container).Objects[3].(*widget.Button)
			btn.OnTapped = func() {
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
							price := services.PriceFromSwaps(eth, swaps)
							f.Set(price)
							// fmt.Println("Current price: ", price)
							// trackSwap(c2, f)
						case msg2 := <-c2:
							json.Unmarshal([]byte(msg2), &swaps)
							price := services.PriceFromSwaps(eth, swaps)
							f.Set(price)
							// fmt.Println("Current price: ", price)
							// trackSwap(c2, f)
						}
					}
				}()
			}
		})
	listPanel := container.NewBorder(nil, append, nil, nil, list)
	return container.NewGridWithColumns(1, listPanel)
}

func newFormWithData(data binding.DataMap) *widget.Form {
	keys := data.Keys()
	items := make([]*widget.FormItem, len(keys))
	for i, k := range keys {
		data, err := data.GetItem(k)
		if err != nil {
			items[i] = widget.NewFormItem(k, widget.NewLabel(err.Error()))
		}
		items[i] = widget.NewFormItem(k, createBoundItem(data))
	}

	return widget.NewForm(items...)
}

func createBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}

func trackSwap(pings <-chan string, price binding.Float) {
	msg := <-pings
	var swaps utils.Swaps
	json.Unmarshal([]byte(msg), &swaps)

	min, max, minTarget, maxTarget, minTime, maxTime := services.MinAndMax(swaps)
	fmt.Println("Min price: ", min, minTarget, minTime)
	fmt.Println("Max price: ", max, maxTarget, maxTime)

	last := services.LastPrice(swaps)
	_ = price.Set(last)
	fmt.Println("Last price: ", last)

	ts, tl, period := services.PeriodOfSwaps(swaps)
	fmt.Println("Timeframe of 100 swaps: ", period)
	fmt.Println("Start and End time of the above time frame: ", ts, tl)
	if (max-min)/last > 0.5 {
		fmt.Println("$$$$$ This is a tradable token! $$$$$")
	}
}

func trackPairs(pings <-chan string) {
	msg := <-pings
	var pairs utils.Pairs

	json.Unmarshal([]byte(msg), &pairs)

	var wg sync.WaitGroup
	wg.Add(len(pairs.Data.Pairs))
	go services.TradableTokens(&wg, pairs)
	wg.Wait()
}
