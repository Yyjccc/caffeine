package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("边框布局")

	top := canvas.NewText("顶部栏", color.White)
	left := canvas.NewText("左侧", color.White)
	middle := canvas.NewText("内容", color.White)
	content := container.NewBorder(top, nil, left, nil, middle)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
