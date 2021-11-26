package data

import "fyne.io/fyne/v2"

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
