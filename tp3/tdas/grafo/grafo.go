package grafo

type Grafo[T comparable] interface {
	// EsDirigido devuelve true si el grafo es dirigido
	EsDirigido() bool

	// AgregarVertice agrega un nuevo vértice al grafo
	AgregarVertice(vertice T)

	// BorrarVertice elimina un vértice y todas sus aristas asociadas
	BorrarVertice(vertice T)

	// BorrarArista elimina la arista entre v1 y v2
	BorrarArista(v1 T, v2 T)

	// HayArista verifica si existe una arista de v1 a v2
	HayArista(v1 T, v2 T) bool

	// Existe verifica si un vértice existe en el grafo
	Existe(vertice T) bool

	// ObtenerVertices devuelve todos los vértices del grafo
	ObtenerVertices() []T

	// Cantidad devuelve la cantidad de vértices en el grafo
	Cantidad() int

	// Adyacentes devuelve los vértices adyacentes a v
	Adyacentes(v T) []T
}

// GrafoNoPesado representa un grafo sin pesos en las aristas
type GrafoNoPesado[T comparable] interface {
	Grafo[T]

	// AgregarArista agrega una arista entre v1 y v2
	AgregarArista(v1 T, v2 T)
}

// GrafoPesado representa un grafo con pesos en las aristas
type GrafoPesado[T comparable, P any] interface {
	Grafo[T]

	// AgregarArista agrega una arista entre v1 y v2 con un peso
	AgregarArista(v1 T, v2 T, peso P)

	// VerPeso devuelve el peso de la arista entre v1 y v2
	VerPeso(v1, v2 T) P
}
