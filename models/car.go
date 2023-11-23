package models

import (
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

// Car representa un automóvil en el estacionamiento.
type Car struct {
	id              int
	tiempoLim       time.Duration
	espacioAsignado int
	imageEnter      *canvas.Image
	imageWait       *canvas.Image
	imageLeave      *canvas.Image
}

// NewCar crea un nuevo Car con un id dado.
// Inicializa las imágenes del auto en diferentes estados (entrar, esperar, salir).
func NewCar(id int) *Car {
	return &Car{
		id:         id,
		tiempoLim:  generateRandomDuration(),
		imageEnter: loadImage("./assets/auto_entrada.png"),
		imageWait:  loadImage("./assets/auto_espera.png"),
		imageLeave: loadImage("./assets/auto_salida.png"),
	}
}

// Enter maneja la entrada de un auto al estacionamiento.
// Obtiene un lugar y lo asigna al auto si está disponible.
func (a *Car) Enter(p *Park, contenedor *fyne.Container) {
	p.Getplace() <- a.GetId()
	p.GetdoorMu().Lock()
	a.checkAndAssignPlace(p, contenedor)
	p.GetdoorMu().Unlock()
	contenedor.Refresh()
}

// Leave maneja la salida de un auto del estacionamiento.
// Actualiza el estado del lugar ocupado y prepara el auto para la salida.
func (a *Car) Leave(p *Park, contenedor *fyne.Container) {
	<-p.Getplace()
	p.GetdoorMu().Lock()
	a.updateParkingSpace(p)
	p.GetdoorMu().Unlock()
	a.prepareForExit(contenedor)
}

// Start inicia el proceso de estacionamiento del auto.
// Avanza al auto, entra al estacionamiento, espera y luego sale.
func (a *Car) Start(p *Park, contenedor *fyne.Container, wg *sync.WaitGroup) {
	a.Advance(9)
	a.Enter(p, contenedor)
	time.Sleep(a.tiempoLim)
	a.exitProcedure(p, contenedor, wg)
}

// Advance mueve el auto hacia adelante un número dado de pasos.
func (a *Car) Advance(steps int) {
	for i := 0; i < steps; i++ {
		a.moveCarForward()
		time.Sleep(time.Millisecond * 200)
	}
}

// GetId devuelve el id del auto.
func (a *Car) GetId() int {
	return a.id
}

// GetimageEnter devuelve la imagen del auto al entrar.
func (a *Car) GetimageEnter() *canvas.Image {
	return a.imageEnter
}

// loadImage carga una imagen desde una ruta de archivo dada.
func loadImage(path string) *canvas.Image {
	return canvas.NewImageFromURI(storage.NewFileURI(path))
}

// generateRandomDuration genera una duración aleatoria para el tiempo de estacionamiento.
func generateRandomDuration() time.Duration {
	return time.Duration(rand.Intn(40)+5) * time.Second
}

// checkAndAssignPlace verifica la disponibilidad de lugares y asigna uno al auto.
func (a *Car) checkAndAssignPlace(p *Park, contenedor *fyne.Container) {
	place := p.GetplaceArray()
	for i, isTaken := range place {
		if !isTaken {
			place[i] = true
			a.espacioAsignado = i
			a.imageEnter.Move(fyne.NewPos(float32(650-(i*30)), 330))
			break
		}
	}
	p.SetplaceArray(place)
}

// updateParkingSpace actualiza el estado de los espacios de estacionamiento cuando un auto sale.
func (a *Car) updateParkingSpace(p *Park) {
	spacesArray := p.GetplaceArray()
	spacesArray[a.espacioAsignado] = false
	p.SetplaceArray(spacesArray)
}

// prepareForExit prepara el auto para la salida del estacionamiento.
func (a *Car) prepareForExit(contenedor *fyne.Container) {
	contenedor.Remove(a.imageWait)
	a.imageLeave.Resize(fyne.NewSize(30, 50))
	a.imageLeave.Move(fyne.NewPos(90, 290))
	contenedor.Add(a.imageLeave)
	contenedor.Refresh()
	for i := 0; i < 10; i++ {
		a.imageLeave.Move(fyne.NewPos(a.imageLeave.Position().X, a.imageLeave.Position().Y-30))
		time.Sleep(time.Millisecond * 200)
	}
	contenedor.Remove(a.imageLeave)
	contenedor.Refresh()
}

// exitProcedure maneja el proceso de salida del auto del estacionamiento.
func (a *Car) exitProcedure(p *Park, contenedor *fyne.Container, wg *sync.WaitGroup) {
	contenedor.Remove(a.imageEnter)
	a.imageWait.Resize(fyne.NewSize(50, 30))
	p.WaitLeave(contenedor, a.imageWait)
	a.Leave(p, contenedor)
	wg.Done()
}

// moveCarForward avanza el auto hacia adelante.
func (a *Car) moveCarForward() {
	a.imageEnter.Move(fyne.NewPos(a.imageEnter.Position().X, a.imageEnter.Position().Y+15))
}
