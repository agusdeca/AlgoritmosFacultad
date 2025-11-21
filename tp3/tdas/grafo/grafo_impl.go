package grafo

import (
	"tp3/tdas/diccionario"
	"tp3/tdas/lista"
)

// grafoNoPesado implementa un grafo sin pesos usando diccionario de listas
type grafoNoPesado[T comparable] struct {
	dirigido   bool
	vertices   diccionario.Diccionario[T, bool]
	adyacencia diccionario.Diccionario[T, diccionario.Diccionario[T, bool]]
}

// arista representa una arista con peso para el grafo pesado
type arista[T comparable, P any] struct {
	destino T
	peso    P
}

// grafoPesado implementa un grafo con pesos usando diccionario de listas
type grafoPesado[T comparable, P any] struct {
	dirigido   bool
	vertices   diccionario.Diccionario[T, bool]
	adyacencia diccionario.Diccionario[T, diccionario.Diccionario[T, P]]
}

// CrearGrafoNoPesado crea un nuevo grafo no pesado
func CrearGrafoNoPesado[T comparable](dirigido bool) GrafoNoPesado[T] {
	return &grafoNoPesado[T]{
		dirigido:   dirigido,
		vertices:   diccionario.CrearHash[T, bool](),
		adyacencia: diccionario.CrearHash[T, diccionario.Diccionario[T, bool]](),
	}
}

// CrearGrafoPesado crea un nuevo grafo pesado
func CrearGrafoPesado[T comparable, P any](dirigido bool) GrafoPesado[T, P] {
	return &grafoPesado[T, P]{dirigido:dirigido, vertices:diccionario.CrearHash[T, bool](), adyacencia: diccionario.CrearHash[T, diccionario.Diccionario[T, P]](),
	}
}

// Grafo no pesado

func (g *grafoNoPesado[T]) EsDirigido() bool {
	return g.dirigido
}

func (g *grafoNoPesado[T]) AgregarVertice(vertice T) {
	if !g.Existe(vertice) {
		g.vertices.Guardar(vertice, true)
		g.adyacencia.Guardar(vertice, diccionario.CrearHash[T, bool]())
	}
}

func (g *grafoNoPesado[T]) BorrarVertice(vertice T) {
	if !g.Existe(vertice) {
		return
	}

	// Elimino las aristas que llegan a este vertice
	iter := g.vertices.Iterador()
	for iter.HaySiguiente() {
		v, _ := iter.VerActual()
		adyacentes := g.adyacencia.Obtener(v)
		if adyacentes.Pertenece(vertice) {
			adyacentes.Borrar(vertice)
		}
		iter.Siguiente()
	}

	// Elimino el vertice
	g.vertices.Borrar(vertice)
	g.adyacencia.Borrar(vertice)
}


func (g *grafoNoPesado[T]) AgregarArista(v1 T, v2 T) {
	if !g.Existe(v1) {
		g.AgregarVertice(v1)
	}
	if !g.Existe(v2) {
		g.AgregarVertice(v2)
	}

	adyacentes := g.adyacencia.Obtener(v1)
	adyacentes.Guardar(v2, true)

	if !g.dirigido {
		adyacentes2 := g.adyacencia.Obtener(v2)
		adyacentes2.Guardar(v1, true)
	}
}

func (g *grafoNoPesado[T]) BorrarArista(v1 T, v2 T) {
	if !g.Existe(v1) || !g.Existe(v2) {
		return
	}

	adyacentes := g.adyacencia.Obtener(v1)
	if adyacentes.Pertenece(v2) {
		adyacentes.Borrar(v2)
	}

	if !g.dirigido {
		adyacentes2 := g.adyacencia.Obtener(v2)
		if adyacentes2.Pertenece(v1) {
			adyacentes2.Borrar(v1)
		}
	}
}

func (g *grafoNoPesado[T]) HayArista(v1 T, v2 T) bool {
	if !g.Existe(v1) || !g.Existe(v2) {
		return false
	}
	adyacentes := g.adyacencia.Obtener(v1)
	return adyacentes.Pertenece(v2)
}

func (g *grafoNoPesado[T]) Existe(vertice T) bool {
	return g.vertices.Pertenece(vertice)
}

func (g *grafoNoPesado[T]) ObtenerVertices() []T {
	vertices := make([]T, 0, g.vertices.Cantidad())
	iter := g.vertices.Iterador()
	for iter.HaySiguiente() {
		v, _ := iter.VerActual()
		vertices = append(vertices, v)
		iter.Siguiente()
	}
	return vertices
}

