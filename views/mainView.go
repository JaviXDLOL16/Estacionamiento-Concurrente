package views

import (
	"estacionamiento/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type MainView struct{} //Aqui se definen los datos de la ventana

func NewMainView() *MainView {
	return &MainView{}
}

func (v *MainView) Run() {
	myApp := app.New()
	window := myApp.NewWindow("Parking") //Nombre de la ventana
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(700, 400)) //Tamaño

	mainScene := scenes.NewMainScene(window)
	mainScene.Show()
	go mainScene.Run()
	window.ShowAndRun()
}
