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
	contorno := canvas.NewRectangle(color.Transparent)
	contorno.StrokeWidth = 3
	contorno.StrokeColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	contorno.Resize(fyne.NewSize(690, 140))
	contorno.Move(fyne.NewPos(5, 260))

	puerta := canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	puerta.Resize(fyne.NewSize(100, 10))
	puerta.Move(fyne.NewPos(30, 245))

	// Dibujar rectángulos alrededor de los espacios de estacionamiento
	espacioInicialX := float32(35)  // Posición inicial X para el primer espacio
	espacioInicialY := float32(300) // Posición inicial Y para los espacios
	for i := 0; i < 20; i++ {
		spaceOutline := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 255, A: 255})
		spaceOutline.StrokeWidth = 2
		spaceOutline.StrokeColor = color.White
		spaceOutline.FillColor = color.RGBA{R: 150, G: 150, B: 255, A: 50}
		spaceOutline.Resize(fyne.NewSize(30, 50))
		spaceOutline.Move(fyne.NewPos(espacioInicialX+float32(i*35), espacioInicialY)) // Ajusta la posición según sea necesario
		contenedor.Add(spaceOutline)
	}

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
