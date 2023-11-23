package views

import (
	"estacionamiento/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// MainView representa la vista principal de la aplicación.
type MainView struct{}

// NewMainView crea una nueva instancia de MainView.
func NewMainView() *MainView {
	return &MainView{}
}

// Run inicializa y ejecuta la vista principal.
func (v *MainView) Run() {
	window := setupWindow()
	mainScene := scenes.NewMainScene(window)
	displayMainScene(window, mainScene)
}

// setupWindow configura y devuelve una nueva ventana de Fyne.
func setupWindow() fyne.Window {
	myApp := app.New()
	window := myApp.NewWindow("Parking")
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(700, 400))
	return window
}

// displayMainScene muestra la escena principal en la ventana.
func displayMainScene(window fyne.Window, mainScene *scenes.MainScene) {
	mainScene.Show()
	go mainScene.Run()
	window.ShowAndRun()
}
