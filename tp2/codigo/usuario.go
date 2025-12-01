package codigo

import (
	"math"
	"tdas/cola_prioridad"
)

type Usuario interface {
	Nombre() string
	Posicion() int
	RecibirPost(Post)
	ProximoPost() Post
}

type entradaFeed struct {
	post      Post
	prioridad int
	id_post   int
}

type usuario struct {
	nombre   string
	posicion int
	feed     cola_prioridad.ColaPrioridad[entradaFeed]
}

func cmpFeed(a, b entradaFeed) int {
	if a.prioridad != b.prioridad {
		return b.prioridad - a.prioridad
	}
	return b.id_post - a.id_post
}

func nuevoUsuario(nombre string, pos int) Usuario {
	return &usuario{
		nombre:   nombre,
		posicion: pos,
		feed:     cola_prioridad.CrearHeap[entradaFeed](cmpFeed),
	}
}

func (u *usuario) Nombre() string {
	return u.nombre
}

func (u *usuario) Posicion() int {
	return u.posicion
}

func (u *usuario) calcularAfinidad(otro *usuario) int {
	return int(math.Abs(float64(u.posicion - otro.posicion)))
}

func (u *usuario) RecibirPost(p Post) {
	otro := p.(*post).autor
	afinidad := u.calcularAfinidad(otro)
	u.feed.Encolar(entradaFeed{
		post:      p,
		prioridad: afinidad,
		id_post:   p.ID(),
	})
}

func (u *usuario) ProximoPost() Post {
	if u.feed.EstaVacia() {
		return nil
	}
	return u.feed.Desencolar().post
}
