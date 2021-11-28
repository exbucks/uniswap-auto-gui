package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	uniswap "github.com/hirokimoto/uniswap-api"
)

func UniswapMarkketPairs(target chan<- []uniswap.Pair) {
	var pairs []uniswap.Pair
	skip := 0

	for {
		var wg sync.WaitGroup
		wg.Add(1)

		cc := make(chan string, 1)
		go uniswap.RequestPairs(cc, 1000, 1000*skip)
		msg := <-cc
		var data uniswap.Pairs
		json.Unmarshal([]byte(msg), &data)
		counts := len(data.Data.Pairs)
		fmt.Println("Got ", counts, "...")
		if counts == 0 {
			fmt.Println("Completed to get pairs ", len(pairs), ".....")
			target <- pairs
			defer wg.Done()
			return
		}
		pairs = append(pairs, data.Data.Pairs...)
		fmt.Println(skip, "* 1000 pairs is coming...")
		skip += 1

		defer wg.Done()
	}
}

func Notify(title string, message string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: message,
	})
}
