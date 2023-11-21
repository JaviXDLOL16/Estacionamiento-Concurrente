package scenes

import (
	"estacionamiento/models"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"gonum.org/v1/gonum/stat/distuv"
)

type MainScene struct {
	window fyne.Window
}

func NewMainScene(window fyne.Window) *MainScene {
	return &MainScene{window: window}
}

var contenedor = container.NewWithoutLayout()

func (s *MainScene) Show() {
	// Cargar imagen de fondo
	fondo := canvas.NewImageFromFile("./assets/estacionamiento_suelo.png")
	fondo.FillMode = canvas.ImageFillStretch
	fondo.Resize(fyne.NewSize(700, 400))
	contenedor.Add(fondo)

	contorno := canvas.NewRectangle(color.Transparent)
	contorno.StrokeWidth = 10
	contorno.StrokeColor = color.Transparent //borde estacionamiento
	contorno.Resize(fyne.NewSize(690, 140))
	contorno.Move(fyne.NewPos(0, 250))

	// Delimitar los espacios de estacionamiento
	for i := 0; i < 20; i++ {
		spaceOutline := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
		spaceOutline.StrokeWidth = 2
		spaceOutline.StrokeColor = color.RGBA{R: 231, G: 231, B: 10, A: 255} // borde contenedor
		spaceOutline.FillColor = color.RGBA{R: 108, G: 107, B: 112, A: 1}    // dentro contenedor
		spaceOutline.Resize(fyne.NewSize(25, 50))
		spaceOutline.Move(fyne.NewPos(83+float32(i*30), 330))
		contenedor.Add(spaceOutline)
	}

	puerta := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 255, A: 255}) //Entrada puerta
	puerta.Resize(fyne.NewSize(100, 10))
	puerta.Move(fyne.NewPos(30, 245))

	contenedor.Add(contorno)
	contenedor.Add(puerta)
	s.window.SetContent(contenedor)
}

func (s *MainScene) Run() {
	p := models.NewEstacionamiento(make(chan int, 20), &sync.Mutex{})
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			auto := models.NewAuto(id)
			imagen := auto.GetImagenEntrada()
			imagen.Resize(fyne.NewSize(30, 50))
			imagen.Move(fyne.NewPos(40, -10))

			contenedor.Add(imagen)
			contenedor.Refresh()

			auto.Iniciar(p, contenedor, &wg)
		}(i)
		var poisson = generarPoisson(float64(2))
		time.Sleep(time.Second * time.Duration(poisson))
	}

	wg.Wait()
}

func generarPoisson(lambda float64) float64 {
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	return poisson.Rand()
}
