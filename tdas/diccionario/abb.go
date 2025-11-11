package diccionario

import (
	TDAPila "tdas/pila"
)

// funcion de comparacion
type funcCmp[K any] func(K, K) int

type nodoAbb[K any, V any] struct {
	izquierdo *nodoAbb[K, V]
	derecho   *nodoAbb[K, V]
	clave     K
	dato      V
}

type abb[K any, V any] struct {
	raiz     *nodoAbb[K, V]
	cantidad int
	cmp      funcCmp[K]
}

type iterAbb[K any, V any] struct {
	pila  TDAPila.Pila[*nodoAbb[K, V]]
	cmp   func(K, K) int
	desde *K
	hasta *K
}

// CrearABB crea un nuevo ABB vacío con la func de cmp
func CrearABB[K any, V any](cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: cmp}
}

// Guardar inserta o reemplaza una clave
func (a *abb[K, V]) Guardar(clave K, dato V) {
	padre, nodo := a.buscarNodo(a.raiz, clave)

	if nodo != nil {
		// Caso 1: La clave ya existe. Actualizamos el dato.
		nodo.dato = dato
		return
	}

	// Caso 2: La clave no existe. Creamos un nuevo nodo.
	nuevoNodo := &nodoAbb[K, V]{clave: clave, dato: dato}
	a.cantidad++

	if padre == nil {
		a.raiz = nuevoNodo
	} else {
		comp := a.cmp(clave, padre.clave)
		if comp < 0 {
			padre.izquierdo = nuevoNodo
		} else {
			padre.derecho = nuevoNodo
		}
	}
}

// Pertenece indica si una clave está en el ABB
func (a *abb[K, V]) Pertenece(clave K) bool {
	_, nodo := a.buscarNodo(a.raiz, clave)
	return nodo != nil
}

// Obtener devuelve el valor asociado a una clave
func (a *abb[K, V]) Obtener(clave K) V {
	_, nodo := a.buscarNodo(a.raiz, clave)
	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	return nodo.dato
}

func (a *abb[K, V]) buscarNodo(nodo *nodoAbb[K, V], clave K) (*nodoAbb[K, V], *nodoAbb[K, V]) {
	var padre *nodoAbb[K, V]
	actual := nodo

	for actual != nil {
		comp := a.cmp(clave, actual.clave)
		if comp == 0 {
			return padre, actual
		}
		padre = actual
		if comp < 0 {
			actual = actual.izquierdo
		} else {
			actual = actual.derecho
		}
	}
	return padre, nil
}

func (a *abb[K, V]) buscarNodoAux(padre *nodoAbb[K, V], nodo *nodoAbb[K, V], clave K) (*nodoAbb[K, V], *nodoAbb[K, V]) {
	if nodo == nil {
		return padre, nil
	}

	comp := a.cmp(clave, nodo.clave)
	if comp == 0 {
		return padre, nodo
	}

	if comp < 0 {
		return a.buscarNodoAux(nodo, nodo.izquierdo, clave)
	}
	return a.buscarNodoAux(nodo, nodo.derecho, clave)
}

// Borrar elimina una clave del ABB y devuelve su valor
func (a *abb[K, V]) Borrar(clave K) V {
	padre, nodo := a.buscarNodo(a.raiz, clave)

	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}

	dato := nodo.dato
	a.cantidad--

	if nodo.izquierdo == nil {
		// Caso 1: 0 hijos o 1 hijo (derecho)
		a.reemplazarHijo(padre, nodo, nodo.derecho)
	} else if nodo.derecho == nil {
		// Caso 2: 1 hijo (izquierdo)
		a.reemplazarHijo(padre, nodo, nodo.izquierdo)
	} else {
		// Caso 3: 2 hijos
		padreSucesor, sucesor := a.buscarMinimoConPadre(nodo.derecho, nodo)

		nodo.clave, nodo.dato = sucesor.clave, sucesor.dato

		a.reemplazarHijo(padreSucesor, sucesor, sucesor.derecho)
	}
	return dato
}

