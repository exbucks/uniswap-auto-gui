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

func tradableScreen(_ fyne.Window) fyne.CanvasObject {
	dataList := binding.BindStringList(&[]string{})

	find := widget.NewButton("Find", func() {
		go func() {
			for {
				c1 := make(chan string, 1)
				go utils.Post(c1, "pairs", "")
				trackTradables(c1, dataList)
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
			f := item.(binding.String)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.Bind(f)

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				fmt.Println("Ok!")
			}
		})

	return container.NewBorder(find, nil, nil, nil, list)
}

func trackTradables(pings <-chan string, list binding.ExternalStringList) {
	msg := <-pings
	var pairs utils.Pairs

	json.Unmarshal([]byte(msg), &pairs)

	var wg sync.WaitGroup
	wg.Add(len(pairs.Data.Pairs))
	go services.TradableTokens(&wg, pairs, list)
	wg.Wait()
}
