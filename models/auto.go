package models

import (
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

type Auto struct {
	id              int
	tiempoLim       time.Duration
	espacioAsignado int
	imagenEntrada   *canvas.Image
	imagenEspera    *canvas.Image
	imagenSalida    *canvas.Image
}

func NewAuto(id int) *Auto {
	imagenEntrada := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_entrada.png"))
	imagenEspera := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_espera.png"))
	imagenSalida := canvas.NewImageFromURI(storage.NewFileURI("./assets/auto_salida.png"))
	return &Auto{
		id:              id,
		tiempoLim:       time.Duration(rand.Intn(40)+5) * time.Second,
		espacioAsignado: 0,
		imagenEntrada:   imagenEntrada,
		imagenEspera:    imagenEspera,
		imagenSalida:    imagenSalida,
	}
}

func (a *Auto) Entrar(p *Estacionamiento, contenedor *fyne.Container) {
	p.GetEspacios() <- a.GetId()
	p.GetPuertaMu().Lock()
	espacios := p.GetEspaciosArray()

	anchoEspacio := 30
	inicioEspacios := (700 - float32(20*anchoEspacio)) / 2

	for i := 0; i < len(espacios); i++ {
		if !espacios[i] {
			espacios[i] = true
			a.espacioAsignado = i
			a.imagenEntrada.Move(fyne.NewPos(inicioEspacios+float32(i*anchoEspacio), 350))
			break
		}
	}
	p.SetEspaciosArray(espacios)
	p.GetPuertaMu().Unlock()
	contenedor.Refresh()
}

func (a *Auto) Salir(p *Estacionamiento, contenedor *fyne.Container) {
	<-p.GetEspacios()
	p.GetPuertaMu().Lock()

	spacesArray := p.GetEspaciosArray()
	spacesArray[a.espacioAsignado] = false
	p.SetEspaciosArray(spacesArray)

	p.GetPuertaMu().Unlock()

	contenedor.Remove(a.imagenEntrada)
	a.imagenEspera.Resize(fyne.NewSize(30, 50))
	a.imagenEspera.Move(fyne.NewPos(40+float32(a.espacioAsignado*35), 350))
	contenedor.Add(a.imagenEspera)
	contenedor.Refresh()

	time.Sleep(a.tiempoLim)

	contenedor.Remove(a.imagenEspera)
	contenedor.Refresh()

	p.GetEspacios() <- a.GetId()
}

func (a *Auto) Iniciar(p *Estacionamiento, contenedor *fyne.Container, wg *sync.WaitGroup) {
	defer wg.Done()
	a.Entrar(p, contenedor)
	time.Sleep(a.tiempoLim)
	a.Salir(p, contenedor)
}

func (a *Auto) Avanzar(pasos int) {
	for i := 0; i < pasos; i++ {
		a.imagenEntrada.Move(fyne.NewPos(a.imagenEntrada.Position().X, a.imagenEntrada.Position().Y+20))
		time.Sleep(time.Millisecond * 200)
	}
}

func (a *Auto) GetId() int {
	return a.id
}

func (a *Auto) GetImagenEntrada() *canvas.Image {
	return a.imagenEntrada
}
