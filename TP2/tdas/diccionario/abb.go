package diccionario

import (
	TDAPila "tp2/tdas/pila"
)

// funcion de comparacion
type funcCmp[K comparable] func(K, K) int

type nodoAbb[K comparable, V any] struct {
	izquierdo *nodoAbb[K, V]
	derecho   *nodoAbb[K, V]
	clave     K
	dato      V
}

type abb[K comparable, V any] struct {
	raiz     *nodoAbb[K, V]
	cantidad int
	cmp      funcCmp[K]
}

type iterAbb[K comparable, V any] struct {
	pila  TDAPila.Pila[*nodoAbb[K, V]]
	cmp   func(K, K) int
	desde *K
	hasta *K
}

// CrearABB crea un nuevo ABB vacío con la func de cmp
func CrearABB[K comparable, V any](cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: cmp}
}

// Guardar inserta o reemplaza una clave
func (a *abb[K, V]) Guardar(clave K, dato V) {
	padre, nodo := buscarNodo(a.raiz, clave, a.cmp)

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
	_, nodo := buscarNodo(a.raiz, clave, a.cmp)
	return nodo != nil
}

// Obtener devuelve el valor asociado a una clave
func (a *abb[K, V]) Obtener(clave K) V {
	_, nodo := buscarNodo(a.raiz, clave, a.cmp)
	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	return nodo.dato
}

func buscarNodo[K comparable, V any](raiz *nodoAbb[K, V], clave K, cmp funcCmp[K]) (padre *nodoAbb[K, V], nodo *nodoAbb[K, V]) {
	padre = nil
	nodo = raiz

	for nodo != nil {
		comp := cmp(clave, nodo.clave)
		if comp == 0 {
			return padre, nodo
		}
		padre = nodo
		if comp < 0 {
			nodo = nodo.izquierdo
		} else {
			nodo = nodo.derecho
		}
	}
	return padre, nil
}

// Borrar elimina una clave del ABB y devuelve su valor
func (a *abb[K, V]) Borrar(clave K) V {
	padre, nodo := buscarNodo(a.raiz, clave, a.cmp)

	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}

	dato := nodo.dato
	a.cantidad--

	if nodo.izquierdo == nil {
		// Caso 1: 0 hijos o 1 hijo (derecho)
		reemplazarHijo(a, padre, nodo, nodo.derecho)
	} else if nodo.derecho == nil {
		// Caso 2: 1 hijo (izquierdo)
		reemplazarHijo(a, padre, nodo, nodo.izquierdo)
	} else {
		// Caso 3: 2 hijos
		padreSucesor, sucesor := buscarMinimoConPadre(nodo.derecho, nodo)

		nodo.clave, nodo.dato = sucesor.clave, sucesor.dato

		reemplazarHijo(a, padreSucesor, sucesor, sucesor.derecho)
	}
	return dato
}

func reemplazarHijo[K comparable, V any](a *abb[K, V], padre *nodoAbb[K, V], nodo *nodoAbb[K, V], reemplazo *nodoAbb[K, V]) {
	if padre == nil {
		// 'nodo' era la raiz
		a.raiz = reemplazo
	} else if padre.izquierdo == nodo {
		padre.izquierdo = reemplazo
	} else {
		padre.derecho = reemplazo
	}
}

func buscarMinimoConPadre[K comparable, V any](nodo *nodoAbb[K, V], padre *nodoAbb[K, V]) (*nodoAbb[K, V], *nodoAbb[K, V]) {
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
func iterarRango[K comparable, V any](nodo *nodoAbb[K, V], cmp funcCmp[K], desde, hasta *K, visitar func(K, V) bool) bool {
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
		panic("El iterador termino de iterar")
	}
	nodo := iter.pila.Desapilar()

	if nodo.derecho != nil {
		iter.apilarIzquierdos(nodo.derecho, nil)
	}
}