func (g *grafoNoPesado[T]) Cantidad() int {
	return g.vertices.Cantidad()
}

func (g *grafoNoPesado[T]) Adyacentes(v T) []T {
	if !g.Existe(v) {
		return []T{}
	}

	adyacentes := g.adyacencia.Obtener(v)
	resultado := make([]T, 0, adyacentes.Cantidad())
	iter := adyacentes.Iterador()
	for iter.HaySiguiente() {
		ady, _ := iter.VerActual()
		resultado = append(resultado, ady)
		iter.Siguiente()
	}
	return resultado
}

// Grafo Pesado

func (g *grafoPesado[T, P]) EsDirigido() bool {
	return g.dirigido
}

func (g *grafoPesado[T, P]) AgregarVertice(vertice T) {
	if !g.Existe(vertice) {
		g.vertices.Guardar(vertice, true)
		g.adyacencia.Guardar(vertice, diccionario.CrearHash[T, P]())
	}
}

func (g *grafoPesado[T, P]) BorrarVertice(vertice T) {
	if !g.Existe(vertice) {
		return
	}

	//Elimino todas las aristas que lleguen al vertice que va a ser borrado
	iter := g.vertices.Iterador()
	for iter.HaySiguiente() {
		v, _ := iter.VerActual()
		adyacentes := g.adyacencia.Obtener(v)
		if adyacentes.Pertenece(vertice) {
			adyacentes.Borrar(vertice)
		}
		iter.Siguiente()
	}

	//Ahora si lo elimino
	g.vertices.Borrar(vertice)
	g.adyacencia.Borrar(vertice)
}

func (g *grafoPesado[T, P]) AgregarArista(v1 T, v2 T, peso P) {
	if !g.Existe(v1) {
		g.AgregarVertice(v1)
	}
	if !g.Existe(v2) {
		g.AgregarVertice(v2)
	}

	adyacentes := g.adyacencia.Obtener(v1)
	adyacentes.Guardar(v2, peso)

	if !g.dirigido {
		adyacentes2 := g.adyacencia.Obtener(v2)
		adyacentes2.Guardar(v1, peso)
	}
}

func (g *grafoPesado[T, P]) BorrarArista(v1 T, v2 T) {
	if !g.Existe(v1) || !g.Existe(v2) {
		return
	}

	adyacentes := g.adyacencia.Obtener(v1)
	if adyacentes.Pertenece(v2) {
		adyacentes.Borrar(v2)
	}

	if !g.dirigido {
		adyacentes2 := g.adyacencia.Obtener(v2)
		if adyacentes2.Pertenece(v1) {
			adyacentes2.Borrar(v1)
		}
	}
}

func (g *grafoPesado[T, P]) HayArista(v1 T, v2 T) bool {
	if !g.Existe(v1) || !g.Existe(v2) {
		return false
	}
	adyacentes := g.adyacencia.Obtener(v1)
	return adyacentes.Pertenece(v2)
}

func (g *grafoPesado[T, P]) Existe(vertice T) bool {
	return g.vertices.Pertenece(vertice)
}

func (g *grafoPesado[T, P]) ObtenerVertices() []T {
	vertices := make([]T, 0, g.vertices.Cantidad())
	iter := g.vertices.Iterador()
	for iter.HaySiguiente() {
		v, _ := iter.VerActual()
		vertices = append(vertices, v)
		iter.Siguiente()
	}
	return vertices
}

func (g *grafoPesado[T, P]) Cantidad() int {
	return g.vertices.Cantidad()
}

func (g *grafoPesado[T, P]) Adyacentes(v T) []T {
	if !g.Existe(v) {
		return []T{}
	}

	adyacentes := g.adyacencia.Obtener(v)
	resultado := make([]T, 0, adyacentes.Cantidad())
	iter := adyacentes.Iterador()
	for iter.HaySiguiente() {
		ady, _ := iter.VerActual()
		resultado = append(resultado, ady)
		iter.Siguiente()
	}
	return resultado
}

func (g *grafoPesado[T, P]) VerPeso(v1, v2 T) P {
	if !g.Existe(v1) || !g.Existe(v2) {
		var cero P
		return cero
	}
	adyacentes := g.adyacencia.Obtener(v1)
	if !adyacentes.Pertenece(v2) {
		var cero P
		return cero
	}
	return adyacentes.Obtener(v2)
}

