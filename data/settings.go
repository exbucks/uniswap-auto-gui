package data

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
)

type Setting struct {
	min float64
	max float64
}

func SaveTrackSettings(settings map[string]Setting) {
	path := absolutePath() + "/settings.txt"

	t, _ := json.Marshal(settings)

	err := writeBytes(t, path)

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

func ReadTrackSettings() (map[string]Setting, error) {
	path := absolutePath() + "/settings.txt"
	bytes, err := readBytes(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var settings map[string]Setting
	json.Unmarshal([]byte(bytes), &settings)

	return settings, nil
}
