package pages

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	uniswap "github.com/hirokimoto/uniswap-api"
	unitrade "github.com/hirokimoto/uniswap-api/swap"
	"github.com/uniswap-auto-gui/services"
)

type Trade struct {
	pair  uniswap.Pair
	swaps uniswap.Swaps
}

func tradesScreen(_ fyne.Window) fyne.CanvasObject {
	var pairs []uniswap.Pair
	trades := map[string]Trade{}

	table := widget.NewTable(
		func() (int, int) { return len(pairs), 7 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row+1))
			case 1:
				label.SetText(pairs[id.Row].Token0.Symbol)
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	table.SetColumnWidth(0, 34)
	table.SetColumnWidth(1, 102)

	infProgress := widget.NewProgressBarInfinite()
	infProgress.Stop()

	find := widget.NewButton("Find Trading Pairs", func() {
		infProgress.Start()

		go func() {
			for _, v := range pairs {
				var wg sync.WaitGroup
				wg.Add(1)

				var t Trade
				var s uniswap.Swaps

				sc := make(chan string, 1)
				go uniswap.SwapsByCounts(sc, 2, v.Id)
				msg := <-sc
				json.Unmarshal([]byte(msg), &s)

				if len(s.Data.Swaps) > 0 {
					t.pair = v
					t.swaps = s
					trades[v.Id] = t
					table.Refresh()
					fmt.Println(unitrade.Name(s.Data.Swaps[0]))
				}

				defer wg.Done()
			}

			infProgress.Stop()
		}()
	})

	go func() {
		pc := make(chan []uniswap.Pair, 1)
		go services.UniswapMarkketPairs(pc)
		pairs = <-pc
	}()

	controls := container.NewVBox(find, infProgress)
	return container.NewBorder(controls, nil, nil, nil, table)
}
