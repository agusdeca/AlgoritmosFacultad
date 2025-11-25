package tp3

import (
	"sort"
	"tp3/tdas/cola"
	"tp3/tdas/diccionario"
	"tp3/tdas/grafo"
)

const (
	MAX_PASOS  = 20
	D_PAGERANK = 0.85
	ITER_PR    = 40
)

type ParPR struct {
	pagina string
	valor  float64
}

// CaminoMinimo busca la forma mas rapida para llegar entre origen y destino con un bfs

func CaminoMinimo[T comparable](g grafo.Grafo[T], origen, destino T) ([]T, int) {
	if !g.Existe(origen) || !g.Existe(destino) {
		return nil, -1
	}

	if origen == destino {
		return []T{origen}, 0
	}

	visitados := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })
	padres := diccionario.CrearHash[T, T](func(a, b T) bool { return a == b })

	q := cola.CrearColaEnlazada[T]()

	q.Encolar(origen)
	visitados.Guardar(origen, true)

	encontrado := false

	for !q.EstaVacia() {
		actual := q.Desencolar()

		if actual == destino {
			encontrado = true
			break
		}

		adyacentes := g.Adyacentes(actual)
		for _, ady := range adyacentes {
			if !visitados.Pertenece(ady) {
				visitados.Guardar(ady, true)
				padres.Guardar(ady, actual)
				q.Encolar(ady)
			}
		}
	}

	if !encontrado {
		return nil, -1
	}

	res := reconstruirCamino(padres, origen, destino)
	return res, len(res) - 1
}

// Funcion aux para construirme el camino de origen a destino
func reconstruirCamino[T comparable](padres diccionario.Diccionario[T, T], origen, destino T) []T {
	camino := []T{}
	actual := destino

	for actual != origen {
		camino = append(camino, actual)
		actual = padres.Obtener(actual)
	}
	camino = append(camino, origen)

	// Invertir
	for i, j := 0, len(camino)-1; i < j; i, j = i+1, j-1 {
		camino[i], camino[j] = camino[j], camino[i]
	}

	return camino
}

// EnRango cuenta la cantidad de vertices que están a exactamente n aristas

func EnRango[T comparable](g grafo.Grafo[T], origen T, n int) int {
	if !g.Existe(origen) {
		return 0
	}

	if n == 0 {
		return 1 // Es el origen
	}

	distancias := diccionario.CrearHash[T, int](func(a, b T) bool { return a == b })

	q := cola.CrearColaEnlazada[T]()

	q.Encolar(origen)
	distancias.Guardar(origen, 0)

	for !q.EstaVacia() {
		actual := q.Desencolar()
		distActual := distancias.Obtener(actual)

		if distActual >= n {
			continue
		}

		adyacentes := g.Adyacentes(actual)
		for _, ady := range adyacentes {
			if !distancias.Pertenece(ady) {
				distancias.Guardar(ady, distActual+1)
				q.Encolar(ady)
			}
		}
	}

	// Ahora contamos cuantos son los vertices que cumplen con la dist n
	contador := 0
	iter := distancias.Iterador()
	for iter.HaySiguiente() {
		_, dist := iter.VerActual()
		if dist == n {
			contador++
		}
		iter.Siguiente()
	}

	return contador
}

// Navegacion navega desde el origen siguiendo siempre el primer link

func Navegacion[T comparable](g grafo.Grafo[T], origen T) []T {
	if !g.Existe(origen) {
		return []T{}
	}

	camino := []T{origen}
	actual := origen

	for len(camino) < MAX_PASOS {
		adyacentes := g.Adyacentes(actual)

		// Si no tiene links corto
		if len(adyacentes) == 0 {
			break
		}

		siguiente := adyacentes[0]

		camino = append(camino, siguiente)
		actual = siguiente
	}

	return camino
}

// Diametro encuentra el camino más largo entre todos los caminos mínimos de la red
func Diametro[T comparable](g grafo.Grafo[T]) ([]T, int) {
	vertices := g.ObtenerVertices()
	if len(vertices) == 0 {
		return nil, 0
	}

	var maxDistanciaGlobal int = -1
	var origenMax, destinoMax T
	hayCamino := false

	// itero todos los vertices hasta encontrar el mayor diametro
	for _, origen := range vertices {
		distancia, destino := bfsMasLejano(g, origen)

		if distancia > maxDistanciaGlobal {
			maxDistanciaGlobal = distancia
			origenMax = origen
			destinoMax = destino
			hayCamino = true
		}
	}

	if !hayCamino {
		return nil, 0
	}

	return CaminoMinimo(g, origenMax, destinoMax)
}

