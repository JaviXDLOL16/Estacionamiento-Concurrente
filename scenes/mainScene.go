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
	return &MainScene{
		window: window,
	}
}

var contenedor = container.NewWithoutLayout()

func (s *MainScene) Show() {
	contorno := canvas.NewRectangle(color.Transparent)
	contorno.StrokeWidth = 3
	contorno.StrokeColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	contorno.Resize(fyne.NewSize(690, 140))
	contorno.Move(fyne.NewPos(0, 250))

	puerta := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
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
