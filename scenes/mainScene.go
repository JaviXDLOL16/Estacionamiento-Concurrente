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

// MainScene representa la escena principal de la aplicación.
type MainScene struct {
	window fyne.Window
}

// NewMainScene crea una nueva instancia de MainScene.
func NewMainScene(window fyne.Window) *MainScene {
	return &MainScene{window: window}
}

var contenedor = container.NewWithoutLayout()

// Show configura y muestra los componentes de la escena.
func (s *MainScene) Show() {
	loadBackground()
	createParkingBorder()
	createParkingSpaces()
	createEntrance()
	s.window.SetContent(contenedor)
}

// loadBackground carga y muestra la imagen de fondo de la escena.
func loadBackground() {
	background := canvas.NewImageFromFile("./assets/estacionamiento_suelo.png")
	background.FillMode = canvas.ImageFillStretch
	background.Resize(fyne.NewSize(700, 400))
	contenedor.Add(background)
}

// createParkingBorder crea y añade el borde del área de estacionamiento.
func createParkingBorder() {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 10
	border.StrokeColor = color.Transparent
	border.Resize(fyne.NewSize(690, 140))
	border.Move(fyne.NewPos(0, 250))
	contenedor.Add(border)
}

// createParkingSpaces crea y añade las líneas que delimitan los espacios de estacionamiento.
func createParkingSpaces() {
	for i := 0; i < 20; i++ {
		spaceOutline := createSpaceOutline(i)
		contenedor.Add(spaceOutline)
	}
}

// createSpaceOutline crea un rectángulo que representa un espacio de estacionamiento.
func createSpaceOutline(index int) *canvas.Rectangle {
	spaceOutline := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 255, A: 255})
	spaceOutline.StrokeWidth = 2
	spaceOutline.StrokeColor = color.RGBA{R: 231, G: 231, B: 10, A: 255}
	spaceOutline.FillColor = color.RGBA{R: 108, G: 107, B: 112, A: 1}
	spaceOutline.Resize(fyne.NewSize(25, 50))
	spaceOutline.Move(fyne.NewPos(83+float32(index*30), 330))
	return spaceOutline
}

// createEntrance crea y añade la representación gráfica de la entrada del estacionamiento.
func createEntrance() {
	entrance := canvas.NewRectangle(color.RGBA{R: 231, G: 231, B: 10, A: 231})
	entrance.Resize(fyne.NewSize(100, 10))
	entrance.Move(fyne.NewPos(30, 245))
	contenedor.Add(entrance)
}

// Run inicia la lógica de la escena principal, creando autos y gestionando su comportamiento.
func (s *MainScene) Run() {
	p := models.NewPark(make(chan int, 20), &sync.Mutex{})
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go generateCars(p, &wg, i)
		waitRandomTime()
	}

	wg.Wait()
}

// generateCars genera un auto y lo añade a la escena.
func generateCars(park *models.Park, wg *sync.WaitGroup, id int) {
	car := models.NewCar(id)
	image := car.GetimageEnter()
	image.Resize(fyne.NewSize(30, 50))
	image.Move(fyne.NewPos(40, -10))

	contenedor.Add(image)
	contenedor.Refresh()

	car.Start(park, contenedor, wg)
}

// waitRandomTime espera un tiempo aleatorio entre la generación de autos.
func waitRandomTime() {
	time.Sleep(time.Second * time.Duration(generatePoisson(2)))
}

// generarPoisson utiliza una distribución de Poisson para generar un valor aleatorio.
func generatePoisson(lambda float64) float64 {
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	return poisson.Rand()
}
