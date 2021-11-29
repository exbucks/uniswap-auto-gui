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

// func analyzePairs(command <-chan string, progress chan<- int, t *Tokens) {
// 	pairs, _ := ReadAllPairs()
// 	t.SetTotal(len(pairs))
// 	var status = "Play"
// 	for index, pair := range pairs {
// 		select {
// 		case cmd := <-command:
// 			fmt.Println(cmd)
// 			switch cmd {
// 			case "Stop":
// 				return
// 			case "Pause":
// 				status = "Pause"
// 			default:
// 				status = "Play"
// 			}
// 		default:
// 			if status == "Play" {
// 				trackPair(pair, index, t, progress)
// 			}
// 		}
// 	}
// }

// func trackPair(pair string, index int, t *Tokens, progress chan<- int) {
// 	duration, _ := strconv.Atoi(os.Getenv("SWAP_DURATION"))
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	ch := make(chan string, 1)
// 	if duration > 100 {
// 		go utils.SwapsByCounts(ch, duration, pair)
// 	} else {
// 		go utils.SwapsByDays(ch, duration, pair)
// 	}

// 	msg := <-ch
// 	var swaps utils.Swaps
// 	json.Unmarshal([]byte(msg), &swaps)

// 	if len(swaps.Data.Swaps) > 0 {
// 		name, price, change, period, average, _ := SwapsInfo(swaps, 0.1)

// 		min, max, _, _, _, _ := minMax(swaps)
// 		howOld := howMuchOld(swaps)

// 		// Filter our some tokens which is in the active trading in recent3 days.
// 		if howOld < 3*24 && price > 0.0001 {
// 			slope, _, _ := testRegression(swaps)
// 			var isGoingUp = slope > 0
// 			var isGoingDown = slope < 0
// 			// var isGoingUp = checkupOfSwaps(swaps)
// 			// var isGoingDown = checkdownOfSwaps(swaps)
// 			var isStable = math.Abs((average-price)/price) < 0.1
// 			var isUnStable = math.Abs((average-price)/price) > 0.1

// 			target := ""
// 			updown := ""
// 			if isUnStable {
// 				target = "unstable"
// 				// Notify("Unstable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/", gosxnotifier.Blow)
// 				fmt.Println("Unstable token ", name, price, average, change, period)
// 			}
// 			if isStable {
// 				target = "stable"
// 				// Notify("Stable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/", gosxnotifier.Blow)
// 				fmt.Println("Stable token ", name, price, average, change, period)
// 			}
// 			if isGoingUp {
// 				updown = "up"
// 				fmt.Println("Trending up token ", name, price, average, change, period)
// 			}
// 			if isGoingDown {
// 				updown = "down"
// 				fmt.Println("Trending down token ", name, price, average, change, period)
// 			}

// 			if isUnStable || isStable || isGoingUp || isGoingDown {
// 				ct := &Token{
// 					target:  target,
// 					updown:  updown,
// 					name:    name,
// 					address: pair,
// 					price:   fmt.Sprintf("%f", price),
// 					change:  fmt.Sprintf("%f", change),
// 					min:     fmt.Sprintf("%f", min),
// 					max:     fmt.Sprintf("%f", max),
// 					period:  fmt.Sprintf("%.2f", period),
// 					swaps:   swaps.Data.Swaps,
// 				}
// 				t.Add(ct)
// 			}
// 		}
// 	}
// 	t.SetProgress(index)
// 	fmt.Print(index, "|")

// 	defer wg.Done()
// 	progress <- index
// }

// func testRegression(swaps utils.Swaps) (float64, float64, float64) {
// 	r := regression.New(7)

// 	for i := 0; i < len(swaps.Data.Swaps); i++ {
// 		swap := swaps.Data.Swaps[i]
// 		price, _ := priceOfSwap(swap)
// 		r.Add(float64(i), price)
// 	}

// 	slope, intercept, stdError := r.CalculateWithStdError()
// 	return slope, intercept, stdError
// }
