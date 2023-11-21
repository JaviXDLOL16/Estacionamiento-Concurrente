package views

import (
	"estacionamiento/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type MainView struct{}

func NewMainView() *MainView {
	return &MainView{}
}

func (v *MainView) Run() {
	myApp := app.New()
	window := myApp.NewWindow("Estacionamiento concurrentes")
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(700, 400))

	mainScene := scenes.NewMainScene(window)
	mainScene.Show()
	go mainScene.Run()
	window.ShowAndRun()
}