func (a *abb[K, V]) reemplazarHijo(padre *nodoAbb[K, V], nodo *nodoAbb[K, V], reemplazo *nodoAbb[K, V]) {
	if padre == nil {
		// 'nodo' era la raiz
		a.raiz = reemplazo
	} else if padre.izquierdo == nodo {
		padre.izquierdo = reemplazo
	} else {
		padre.derecho = reemplazo
	}
}

func (a *abb[K, V]) buscarMinimoConPadre(nodo *nodoAbb[K, V], padre *nodoAbb[K, V]) (*nodoAbb[K, V], *nodoAbb[K, V]) {
	actual := nodo
	for actual.izquierdo != nil {
		padre = actual
		actual = actual.izquierdo
	}
	return padre, actual
}

// Cantidad devuelve el num de elem
func (a *abb[K, V]) Cantidad() int {
	return a.cantidad
}

// Iterador interno
func (a *abb[K, V]) Iterar(visitar func(K, V) bool) {
	iterarRango(a.raiz, a.cmp, nil, nil, visitar)
}

func (a *abb[K, V]) IterarRango(desde, hasta *K, visitar func(K, V) bool) {
	iterarRango(a.raiz, a.cmp, desde, hasta, visitar)
}

// Func aux que recorre el arbol recursivamente y ve solo las claves dentro del rango
func iterarRango[K any, V any](nodo *nodoAbb[K, V], cmp funcCmp[K], desde, hasta *K, visitar func(K, V) bool) bool {
	if nodo == nil {
		return true
	}
	if desde == nil || cmp(nodo.clave, *desde) >= 0 {
		if !iterarRango(nodo.izquierdo, cmp, desde, hasta, visitar) {
			return false
		}
	}
	if (desde == nil || cmp(nodo.clave, *desde) >= 0) && (hasta == nil || cmp(nodo.clave, *hasta) <= 0) {
		if !visitar(nodo.clave, nodo.dato) {
			return false
		}
	}
	if hasta == nil || cmp(nodo.clave, *hasta) <= 0 {
		if !iterarRango(nodo.derecho, cmp, desde, hasta, visitar) {
			return false
		}
	}
	return true
}

// Iterador externo
func (a *abb[K, V]) Iterador() IterDiccionario[K, V] {
	return a.IteradorRango(nil, nil)
}
func (a *abb[K, V]) IteradorRango(desde, hasta *K) IterDiccionario[K, V] {
	iter := &iterAbb[K, V]{pila: TDAPila.CrearPilaDinamica[*nodoAbb[K, V]](), cmp: a.cmp, desde: desde, hasta: hasta}
	iter.apilarIzquierdos(a.raiz, desde)
	return iter
}
func (iter *iterAbb[K, V]) apilarIzquierdos(nodo *nodoAbb[K, V], desde *K) {
	if desde == nil {
		for nodo != nil {
			iter.pila.Apilar(nodo)
			nodo = nodo.izquierdo
		}
		return
	}
	for nodo != nil {
		comp := iter.cmp(nodo.clave, *desde)
		if comp >= 0 {
			iter.pila.Apilar(nodo)
			nodo = nodo.izquierdo
		} else {
			nodo = nodo.derecho
		}
	}
}
func (iter *iterAbb[K, V]) HaySiguiente() bool {
	if iter.pila.EstaVacia() {
		return false
	}
	tope := iter.pila.VerTope()
	if iter.hasta != nil && iter.cmp(tope.clave, *iter.hasta) > 0 {
		return false
	}
	return true
}
func (iter *iterAbb[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	nodo := iter.pila.VerTope()
	return nodo.clave, nodo.dato
}

// Siguiente avanza al siguiente elemento en el recorrido inorder
func (iter *iterAbb[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic(MENSAJE_ITER_TERMINADO)
	}
	nodo := iter.pila.Desapilar()
	if nodo.derecho != nil {
		iter.apilarIzquierdos(nodo.derecho, iter.desde)
	}
}
