package diccionario

import (
	TDAPila "tdas/pila"
)

// funcion de comparacion
type funcCmp[K any] func(K, K) int

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
	a.raiz = a.insertarRec(a.raiz, clave, dato)
}

// Insertar recursivamente en el arbol
func (a *abb[K, V]) insertarRec(n *nodoAbb[K, V], clave K, dato V) *nodoAbb[K, V] {
	if n == nil {
		a.cantidad++
		return &nodoAbb[K, V]{clave: clave, dato: dato}
	}

	comp := a.cmp(clave, n.clave)
	if comp == 0 {
		n.dato = dato
	} else if comp < 0 {
		n.izquierdo = a.insertarRec(n.izquierdo, clave, dato)
	} else {
		n.derecho = a.insertarRec(n.derecho, clave, dato)
	}
	return n
}

// Pertenece indica si una clave está en el ABB
func (a *abb[K, V]) Pertenece(clave K) bool {
	return a.buscarRec(a.raiz, clave) != nil
}

// Obtener devuelve el valor asociado a una clave
func (a *abb[K, V]) Obtener(clave K) V {
	nodo := a.buscarRec(a.raiz, clave)
	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	return nodo.dato
}

// Buscar recursivamente de una clave
func (a *abb[K, V]) buscarRec(n *nodoAbb[K, V], clave K) *nodoAbb[K, V] {
	if n == nil {
		return nil
	}
	comp := a.cmp(clave, n.clave)
	if comp == 0 {
		return n
	}
	if comp < 0 {
		return a.buscarRec(n.izquierdo, clave)
	}
	return a.buscarRec(n.derecho, clave)
}

// Borrar elimina una clave del ABB y devuelve su valor
func (a *abb[K, V]) Borrar(clave K) V {
	nuevaRaiz, borrado, ok := a.borrarRec(a.raiz, clave)
	if !ok {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	a.raiz = nuevaRaiz
	a.cantidad--
	return borrado
}

// Borro recursivo con los tres casos
func (a *abb[K, V]) borrarRec(n *nodoAbb[K, V], clave K) (*nodoAbb[K, V], V, bool) {
	if n == nil {
		var cero V
		return nil, cero, false
	}

	comp := a.cmp(clave, n.clave)
	if comp < 0 {
		var borrado V
		var ok bool
		n.izquierdo, borrado, ok = a.borrarRec(n.izquierdo, clave)
		return n, borrado, ok
	}
	if comp > 0 {
		var borrado V
		var ok bool
		n.derecho, borrado, ok = a.borrarRec(n.derecho, clave)
		return n, borrado, ok
	}

	// Caso base: encontramos la clave
	borrado := n.dato
	if n.izquierdo == nil && n.derecho == nil {
		return nil, borrado, true
	}
	if n.izquierdo == nil {
		return n.derecho, borrado, true
	}
	if n.derecho == nil {
		return n.izquierdo, borrado, true
	}

	min := a.minimo(n.derecho)
	n.clave, n.dato = min.clave, min.dato
	n.derecho, _, _ = a.borrarRec(n.derecho, min.clave)
	return n, borrado, true
}

// Devuelve el nodo con la clave min del subarbol
func (a *abb[K, V]) minimo(n *nodoAbb[K, V]) *nodoAbb[K, V] {
	if n.izquierdo == nil {
		return n
	}
	return a.minimo(n.izquierdo)
}

// Cantidad devuelve el num de elem
func (a *abb[K, V]) Cantidad() int {
	return a.cantidad
}

// Iterador interno
func (a *abb[K, V]) Iterar(visitar func(K, V) bool) {
	a.IterarRango(nil, nil, visitar)
}

func (a *abb[K, V]) IterarRango(desde, hasta *K, visitar func(K, V) bool) {
	a.iterarRangoRec(a.raiz, desde, hasta, visitar)
}

// Recorre el arbol recursivamente y mira solo las claves dentro del rango
func (a *abb[K, V]) iterarRangoRec(n *nodoAbb[K, V], desde, hasta *K, visitar func(K, V) bool) bool {
	if n == nil {
		return true
	}
	// Recorro izq si estoy en el rango inf
	if desde == nil || a.cmp(n.clave, *desde) > 0 {
		if !a.iterarRangoRec(n.izquierdo, desde, hasta, visitar) {
			return false
		}
	}
	// Veo el nodo actual si esta en el rango
	if (desde == nil || a.cmp(*desde, n.clave) <= 0) && (hasta == nil || a.cmp(n.clave, *hasta) <= 0) {
		if !visitar(n.clave, n.dato) {
			return false
		}
	}
	// Recorro la der si no estoy en el limite sup
	if hasta == nil || a.cmp(n.clave, *hasta) < 0 {
		if !a.iterarRangoRec(n.derecho, desde, hasta, visitar) {
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
