package data

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func SaveTrackSettings(settings []byte) {
	path := absolutePath() + "/settings.txt"
	fmt.Println(settings)
	err := writeBytes(settings, path)

	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Error",
			Content: "Failed saving settings!",
		})
		return
	}

	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Success",
		Content: "Saved settings successfully!",
	})
}

func ReadTrackSettings() ([]byte, error) {
	path := absolutePath() + "/settings.txt"
	bytes, err := readBytes(path)
	fmt.Println(bytes)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bytes, nil
}
