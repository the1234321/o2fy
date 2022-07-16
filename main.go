package main

import (
	"o2fy/theme"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	o2fy_app := app.New()
	o2fy_app.Settings().SetTheme(&theme.O2fyTheme{})

	o2fy_main_window := o2fy_app.NewWindow("Hello World")

	o2fy_main_window.SetContent(widget.NewLabel("PSO2牛逼助手!"))
	o2fy_main_window.Show()

	// w2 := o2fy_app.NewWindow("Larger")
	// w2.Resize(fyne.NewSize(100, 100))
	// w2.SetContent(widget.NewButton("Open new", func() {
	// 	w3 := o2fy_app.NewWindow("Third")
	// 	w3.SetContent(widget.NewLabel("Third"))
	// 	w3.Show()
	// }))
	// w2.Show()

	o2fy_app.Run()
}
