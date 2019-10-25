package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

//JSON es una función
func JSON(interfaz interface{}) string {
	bytesJSON, _ := json.Marshal(interfaz)
	return string(bytesJSON)
}

//JSONIdentado es una función
func JSONIdentado(interfaz interface{}) string {
	bytesJSON, _ := json.MarshalIndent(interfaz, "", "    ")
	return string(bytesJSON)
}

//Imprimir es una función
func Imprimir(interfaz interface{}) {
	fmt.Println(JSON(interfaz))
}

//ImprimirIdentado es una función
func ImprimirIdentado(interfaz interface{}) {
	fmt.Println(JSONIdentado(interfaz))
}

type Toro struct {
	Bidimensional *Bidimensional
}

func CrearToro(anchura, altura int, valorPorDefecto interface{}) *Toro {
	var toro Toro
	toro.Bidimensional = CrearBidimensional(anchura, altura, valorPorDefecto)
	return &toro
}

func (toro *Toro) Limitar(valor, inferior, superior int) int {
	diferencia := superior - inferior
	for valor < inferior {
		valor += diferencia
	}
	for valor >= superior {
		valor -= diferencia
	}
	return valor
}

func (toro *Toro) ConvertirIndice(fila, columna int) (int, int) {
	fila = toro.Limitar(fila, 0, toro.LeerAltura())
	columna = toro.Limitar(columna, 0, toro.LeerAnchura())
	return fila, columna
}

func (toro *Toro) Leer(fila, columna int) interface{} {
	fila, columna = toro.ConvertirIndice(fila, columna)
	return toro.Bidimensional.Leer(fila, columna)
}

func (toro *Toro) Escribir(fila, columna int, valor interface{}) {
	fila, columna = toro.ConvertirIndice(fila, columna)
	toro.Bidimensional.Escribir(fila, columna, valor)
}

func (toro *Toro) LeerAltura() int {
	return toro.Bidimensional.LeerAltura()
}

func (toro *Toro) LeerAnchura() int {
	return toro.Bidimensional.LeerAnchura()
}

//Bidimensional es una estructura
type Bidimensional struct {
	Interno [][]interface{}
	Anchura int
	Altura  int
}

//CrearBidimensional es una función que crea un nuevo Bidimensional
func CrearBidimensional(anchura, altura int, valorPorDefecto interface{}) *Bidimensional {
	var arreglo Bidimensional
	arreglo.Anchura = anchura
	arreglo.Altura = altura
	arreglo.ConstruirInterno()
	arreglo.Llenar(valorPorDefecto)
	return &arreglo
}

//ConstruirInterno es una función
func (arreglo *Bidimensional) ConstruirInterno() {
	arreglo.Interno = make([][]interface{}, arreglo.Altura)
	for fila := 0; fila < arreglo.Altura; fila++ {
		arreglo.Interno[fila] = make([]interface{}, arreglo.Anchura)
	}
}

//Llenar es una función
func (arreglo *Bidimensional) Llenar(valor interface{}) {
	for fila := 0; fila < arreglo.Altura; fila++ {
		for columna := 0; columna < arreglo.Anchura; columna++ {
			arreglo.Interno[fila][columna] = valor
		}
	}
}

//Leer es una función
func (arreglo *Bidimensional) Leer(fila, columna int) interface{} {
	return arreglo.Interno[fila][columna]
}

//Escribir es una función
func (arreglo *Bidimensional) Escribir(fila, columna int, valor interface{}) {
	arreglo.Interno[fila][columna] = valor
}

func (arreglo *Bidimensional) LeerAltura() int {
	return arreglo.Altura
}

func (arreglo *Bidimensional) LeerAnchura() int {
	return arreglo.Anchura
}

func (arreglo *Bidimensional) LlenarDesdeDatos(datos [][]interface{}) {
	for fila:= 0; fila < arreglo.Altura; fila++ {
		for columna:= 0; columna < arreglo.Anchura; columna++ {
			arreglo.Interno[fila][columna] = datos[fila][columna]
		}
	}
}

type Cuarto struct {
	Ocupantes *Toro
}

func CrearCuarto(anchura, altura int) *Cuarto {
	var cuarto Cuarto
	cuarto.Ocupantes = CrearToro(anchura, altura, nil)
	return &cuarto
}

func (cuarto *Cuarto) AñadirOcupante(ocupante Ocupante) {
	cuarto.Ocupantes.Escribir(ocupante.LeerFila(), ocupante.LeerColumna(), ocupante)
}

func (cuarto *Cuarto) LeerAltura() int {
	return cuarto.Ocupantes.LeerAltura()
}

func (cuarto *Cuarto) LeerAnchura() int {
	return cuarto.Ocupantes.LeerAnchura()
}

type Posición struct {
	Fila int
	Columna int
}

func CrearPosición(fila, columna int) *Posición {
	var posición Posición
	posición.Columna = columna
	posición.Fila = fila
	return &posición
}

type Símbolo struct {
	Celdas *Bidimensional
}

func CrearSímbolo(datos [][]byte) *Símbolo {
	altura := len(datos)
	anchura := len(datos[0])
	var símbolo Símbolo
	símbolo.Celdas = CrearBidimensional(anchura, altura, 0)
	interfaz := make([][]interface{}, len(datos))
	for fila := 0; fila < len(datos); fila++ {
		interfaz[fila] = make([]interface{}, len(datos[fila]))
		for columna := 0; columna < len(datos[fila]); columna++ {
			interfaz[fila][columna] = datos[fila][columna]
		}
	}
	símbolo.Celdas.LlenarDesdeDatos(interfaz)
	return &símbolo
}

