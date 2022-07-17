package main

import (
	"image/color"

	"github.com/the1234321/o2fy/src/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/the1234321/parsecode/src/pso2s"
)

func main() {
	o2fy_app := app.New()
	o2fy_app.Settings().SetTheme(&theme.O2fyTheme{})

	o2fy_main_window := o2fy_app.NewWindow("O2FY")

	// o2fy_main_window.SetContent(widget.NewLabel("PSO2牛逼助手!"))

	label1 := canvas.NewText("Label 1", color.Black)
	value1 := canvas.NewText("Value", color.White)
	label2 := canvas.NewText("Label 2", color.Black)
	value2 := canvas.NewText("Something", color.White)
	grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2)
	o2fy_main_window.SetContent(grid)

	o2fy_main_window.Resize(fyne.NewSize(800, 600))

	go func() { pso2s.StartParse() }()

	o2fy_main_window.ShowAndRun()
}
