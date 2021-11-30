package data

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/hirokimoto/crypto-auto/services"
	"github.com/uniswap-auto-gui/services"
)

func SaveTrackPairs(pairs []string) {
	writePairs("/pairs.txt", pairs)
}

func ReadTrackPairs() []string {
	readPairs("/pairs.txt")
}

func SaveTradePairs(pairs []string) {
	writePairs("/trades.txt", pairs)
}

func ReadTradePairs() {
	readPairs("/trades.txt")
}

func readPairs(file string) []string {
	path := absolutePath() + file
	pairs, err := readLines(path)

	if err != nil {
		fmt.Println(err)
		return []string{
			"0x7a99822968410431edd1ee75dab78866e31caf39",
			"0x3dd49f67e9d5bc4c5e6634b3f70bfd9dc1b6bd74"}
	}

	return pairs
}

func writePairs(file string, pairs []string) {
	path := absolutePath() + file
	err := writeLines(pairs, path)

	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Error",
			Content: "Failed tracking pairs!",
		})
		return
	}

	services.Notify("Success", "Saved tracking pairs successfully!")
}
