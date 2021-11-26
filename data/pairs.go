package data

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne/v2"
)

func SaveTrackPairs(pairs []string) {
	file, err := os.CreateTemp("", "pairs.*.bat")
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(file)
	for _, pair := range pairs {
		fmt.Fprintln(w, pair)
	}
	w.Flush()

	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   "Success",
		Content: "Saved tracking pairs successfully!",
	})
}

func ReadTrackPairs() []string {
	file, err := ioutil.TempFile("", "pairs.*.bat")
	if err != nil {
		fmt.Println(err)
		return []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	}

	var pairs []string
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pairs = append(pairs, scanner.Text())
	}

	if err != nil || len(pairs) == 0 {
		fmt.Println(file.Name())
		fmt.Println(err)
		return []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	}
	return pairs
}
