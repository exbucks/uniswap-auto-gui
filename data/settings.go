package data

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
)

func SaveTrackSettings() {
	records := [][]string{
		{"first_name", "last_name", "occupation"},
		{"John", "Doe", "gardener"},
		{"Lucy", "Smith", "teacher"},
		{"Brian", "Bethamy", "programmer"},
	}

	filePath := absolutePath() + "/settings.csv"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	err = csvWriter.WriteAll(records)

	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

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

func ReadTrackSettings() ([][]string, error) {
	filePath := absolutePath() + "/settings.csv"
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	fmt.Println(records)

	return records, nil
}
