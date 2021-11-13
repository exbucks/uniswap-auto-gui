package pages

import (
	"encoding/json"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/uniswap-auto-gui/services"
	"github.com/uniswap-auto-gui/utils"
)

func tradableScreen(_ fyne.Window) fyne.CanvasObject {
	dataList := binding.BindFloatList(&[]float64{0.1, 0.2, 0.3})

	button := widget.NewButton("Append", func() {
		// dataList.Append(float64(dataList.Length()+1) / 10)

		go func() {
			for {
				c1 := make(chan string, 1)
				go utils.Post(c1, "pairs", "")
				trackTradables(c1)
				time.Sleep(time.Minute * 20)
			}
		}()
	})

	list := widget.NewListWithData(dataList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("+", nil),
				widget.NewLabel("item x.y"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.Bind(binding.FloatToStringWithFormat(f, "item %0.1f"))

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				_ = f.Set(val + 1)
			}
		})

	return container.NewBorder(button, nil, nil, nil, list)
}

func trackTradables(pings <-chan string) {
	msg := <-pings
	var pairs utils.Pairs

	json.Unmarshal([]byte(msg), &pairs)

	var wg sync.WaitGroup
	wg.Add(len(pairs.Data.Pairs))
	go services.TradableTokens(&wg, pairs)
	wg.Wait()
}
