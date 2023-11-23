package models

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Park representa un estacionamiento con un número limitado de plazas.
type Park struct {
	place      chan int    // Canal para manejar la disponibilidad de plazas.
	door       *sync.Mutex // Mutex para controlar el acceso a la puerta del estacionamiento.
	placeArray [20]bool    // Array que representa la disponibilidad de cada plaza en el estacionamiento.
}

// NewPark crea e inicializa una nueva instancia de Park.
func NewPark(place chan int, doorMu *sync.Mutex) *Park {
	return &Park{
		place:      place,
		door:       doorMu,
		placeArray: [20]bool{},
	}
}

// Getplace devuelve el canal utilizado para la gestión de las plazas del estacionamiento.
func (p *Park) Getplace() chan int {
	return p.place
}

// GetdoorMu devuelve el mutex utilizado para controlar el acceso a la puerta del estacionamiento.
func (p *Park) GetdoorMu() *sync.Mutex {
	return p.door
}

// GetplaceArray devuelve el array que representa la disponibilidad de las plazas del estacionamiento.
func (p *Park) GetplaceArray() [20]bool {
	return p.placeArray
}

// SetplaceArray establece el estado actualizado de las plazas del estacionamiento.
func (p *Park) SetplaceArray(placeArray [20]bool) {
	p.placeArray = placeArray
}

// WaitLeave gestiona la espera para la salida de un auto del estacionamiento.
func (p *Park) WaitLeave(contenedor *fyne.Container, imagen *canvas.Image) {
	prepareExit(imagen, contenedor)
}

// prepareExit prepara y ejecuta la animación para la salida de un auto.
func prepareExit(imagen *canvas.Image, contenedor *fyne.Container) {
	imagen.Move(fyne.NewPos(80, 280))
	contenedor.Add(imagen)
	contenedor.Refresh()
}
