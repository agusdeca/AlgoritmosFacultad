package diccionario

import (
	TDAPila "tdas/pila"
)

const (
	MENOR = -1
	IGUAL = 0
	MAYOR = 1
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
	hasta *K
}

// CrearABB crea un nuevo ABB vacío con la func de cmp
func CrearABB[K comparable, V any](cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: cmp}
}

// Guardar inserta o reemplaza una clave
func (a *abb[K, V]) Guardar(clave K, dato V) {
	a.raiz = insertarNodo(a.raiz, clave, dato, a.cmp, &a.cantidad)
}

// Func aux recursiva
func insertarNodo[K comparable, V any](nodo *nodoAbb[K, V], clave K, dato V, cmp funcCmp[K], cant *int) *nodoAbb[K, V] {
	if nodo == nil {
		*cant++
		return &nodoAbb[K, V]{clave: clave, dato: dato}
	}
	comp := cmp(clave, nodo.clave)
	if comp == IGUAL {
		nodo.dato = dato
	} else if comp == MENOR {
		nodo.izquierdo = insertarNodo(nodo.izquierdo, clave, dato, cmp, cant)
	} else {
		nodo.derecho = insertarNodo(nodo.derecho, clave, dato, cmp, cant)
	}
	return nodo
}

// Pertenece indica si una clave está en el ABB
func (a *abb[K, V]) Pertenece(clave K) bool {
	return obtenerNodo(a.raiz, clave, a.cmp) != nil
}

// Obtener devuelve el valor asociado a una clave
func (a *abb[K, V]) Obtener(clave K) V {
	nodo := obtenerNodo(a.raiz, clave, a.cmp)
	if nodo == nil {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	return nodo.dato
}

// Func aux recursiva que recorre el ABB teniendo en cuenta si la clave es menor o mayor
func obtenerNodo[K comparable, V any](nodo *nodoAbb[K, V], clave K, cmp func(K, K) int) *nodoAbb[K, V] {
	if nodo == nil {
		return nil
	}
	comp := cmp(clave, nodo.clave)
	if comp == IGUAL {
		return nodo
	}
	if comp == MENOR {
		return obtenerNodo(nodo.izquierdo, clave, cmp)
	}
	return obtenerNodo(nodo.derecho, clave, cmp)
}

// Borrar elimina una clave del ABB y devuelve su valor
func (a *abb[K, V]) Borrar(clave K) V {
	nuevaRaiz, borrado, ok := borrar(a.raiz, clave, a.cmp)
	if !ok {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	a.raiz = nuevaRaiz
	a.cantidad--
	return borrado
}

// Func aux recursiva que busca el nodo a eliminar
func borrar[K comparable, V any](nodo *nodoAbb[K, V], clave K, cmp funcCmp[K]) (*nodoAbb[K, V], V, bool) {
	if nodo == nil {
		var cero V
		return nil, cero, false
	}
	comp := cmp(clave, nodo.clave)
	if comp == MENOR {
		izq, v, ok := borrar(nodo.izquierdo, clave, cmp)
		nodo.izquierdo = izq
		return nodo, v, ok
	}
	if comp == MAYOR {
		der, v, ok := borrar(nodo.derecho, clave, cmp)
		nodo.derecho = der
		return nodo, v, ok
	}
	return borrarNodoActual(nodo, cmp)
}

// func aux que borra dependiendo de los casos del nodo a borrar
func borrarNodoActual[K comparable, V any](nodo *nodoAbb[K, V], cmp funcCmp[K]) (*nodoAbb[K, V], V, bool) {
	dato := nodo.dato
	// Caso 1: sin hijos, se borra el nodo
	if nodo.izquierdo == nil && nodo.derecho == nil {
		return nil, dato, true
	}
	// Caso 2: un hijo, se reemplaza el nodo por su hijo
	if nodo.izquierdo == nil {
		return nodo.derecho, dato, true
	}
	if nodo.derecho == nil {
		return nodo.izquierdo, dato, true
	}
	// Caso 3: dos hijos
	sucesor := minimoNodo(nodo.derecho)
	nodo.clave, nodo.dato = sucesor.clave, sucesor.dato
	nodo.derecho, _, _ = borrar(nodo.derecho, sucesor.clave, cmp)
	return nodo, dato, true
}

// Func aux que busca al minimo del subarbol
func minimoNodo[K comparable, V any](nodo *nodoAbb[K, V]) *nodoAbb[K, V] {
	for nodo.izquierdo != nil {
		nodo = nodo.izquierdo
	}
	return nodo
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
	iterarRango(a.raiz, a.cmp, desde, hasta, visitar)
}

// Func aux que recorre el arbol recursivamente y ve solo las claves dentro del rango
func iterarRango[K comparable, V any](nodo *nodoAbb[K, V], cmp funcCmp[K], desde, hasta *K, visitar func(K, V) bool) bool {
	if nodo == nil {
		return true
	}

	if fueraDeLimiteInferior(nodo.clave, desde, cmp) {
		if !iterarRango(nodo.izquierdo, cmp, desde, hasta, visitar) {
			return false
		}
	}

	if enRango(nodo.clave, desde, hasta, cmp) {
		if !visitar(nodo.clave, nodo.dato) {
			return false
		}
	}

	if fueraDeLimiteSuperior(nodo.clave, hasta, cmp) {
		if !iterarRango(nodo.derecho, cmp, desde, hasta, visitar) {
			return false
		}
	}
	return true
}

// Func aux que analizo si bajo a la izq
func fueraDeLimiteInferior[K comparable](clave K, desde *K, cmp funcCmp[K]) bool {
	return desde == nil || cmp(clave, *desde) > 0
}

// Func aux que analizo si bajo a la der
func fueraDeLimiteSuperior[K comparable](clave K, hasta *K, cmp funcCmp[K]) bool {
	return hasta == nil || cmp(clave, *hasta) < 0
}

// Func aux que analizo si veo el nodo actual
func enRango[K comparable](clave K, desde, hasta *K, cmp funcCmp[K]) bool {
	return (desde == nil || cmp(*desde, clave) <= 0) &&
		(hasta == nil || cmp(clave, *hasta) <= 0)
}

// Iterador externo
