package services

import (
	"fyne.io/fyne/v2"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func Notify(title string, message string, link string) {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = gosxnotifier.Default
	note.Link = link
	note.AppIcon = "gopher.png"
	note.Push()
}

func Alert(title string, message string) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: message,
	})
}
