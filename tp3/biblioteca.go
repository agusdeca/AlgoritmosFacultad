package biblioteca

import (
	"tp3/tdas/cola"
	"tp3/tdas/diccionario"
	"tp3/tdas/grafo"
)

const MAX_PASOS = 20

// CaminoMinimo busca la forma mas rapida para llegar entre origen y destino con un bfs

func CaminoMinimo[T comparable](g grafo.Grafo[T], origen, destino T) ([]T, int) {
	if !g.Existe(origen) || !g.Existe(destino) {
		return nil, -1
	}

	if origen == destino {
		return []T{origen}, 0
	}

	visitados := diccionario.CrearHash[T, bool]()
	padres := diccionario.CrearHash[T, T]()
	
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

	res:= reconstruirCamino(padres,origen,destino)
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

	distancias := diccionario.CrearHash[T, int]()
	
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

    visitados := diccionario.CrearHash[T, bool]()
    distancias := diccionario.CrearHash[T, int]()
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
func Lectura[T comparable](g grafo.Grafo[T], paginas[]T) []T{
	grados:= diccionario.CrearHash[T,int]()
	actuales:= diccionario.CrearHash[T,bool]()

	for _, p := range paginas{
		grados.Guardar(p,0)
		actuales.Guardar(p,true)
	}

	for _,p := range paginas{
		if !g.Existe(p){
			continue
		}

		adyacentes:= g.Adyacentes(p)

		for _,ady : range adyacentes{
			if actuales.Pertenece(ady){
				grados.Guardar(ady, grados.Obtener(ady)+1)
			}
		}
	}

	q := cola.CrearColaEnlazada[T]()

	for _,p:= range paginas{
		if grados.Obtener(p)==0{
			q.Encolar(p)
		}
	}

	resultado := make([]T, 0, len(paginas))

	for !q.EstaVacia(){
		actual:= q.Desencolar()
		resultado= append(resultado,actual)

		adyacentes:= grafo.adyacentes(actual)
		for _,ady:= range adyacentes{
			if actuales.Pertenece(ady){
				grados.Guardar(ady, grados.Obtener(ady)-1)
				if grados.Obtener(ady) == 0{
					q.Encolar(ady)
				} 
			}
		}
	}

	if len(resultado) != len(paginas) {
        return nil
    }

	//Ahora lo invierto asi veo la inversa de los links
	for i, j := 0, len(resultado)-1; i < j; i, j = i+1, j-1 {
        resultado[i], resultado[j] = resultado[j], resultado[i]
    }
    
    return resultado
}


// Ciclo busca un ciclo simple de largo n 
func Ciclo[T comparable](g grafo.Grafo[T], origen T, n int) []T {
    if !g.Existe(origen) {
        return nil
    }
    
    visitados := diccionario.CrearHash[T, bool]()
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



//Ya hechas: ciclo n articulos, Camino mas corto, Lectura, diametro, todos en rango, navegacion