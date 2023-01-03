package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	content1 := widget.NewLabel("text1")
	content2 := widget.NewLabel("text2")
	content3 := widget.NewLabel("text3")

	content := container.New(layout.NewHBoxLayout(), content1, content2, layout.NewSpacer(), content3)

	w.SetContent(content)
	w.ShowAndRun()
}