func bfsMasLejano[T comparable](g grafo.Grafo[T], origen T) (int, T) {
	if !g.Existe(origen) {
		var nulo T
		return -1, nulo
	}

	visitados := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })
	distancias := diccionario.CrearHash[T, int](func(a, b T) bool { return a == b })
	q := cola.CrearColaEnlazada[T]()

	q.Encolar(origen)
	visitados.Guardar(origen, true)
	distancias.Guardar(origen, 0)

	var ultimoNodo T = origen
	var maxDist int = 0

	for !q.EstaVacia() {
		actual := q.Desencolar()
		distActual := distancias.Obtener(actual)
		if distActual > maxDist {
			maxDist = distActual
			ultimoNodo = actual
		}

		adyacentes := g.Adyacentes(actual)
		for _, ady := range adyacentes {
			if !visitados.Pertenece(ady) {
				visitados.Guardar(ady, true)
				distancias.Guardar(ady, distActual+1)
				q.Encolar(ady)
			}
		}
	}

	return maxDist, ultimoNodo
}

// Lectura devuelve un orden válido para leer las páginas dadas
func Lectura[T comparable](g grafo.Grafo[T], paginas []T) []T {
	grados := diccionario.CrearHash[T, int](func(a, b T) bool { return a == b })

	actuales := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })

	for _, p := range paginas {
		grados.Guardar(p, 0)
		actuales.Guardar(p, true)
	}

	for _, p := range paginas {
		if !g.Existe(p) {
			continue
		}
		adyacentes := g.Adyacentes(p)

		for _, ady := range adyacentes {
			if actuales.Pertenece(ady) {
				grados.Guardar(ady, grados.Obtener(ady)+1)
			}
		}
	}

	q := cola.CrearColaEnlazada[T]()
	for _, p := range paginas {
		if grados.Obtener(p) == 0 {
			q.Encolar(p)
		}
	}

	resultado := make([]T, 0, len(paginas))

	for !q.EstaVacia() {
		actual := q.Desencolar()
		resultado = append(resultado, actual)

		adyacentes := g.Adyacentes(actual)
		for _, ady := range adyacentes {
			if actuales.Pertenece(ady) {
				grados.Guardar(ady, grados.Obtener(ady)-1)
				if grados.Obtener(ady) == 0 {
					q.Encolar(ady)
				}
			}
		}
	}

	if len(resultado) != len(paginas) {
		return nil
	}

	return resultado
}

// Ciclo busca un ciclo simple de largo n
func Ciclo[T comparable](g grafo.Grafo[T], origen T, n int) []T {
	if !g.Existe(origen) {
		return nil
	}

	visitados := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })
	camino := []T{origen}
	visitados.Guardar(origen, true)

	solucion := DFS(g, origen, n, visitados, camino)

	return solucion
}

func DFS[T comparable](g grafo.Grafo[T], actual T, n int, visitados diccionario.Diccionario[T, bool], camino []T) []T {
	if len(camino) == n {
		origen := camino[0]
		if g.HayArista(actual, origen) {
			camino = append(camino, origen)
			return camino
		}
		return nil
	}

	adyacentes := g.Adyacentes(actual)
	for _, ady := range adyacentes {
		if visitados.Pertenece(ady) {
			continue
		}

		visitados.Guardar(ady, true)
		camino = append(camino, ady)

		res := DFS(g, ady, n, visitados, camino)
		if res != nil {
			return res
		}

		camino = camino[:len(camino)-1]
		visitados.Borrar(ady)
	}

	return nil
}

