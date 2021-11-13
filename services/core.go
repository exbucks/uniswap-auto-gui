package services

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/uniswap-auto-gui/data"
)

func Notify(title string, message string, link string) {
	logo := canvas.NewImageFromResource(data.FyneScene)
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = gosxnotifier.Default
	note.Link = link
	note.AppIcon = logo.Resource.Name()
	note.Push()
}

func Alert(title string, message string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: message,
	})
}