//Jugador es una estructura
type Jugador struct {
	Posición *Posición
	Símbolo *Símbolo
}

func CrearJugador(posición *Posición, símbolo *Símbolo) *Jugador {
	var jugador Jugador
	jugador.Posición = posición
	jugador.Símbolo = símbolo
	return &jugador
}

func (jugador *Jugador) LeerFila() int {
	return jugador.Posición.Fila
}

func (jugador *Jugador) LeerColumna() int {
	return jugador.Posición.Columna
}

func (jugador *Jugador) LeerSímbolo() *Símbolo {
	return jugador.Símbolo
}

type Ocupante interface {
	LeerFila() int
	LeerColumna() int
	LeerSímbolo() *Símbolo
}

type Píxel struct {
	Color Color
}

type Color struct {
	Rojo byte
	Verde byte
	Azul byte
}

type Implementación struct {
	Juego *Juego
	Ventana *pixelgl.Window
	Imagen *image.RGBA
}

func CrearImplementación(juego *Juego) *Implementación{
	var implementación Implementación
	implementación.Juego = juego
	return &implementación
}

func (implementación *Implementación) Correr() {
	pixelgl.Run(implementación.HiloPrincipal)
}

func (implementación *Implementación) HiloPrincipal() {
	configuración := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 512),
		VSync:  true,
	}
	ventana, err := pixelgl.NewWindow(configuración)
	if err != nil {
		panic(err)
	}

	imagen := image.NewRGBA(image.Rect(0, 0, 1024, 512))
	for y := 0; y < imagen.Rect.Dy(); y++ {
		for x := 0; x < imagen.Rect.Dx(); x++ {
			color := color.RGBA{R: 255, G: 0, B: 0, A: 255}
			imagen.SetRGBA(x, y, color)
		}
	}

	cuadro := pixel.PictureDataFromImage(imagen)
	sprite := pixel.NewSprite(cuadro, cuadro.Bounds())

	ventana.Clear(colornames.Skyblue)
	sprite.Draw(ventana, pixel.IM.Moved(ventana.Bounds().Center()))

	for !ventana.Closed() {
		ventana.Update()
	}
}

type Juego struct {
	Cuarto *Cuarto
	Jugador *Jugador
	Pixeles *Bidimensional
	TamañoDeSímbolo int
	SímboloVacío *Símbolo
	Colores []Color
	Implementación *Implementación
}

//CrearDatosDelJuego es una función
func CrearJuego() *Juego {
	var juego Juego
	símboloDelJugador := CrearSímbolo([][]byte{
		{0, 0, 0, 1, 0, 0, 0, 1},
		{0, 1, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0, 0},
		{1, 0, 0, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 1},
		{0, 0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 1, 0},
		{1, 0, 0, 0, 1, 0, 0, 0},
	})
	jugador := CrearJugador(CrearPosición(0, 0), símboloDelJugador)
	cuarto := CrearCuarto(128, 64)
	cuarto.AñadirOcupante(jugador)
	juego.Cuarto = cuarto
	juego.Jugador = jugador
	juego.TamañoDeSímbolo = 8
	juego.SímboloVacío = CrearSímbolo([][]byte{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	})
	juego.Implementación = CrearImplementación(&juego)
	juego.Colores = append(juego.Colores, Color{Rojo: 0, Verde: 0, Azul: 0})
	juego.Colores = append(juego.Colores, Color{Rojo: 255, Verde: 255, Azul: 255})
	juego.Pixeles = CrearBidimensional(1024, 512, Píxel{Color: Color{Rojo: 0, Verde: 0, Azul: 0}})
	return &juego
}

func (juego *Juego) DibujarSímbolo(símbolo *Símbolo, fila, columna int) {
	yInicial := fila * juego.TamañoDeSímbolo
	xInicial := columna * juego.TamañoDeSímbolo
	yFinal := yInicial + juego.TamañoDeSímbolo
	xFinal := xInicial + juego.TamañoDeSímbolo
	for y := yInicial; y < yFinal; y++ {
		for x := xInicial; x < xFinal; x++ {
			colorIndex := símbolo.Celdas.Leer(y, x).(int)

			píxel.Color = juego.Colores[colorIndex]
			juego.Pixeles.Escribir(y, x, Píxel{Color: })
		}
	}

}

func (juego *Juego) DibujarOcupante(fila, columna int) {
	var símbolo *Símbolo
	ocupante := juego.Cuarto.Ocupantes.Leer(fila, columna).(Ocupante)
	if ocupante == nil {
		símbolo = juego.SímboloVacío
	} else {
		símbolo = ocupante.LeerSímbolo()
	}
	juego.DibujarSímbolo(símbolo, fila, columna)
}

func (juego *Juego) Dibujar() {
	for fila := 0; fila < juego.Cuarto.LeerAltura(); fila++ {
		for columna := 0; columna < juego.Cuarto.LeerAnchura(); columna++ {
			juego.DibujarOcupante(fila, columna)
		}
	}
}

func (juego *Juego) Jugar() {
	juego.Implementación.Correr()
}

func main() {
	juego := CrearJuego()
	juego.Jugar()
	//arreglo := CrearBidimensional(4, 2, nil)
	//ImprimirIdentado(arreglo)
}


