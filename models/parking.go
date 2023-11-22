package models

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Park struct {
	place      chan int    //canal
	door       *sync.Mutex //Puerta
	placeArray [20]bool
}

func NewPark(place chan int, doorMu *sync.Mutex) *Park {
	return &Park{
		place:      place,
		door:       doorMu,     //Datos
		placeArray: [20]bool{}, //20 lugares
	}
}

func (p *Park) Getplace() chan int {
	return p.place //Retorna canal
}

func (p *Park) GetdoorMu() *sync.Mutex {
	return p.door //Retorna el mutex de la puerta
}

func (p *Park) GetplaceArray() [20]bool {
	return p.placeArray
} //+Obtiene los datos

func (p *Park) SetplaceArray(placeArray [20]bool) {
	p.placeArray = placeArray
} //setea los datos

func (p *Park) WaitLeave(contenedor *fyne.Container, imagen *canvas.Image) {
	imagen.Move(fyne.NewPos(80, 280))
	contenedor.Add(imagen) //Metodo para la espera para la salida
	contenedor.Refresh()
}
