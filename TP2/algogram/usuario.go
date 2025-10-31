package algogram

import (
	"tp2/tdas/cola_prioridad"
	"math"
)

// entradaFeed guarda la info necesaria para priorizar posts en el feed
type entradaFeed struct {
	post      *post
	prioridad int 
	id_post   int
}

type usuario struct {
	nombre    string
	posicion  int 
	feed      cola_prioridad.ColaPrioridad[entradaFeed]
}

// Func de comparacion
func cmpFeed(a, b entradaFeed) int {
	if a.prioridad != b.prioridad {
		return b.prioridad - a.prioridad
	}
	return b.id_post - a.id_post
}

func nuevoUsuario(nombre string, pos int) *usuario {
	return &usuario{
		nombre:   nombre,
		posicion: pos,
		feed:     cola_prioridad.CrearHeap[entradaFeed](cmpFeed),
	}
}

// agregarPostAlFeed agrega una entrada de post con la prioridad calculada
func (u *usuario) agregarPostAlFeed(p *post, prioridad int) {
	u.feed.Encolar(entradaFeed{
		post:      p,
		prioridad: prioridad,
		id_post:   p.id,
	})
}

// proximoPost nos devuelve el siguiente post del feed
func (u *usuario) proximoPost() *post {
	if u.feed.EstaVacia() {
		return nil
	}
	entrada := u.feed.Desencolar()
	return entrada.post
}

// calcularAfinidad devuelve la distancia entre este usuario y otro.
func (u *usuario) calcularAfinidad(otro *usuario) int {
	return int(math.Abs(float64(u.posicion - otro.posicion)))
}