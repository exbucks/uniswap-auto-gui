package pages

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/uniswap-auto-gui/services"
	"github.com/uniswap-auto-gui/utils"
)

func stableScreen(_ fyne.Window) fyne.CanvasObject {
	table := widget.NewTable(
		func() (int, int) { return 500, 150 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row+1))
			case 1:
				label.SetText("A longer cell")
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	table.SetColumnWidth(0, 34)
	table.SetColumnWidth(1, 102)

	button := widget.NewButton("Find", func() {
		go func() {
			for {
				c1 := make(chan string)
				go utils.Post(c1, "pairs", "")
				trackPairs(c1)
			}
		}()
	})

	return container.NewBorder(button, nil, nil, nil, table)
}

func trackPairs(pings <-chan string) {
	msg := <-pings
	var pairs utils.Pairs

	json.Unmarshal([]byte(msg), &pairs)

	var wg sync.WaitGroup
	wg.Add(len(pairs.Data.Pairs))
	go services.StableTokens(&wg, pairs)
	wg.Wait()
}
