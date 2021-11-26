package data

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func SaveTrackPairs(pairs []string) {
	path := absolutePath() + "/pairs.txt"
	err := writeLines(pairs, path)

	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Error",
			Content: "Failed tracking pairs!",
		})
		return
	}

	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Success",
		Content: "Saved tracking pairs successfully!",
	})
}

func ReadTrackPairs() []string {
	path := absolutePath() + "/pairs.txt"
	pairs, err := readLines(path)

	if err != nil {
		fmt.Println(err)
		return []string{
			"0x7a99822968410431edd1ee75dab78866e31caf39",
			"0x3dd49f67e9d5bc4c5e6634b3f70bfd9dc1b6bd74"}
	}

	return pairs
}
