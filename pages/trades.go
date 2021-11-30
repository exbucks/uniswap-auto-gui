package pages

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	regression "github.com/gaillard/go-online-linear-regression/v1"
	uniswap "github.com/hirokimoto/uniswap-api"
	unitrade "github.com/hirokimoto/uniswap-api/swap"
	unitrades "github.com/hirokimoto/uniswap-api/swaps"
	"github.com/hirokimoto/uniswap-auto-gui/data"
	"github.com/hirokimoto/uniswap-auto-gui/services"
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

		var list []string
		for _, item := range pairs {
			list = append(list, item.Id)
		}

		data.SaveTradePairs(list)
	}()

	controls := container.NewVBox(find, infProgress)
	return container.NewBorder(controls, nil, nil, nil, table)
}

func trackTrade(pair string, index int, progress chan<- int) {
	duration, _ := strconv.Atoi(os.Getenv("SWAP_DURATION"))
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan string, 1)
	if duration > 100 {
		go uniswap.SwapsByCounts(ch, duration, pair)
	} else {
		go uniswap.SwapsByDays(ch, duration, pair)
	}

	msg := <-ch
	var swaps uniswap.Swaps
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		n := unitrade.Name(swaps.Data.Swaps[0])
		p, _ := unitrade.Price(swaps.Data.Swaps[0])
		_, c := unitrades.WholePriceChanges(swaps)
		_, _, d := unitrades.Duration(swaps)
		old, _ := unitrade.Old(swaps.Data.Swaps[0])
		average := unitrades.AveragePrice(swaps.Data.Swaps)

		// Filter our some tokens which is in the active trading in recent3 days.
		if old < 3*24 && p > 0.0001 {
			slope, _, _ := testRegression(swaps)
			var isGoingUp = slope > 0
			var isGoingDown = slope < 0
			// var isGoingUp = checkupOfSwaps(swaps)
			// var isGoingDown = checkdownOfSwaps(swaps)
			var isStable = math.Abs((average-p)/p) < 0.1
			var isUnStable = math.Abs((average-p)/p) > 0.1

			target := ""
			updown := ""

			if isUnStable {
				target = "unstable"
				fmt.Println("Unstable token ", n, p, average, c, d)
			}
			if isStable {
				target = "stable"
				fmt.Println("Stable token ", n, p, average, c, d)
			}
			if isGoingUp {
				updown = "up"
				fmt.Println("Trending up token ", n, p, average, c, d)
			}
			if isGoingDown {
				updown = "down"
				fmt.Println("Trending down token ", n, p, average, c, d)
			}

			if isUnStable || isStable || isGoingUp || isGoingDown {

			}
		}
	}
	fmt.Print(index, "|")

	defer wg.Done()
	progress <- index
}

func testRegression(swaps uniswap.Swaps) (float64, float64, float64) {
	r := regression.New(7)

	for i := 0; i < len(swaps.Data.Swaps); i++ {
		swap := swaps.Data.Swaps[i]
		price, _ := unitrade.Price(swap)
		r.Add(float64(i), price)
	}

	slope, intercept, stdError := r.CalculateWithStdError()
	return slope, intercept, stdError
}
