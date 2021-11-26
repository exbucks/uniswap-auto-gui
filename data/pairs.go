package data

import (
	"fyne.io/fyne/v2"
)

const TRADABLES string = "/tradables.txt"
const ALL_PAIRS string = "/allpairs.txt"

func SaveTrackPairs(pairs []string) {
	path := absolutePath() + "/tracks.txt"
	err := writeLines(pairs, path)
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Warning!",
			Content: "Failed to save tracking pairs!",
		})
	}
}

func ReadTrackPairs() []string {
	path := absolutePath() + "/tracks.txt"
	pairs, err := readLines(path)
	if err != nil || len(pairs) == 0 {
		return []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	}
	return pairs
}
