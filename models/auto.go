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
		tiempoLim:       time.Duration(rand.Intn(50)+10) * time.Second,
		espacioAsignado: 0,
		imagenEntrada:   imagenEntrada,
		imagenEspera:    imagenEspera,
		imagenSalida:    imagenSalida,
	}
}

func (a *Auto) Iniciar(p *Estacionamiento, contenedor *fyne.Container, wg *sync.WaitGroup) {
	defer wg.Done()

	// Esperar por un espacio disponible y luego entrar
	<-p.GetEspacios()
	a.Entrar(p, contenedor)

	// Simular tiempo estacionado
	time.Sleep(a.tiempoLim)

	// Iniciar el proceso de salida
	a.Salir(p, contenedor)
}

func (a *Auto) Entrar(p *Estacionamiento, contenedor *fyne.Container) {
	p.GetPuertaMu().Lock()
	espacios := p.GetEspaciosArray()

	for i := 0; i < len(espacios); i++ {
		if !espacios[i] {
			espacios[i] = true
			a.espacioAsignado = i
			a.imagenEntrada.Move(fyne.NewPos(float32(650-(i*30)), 330))
			break
		}
	}
	p.SetEspaciosArray(espacios)
	p.GetPuertaMu().Unlock()
	contenedor.Refresh()

	// Aquí debe avanzar hacia su espacio asignado
	a.AvanzarHaciaEspacio(contenedor)
}

func (a *Auto) AvanzarHaciaEspacio(contenedor *fyne.Container) {
	for a.imagenEntrada.Position().Y < 290 { // Asumiendo que 290 es la posición Y final
		a.imagenEntrada.Move(fyne.NewPos(a.imagenEntrada.Position().X, a.imagenEntrada.Position().Y+10))
		time.Sleep(time.Millisecond * 200)
		contenedor.Refresh()
	}
}

func (a *Auto) Salir(p *Estacionamiento, contenedor *fyne.Container) {
	p.GetPuertaMu().Lock()
	espacios := p.GetEspaciosArray()
	espacios[a.espacioAsignado] = false
	p.SetEspaciosArray(espacios)
	p.GetPuertaMu().Unlock()

	contenedor.Remove(a.imagenEntrada)
	a.imagenEspera.Resize(fyne.NewSize(50, 30))
	a.imagenEspera.Move(fyne.NewPos(float32(650-(a.espacioAsignado*30)), 330))
	contenedor.Add(a.imagenEspera)
	contenedor.Refresh()

	// Simular tiempo antes de salir
	time.Sleep(2 * time.Second)

	contenedor.Remove(a.imagenEspera)
	contenedor.Refresh()

	// Indicar que el espacio está nuevamente disponible
	p.GetEspacios() <- a.GetId()
}

func (a *Auto) GetId() int {
	return a.id
}

func (a *Auto) GetImagenEntrada() *canvas.Image {
	return a.imagenEntrada
}
