package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	AnchoVentana         = 800
	AltoVentana          = 600
	AnchoEstacionamiento = 600
	AltoEstacionamiento  = 400
	AnchoEntrada         = 100
	NumEspacios          = 20
	AnchoEspacio         = AnchoEstacionamiento / (NumEspacios / 2)
	AltoEspacio          = AltoEstacionamiento / (NumEspacios / 2)
	AnchoLinea           = 2
	AnchoCoche           = AnchoEspacio - 10
	AltoCoche            = AltoEspacio - AnchoLinea - 10
	VelocidadCarro       = 2.0
)

type Direccion int

const (
	Entrando Direccion = iota
	Saliendo
)

type Carro struct {
	x, y, targetX, targetY float64
	tiempoRestante         int
	direccion              Direccion
}

type EntradaSalida struct {
	canal  chan Direccion
	actual Direccion
}

func (e *EntradaSalida) SolicitarAcceso(d Direccion) {
	if d == e.actual {
		return
	}

	e.canal <- d
	e.actual = d
}

func (e *EntradaSalida) LiberarAcceso() {
	<-e.canal
}

type Game struct {
	carrosEstacionados   []Carro
	carrosEsperando      []Carro
	tiempoParaNuevoCarro int
	entradaSalida        EntradaSalida
}

func (g *Game) Update() error {
	g.tiempoParaNuevoCarro--
	if g.tiempoParaNuevoCarro <= 0 {
		g.tiempoParaNuevoCarro = rand.Intn(60) + 30
		g.carrosEsperando = append(g.carrosEsperando, Carro{x: (AnchoVentana - AnchoEntrada) / 2, y: 0, targetX: -AnchoCoche, targetY: -AltoCoche, direccion: Entrando})
	}

	// Actualizar los carros estacionados
	for i, carro := range g.carrosEstacionados {
		// Si el carro todavía no ha llegado a su destino
		if carro.x != carro.targetX || carro.y != carro.targetY {
			if carro.x < carro.targetX {
				carro.x += VelocidadCarro
			} else if carro.x > carro.targetX {
				carro.x -= VelocidadCarro
			}

			if carro.y < carro.targetY {
				carro.y += VelocidadCarro
			} else if carro.y > carro.targetY {
				carro.y -= VelocidadCarro
			}

			g.carrosEstacionados[i] = carro

		} else {
			// Si el carro ha llegado a su destino
			if carro.direccion == Saliendo {
				// Si el carro está saliendo, liberar el espacio y el acceso
				g.entradaSalida.LiberarAcceso()
				g.carrosEstacionados = append(g.carrosEstacionados[:i], g.carrosEstacionados[i+1:]...)
				i--
			} else if carro.tiempoRestante <= 0 {
				// Si el carro ha terminado su tiempo, cambiar su dirección a Saliendo
				carro.direccion = Saliendo
				carro.targetX = (AnchoVentana + AnchoEntrada) / 2
				carro.targetY = 0
				g.entradaSalida.SolicitarAcceso(Saliendo)
				g.carrosEstacionados[i] = carro
			} else {
				carro.tiempoRestante--
				g.carrosEstacionados[i] = carro
			}
		}
	}

	// Procesar carros esperando entrar
	for _, carro := range g.carrosEsperando {
		g.entradaSalida.SolicitarAcceso(Entrando)
		// Encontrar un espacio vacío para el carro
		for i := 0; i < 10; i++ {
			x := float64((AnchoVentana-AnchoEstacionamiento)/2) + 5
			y := float64((AltoVentana-AltoEstacionamiento)/2) + float64(i)*AltoEspacio + 5
			if !g.carroEnPosicion(x, y) {
				carro.targetX = x
				carro.targetY = y
				carro.tiempoRestante = 360
				g.carrosEstacionados = append(g.carrosEstacionados, carro)
				break
			}

			x = float64((AnchoVentana+AnchoEstacionamiento)/2) - AnchoEspacio + 5
			if !g.carroEnPosicion(x, y) {
				carro.targetX = x
				carro.targetY = y
				carro.tiempoRestante = 360
				g.carrosEstacionados = append(g.carrosEstacionados, carro)
				break
			}
		}
	}

	g.carrosEsperando = nil

	return nil
}

