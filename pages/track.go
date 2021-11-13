package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func trackScreen(_ fyne.Window) fyne.CanvasObject {
	dataList := binding.BindStringList(&[]string{"0x295b42684f90c77da7ea46336001010f2791ec8c", "0xf8e9f10c22840b613cda05a0c5fdb59a4d6cd7ef", "0xe9cb6838902ccf711f16a9ea5a1170f8e9853c02"})

	append := widget.NewButton("Append", func() {
		dataList.Append("0x295b42684f90c77da7ea46336001010f2791ec8c")
	})

	list := widget.NewListWithData(dataList,
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Address: "), widget.NewLabel("address"), widget.NewLabel("Price x"), widget.NewButton("Track", nil))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s := item.(binding.String)
			address := obj.(*fyne.Container).Objects[1].(*widget.Label)
			address.Bind(s)

			f := binding.NewFloat()
			f.Set(0.1)
			price := obj.(*fyne.Container).Objects[2].(*widget.Label)
			price.Bind(binding.FloatToStringWithFormat(f, "Price %f"))

			btn := obj.(*fyne.Container).Objects[3].(*widget.Button)
			btn.OnTapped = func() {
				// val, _ := f.Get()
				// _ = f.Set(val + 1)
			}
		})
	listPanel := container.NewBorder(nil, append, nil, nil, list)
	return container.NewGridWithColumns(1, listPanel)
}

func newFormWithData(data binding.DataMap) *widget.Form {
	keys := data.Keys()
	items := make([]*widget.FormItem, len(keys))
	for i, k := range keys {
		data, err := data.GetItem(k)
		if err != nil {
			items[i] = widget.NewFormItem(k, widget.NewLabel(err.Error()))
		}
		items[i] = widget.NewFormItem(k, createBoundItem(data))
	}

	return widget.NewForm(items...)
}

func createBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}