// MasImportantes COMPLEJIDAD: O(K(P+L) + P log P).
func MasImportantes(g grafo.Grafo[string], k int) []string {
	vertices := g.ObtenerVertices()
	n := len(vertices)
	if n == 0 {
		return []string{}
	}

	pr := diccionario.CrearHash[string, float64](func(a, b string) bool { return a == b })
	for _, v := range vertices {
		pr.Guardar(v, 1.0/float64(n))
	}

	for iter := 0; iter < ITER_PR; iter++ {
		nuevo := diccionario.CrearHash[string, float64](func(a, b string) bool { return a == b })
		base := (1 - D_PAGERANK) / float64(n)

		// inicializo con el valor base
		for _, v := range vertices {
			nuevo.Guardar(v, base)
		}

		for _, v := range vertices {
			ady := g.Adyacentes(v)
			if len(ady) == 0 {
				continue
			}
			contrib := pr.Obtener(v) * D_PAGERANK / float64(len(ady))
			for _, w := range ady {
				actual := nuevo.Obtener(w)
				nuevo.Guardar(w, actual+contrib)
			}
		}

		pr = nuevo
	}

	// Ordenar por PageRank
	aux := make([]ParPR, 0, n)
	for _, v := range vertices {
		aux = append(aux, ParPR{v, pr.Obtener(v)})
	}

	sort.Slice(aux, func(i, j int) bool {
		return aux[i].valor > aux[j].valor
	})

	if k > n {
		k = n
	}

	res := make([]string, k)
	for i := 0; i < k; i++ {
		res[i] = aux[i].pagina
	}
	return res
}

func Conectados[T comparable](g grafo.Grafo[T], inicio T) []T {
	if !g.Existe(inicio) {
		return []T{}
	}

	visitados := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })
	q := cola.CrearColaEnlazada[T]()

	res := []T{}

	q.Encolar(inicio)
	visitados.Guardar(inicio, true)

	for !q.EstaVacia() {
		v := q.Desencolar()
		res = append(res, v)

		for _, w := range g.Adyacentes(v) {
			if !visitados.Pertenece(w) {
				visitados.Guardar(w, true)
				q.Encolar(w)
			}
		}
	}

	return res
}

func Comunidad[T comparable](g grafo.Grafo[T], inicio T) []T {
	vertices := g.ObtenerVertices()
	if len(vertices) == 0 || !g.Existe(inicio) {
		return []T{}
	}

	etiqueta := diccionario.CrearHash[T, T](func(a, b T) bool { return a == b })
	for _, v := range vertices {
		etiqueta.Guardar(v, v)
	}

	for iter := 0; iter < 20; iter++ {
		huboCambio := false

		for _, v := range vertices {
			frecuencias := diccionario.CrearHash[T, int](func(a, b T) bool { return a == b })

			for _, vecino := range g.Adyacentes(v) {
				etiquetaVecino := etiqueta.Obtener(vecino)
				if frecuencias.Pertenece(etiquetaVecino) {
					frecuencias.Guardar(etiquetaVecino, frecuencias.Obtener(etiquetaVecino)+1)
				} else {
					frecuencias.Guardar(etiquetaVecino, 1)
				}
			}

			if frecuencias.Cantidad() == 0 {
				continue
			}

			// busco la etiqueta + frecuente
			mejorEtiqueta := etiqueta.Obtener(v)
			maxFrecuencia := -1

			iterFreq := frecuencias.Iterador()
			for iterFreq.HaySiguiente() {
				etq, cant := iterFreq.VerActual()
				if cant > maxFrecuencia {
					maxFrecuencia = cant
					mejorEtiqueta = etq
				}
				iterFreq.Siguiente()
			}

			if mejorEtiqueta != etiqueta.Obtener(v) {
				etiqueta.Guardar(v, mejorEtiqueta)
				huboCambio = true
			}
		}

		if !huboCambio {
			break
		}
	}

	etiquetaObjetivo := etiqueta.Obtener(inicio)

	comunidad := []T{}
	for _, v := range vertices {
		if etiqueta.Obtener(v) == etiquetaObjetivo {
			comunidad = append(comunidad, v)
		}
	}

	return comunidad
}

func Clustering[T comparable](g grafo.Grafo[T], v T) float64 {
	vecinos := g.Adyacentes(v)
	k := len(vecinos)
	if k < 2 {
		return 0
	}

	set := diccionario.CrearHash[T, bool](func(a, b T) bool { return a == b })
	for _, x := range vecinos {
		set.Guardar(x, true)
	}

	E := 0
	for _, x := range vecinos {
		for _, y := range g.Adyacentes(x) {
			if set.Pertenece(y) {
				E++
			}
		}
	}

	// Como contamos en grafo no dirigido, divido por 2
	return float64(E/2) / float64(k*(k-1)/2)
}
