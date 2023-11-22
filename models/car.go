package models

import (
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

func NewCar(id int) *Car {
	imageEnter := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_entrada.png"))
	imageWait := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_espera.png")) //Mandas a llamar las iamgenes con sus tres modos
	imageLeave := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_salida.png"))
	return &Car{
		id:              id,
		tiempoLim:       time.Duration(rand.Intn(40)+5) * time.Second,
		espacioAsignado: 0,
		imageEnter:      imageEnter,
		imageWait:       imageWait, //Creas el objeto car (Son los carritos)
		imageLeave:      imageLeave,
	}
}

type Car struct {
	id              int
	tiempoLim       time.Duration
	espacioAsignado int
	imageEnter      *canvas.Image
	imageWait       *canvas.Image //Se le asignan sus respectivos valores
	imageLeave      *canvas.Image
}

func (a *Car) Enter(p *Park, contenedor *fyne.Container) {
	p.Getplace() <- a.GetId()
	p.GetdoorMu().Lock() //Bloquea
	//Verifica si hay espacios disponibles
	place := p.GetplaceArray()

	a.Advance(5) //Avanza si hay lugar

	for i := 0; i < len(place); i++ {
		if place[i] == false {
			place[i] = true //Da un lugar y le asigna un numero si hay espacio
			a.espacioAsignado = i
			a.imageEnter.Move(fyne.NewPos(float32(650-(i*30)), 330))
			break
		}
	}
	p.SetplaceArray(place)

	p.GetdoorMu().Unlock() //Desbloquea la puerta
	contenedor.Refresh()
}

func (a *Car) Leave(p *Park, contenedor *fyne.Container) {
	<-p.Getplace()
	p.GetdoorMu().Lock() //Bloquea la puerta para no dejar entrar

	spacesArray := p.GetplaceArray()
	spacesArray[a.espacioAsignado] = false //Actualiza el arreglo si hay lugar libre
	p.SetplaceArray(spacesArray)

	p.GetdoorMu().Unlock() //desbloquea la puerta

	contenedor.Remove(a.imageWait)
	a.imageLeave.Resize(fyne.NewSize(30, 50)) //Refresca el lugar del auto
	a.imageLeave.Move(fyne.NewPos(90, 290))

	contenedor.Add(a.imageLeave)
	contenedor.Refresh() //Cambia la trancision

	for i := 0; i < 10; i++ {
		a.imageLeave.Move(fyne.NewPos(a.imageLeave.Position().X, a.imageLeave.Position().Y-30))
		time.Sleep(time.Millisecond * 200)
	}

	contenedor.Remove(a.imageLeave) //Remueve la imagen salir
	contenedor.Refresh()
}

func (a *Car) Start(p *Park, contenedor *fyne.Container, wg *sync.WaitGroup) {
	a.Advance(18)

	a.Enter(p, contenedor)

	time.Sleep(a.tiempoLim)

	contenedor.Remove(a.imageEnter)          //remueve la imagen
	a.imageWait.Resize(fyne.NewSize(50, 30)) //Pone el modo espera
	p.WaitLeave(contenedor, a.imageWait)
	a.Leave(p, contenedor)
	wg.Done()
}

func (a *Car) Advance(pasos int) {
	for i := 0; i < pasos; i++ {
		a.imageEnter.Move(fyne.NewPos(a.imageEnter.Position().X, a.imageEnter.Position().Y+8)) //Funcion para avanzar
		time.Sleep(time.Millisecond * 200)
	}
}

func (a *Car) GetId() int {
	return a.id //Regresa el id
}

func (a *Car) GetimageEnter() *canvas.Image {
	return a.imageEnter //Da imagen de entrar
}