func (g *Game) carroEnPosicion(x, y float64) bool {
	for _, carro := range g.carrosEstacionados {
		if carro.targetX == x && carro.targetY == y {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Dibujar fondo del estacionamiento
	rectEstacionamiento := ebiten.NewImage(AnchoEstacionamiento, AltoEstacionamiento)
	rectEstacionamiento.Fill(color.RGBA{128, 128, 128, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((AnchoVentana-AnchoEstacionamiento)/2), float64((AltoVentana-AltoEstacionamiento)/2))
	screen.DrawImage(rectEstacionamiento, op)

	// Dibujar espacios de estacionamiento y líneas de separación
	for i := 0; i < 10; i++ {
		rectLineaIzq := ebiten.NewImage(AnchoLinea, AltoEspacio)
		rectLineaIzq.Fill(color.White)
		opIzq := &ebiten.DrawImageOptions{}
		opIzq.GeoM.Translate(float64((AnchoVentana-AnchoEstacionamiento)/2), float64((AltoVentana-AltoEstacionamiento)/2)+float64(i)*AltoEspacio)
		screen.DrawImage(rectLineaIzq, opIzq)

		rectLineaDer := ebiten.NewImage(AnchoLinea, AltoEspacio)
		rectLineaDer.Fill(color.White)
		opDer := &ebiten.DrawImageOptions{}
		opDer.GeoM.Translate(float64((AnchoVentana+AnchoEstacionamiento)/2)-AnchoLinea, float64((AltoVentana-AltoEstacionamiento)/2)+float64(i)*AltoEspacio)
		screen.DrawImage(rectLineaDer, opDer)

		rectEspacioIzq := ebiten.NewImage(int(AnchoCoche), int(AltoCoche))
		rectEspacioIzq.Fill(color.RGBA{64, 64, 64, 255})
		opIzq = &ebiten.DrawImageOptions{}
		opIzq.GeoM.Translate(float64((AnchoVentana-AnchoEstacionamiento)/2)+5, float64((AltoVentana-AltoEstacionamiento)/2)+float64(i)*AltoEspacio+5)
		screen.DrawImage(rectEspacioIzq, opIzq)

		rectEspacioDer := ebiten.NewImage(int(AnchoCoche), int(AltoCoche))
		rectEspacioDer.Fill(color.RGBA{64, 64, 64, 255})
		opDer = &ebiten.DrawImageOptions{}
		opDer.GeoM.Translate(float64((AnchoVentana+AnchoEstacionamiento)/2)-AnchoEspacio+5, float64((AltoVentana-AltoEstacionamiento)/2)+float64(i)*AltoEspacio+5)
		screen.DrawImage(rectEspacioDer, opDer)
	}

	// Dibujar coches estacionados y en espera
	for _, carro := range g.carrosEstacionados {
		rectCarro := ebiten.NewImage(int(AnchoCoche), int(AltoCoche))
		rectCarro.Fill(color.RGBA{0, 128, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(carro.x, carro.y)
		screen.DrawImage(rectCarro, op)
	}

	for _, carro := range g.carrosEsperando {
		rectCarro := ebiten.NewImage(int(AnchoCoche), int(AltoCoche))
		rectCarro.Fill(color.RGBA{255, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(carro.x, carro.y)
		screen.DrawImage(rectCarro, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return AnchoVentana, AltoVentana
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(AnchoVentana, AltoVentana)
	ebiten.SetWindowTitle("Estacionamiento")
	game := &Game{
		entradaSalida: EntradaSalida{canal: make(chan Direccion, 1)},
	}
	ebiten.RunGame(game)
}
