package pages

import (
	"fyne.io/fyne/v2"
)

// Page defines the data structure for a page
type Page struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Pages defines the metadata for each page
	Pages = map[string]Page{
		"welcome":  {"Welcome", "", welcomeScreen},
		"coins":    {"Coins", "", coinsScreen},
		"track":    {"Track", "", trackScreen},
		"tradable": {"Tradable", "", tradableScreen},
		"stable":   {"Stable", "", stableScreen},
	}

	// PageIndex  defines how our pages should be laid out in the index tree
	PageIndex = map[string][]string{
		"":            {"welcome", "coins", "track", "tradable", "stable"},
		"collections": {"list", "table", "tree"},
		"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split"},
		"widgets":     {"accordion", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
