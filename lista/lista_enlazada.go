package lista

type nodoLista[T any] struct {
	dato      T
	siguiente *nodoLista[T]
}

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

type iterListaEnlazada[T any] struct {
	lista    *listaEnlazada[T]
	actual   *nodoLista[T]
	anterior *nodoLista[T]
}

const (
	MENSAJE_LISTA_VACIA        = "La lista esta vacia"
	MENSAJE_ITERADOR_TERMINADO = "El iterador termino de iterar"
	LARGO_INICIAL              = 0
)

func CrearListaEnlazada[T any]() Lista[T] {
	return &listaEnlazada[T]{
		primero: nil,
		ultimo:  nil,
		largo:   LARGO_INICIAL,
	}
}

func crearNodo[T any](dato T) *nodoLista[T] {
	return &nodoLista[T]{
		dato:      dato,
		siguiente: nil,
	}
}

func (l *listaEnlazada[T]) EstaVacia() bool {
	return l.largo == LARGO_INICIAL
}

func (l *listaEnlazada[T]) InsertarPrimero(dato T) {
	nuevoNodo := crearNodo(dato)

	if l.EstaVacia() {
		l.primero = nuevoNodo
		l.ultimo = nuevoNodo
	} else {
		nuevoNodo.siguiente = l.primero
		l.primero = nuevoNodo
	}
	l.largo++
}

func (l *listaEnlazada[T]) InsertarUltimo(dato T) {
	nuevoNodo := crearNodo(dato)

	if l.EstaVacia() {
		l.primero = nuevoNodo
		l.ultimo = nuevoNodo
	} else {
		l.ultimo.siguiente = nuevoNodo
		l.ultimo = nuevoNodo
	}
	l.largo++
}

func (l *listaEnlazada[T]) BorrarPrimero() T {
	if l.EstaVacia() {
		panic(MENSAJE_LISTA_VACIA)
	}

	dato := l.primero.dato
	l.primero = l.primero.siguiente

	if l.primero == nil {
		l.ultimo = nil
	}

	l.largo--
	return dato
}

func (l *listaEnlazada[T]) VerPrimero() T {
	if l.EstaVacia() {
		panic(MENSAJE_LISTA_VACIA)
	}
	return l.primero.dato
}

func (l *listaEnlazada[T]) VerUltimo() T {
	if l.EstaVacia() {
		panic(MENSAJE_LISTA_VACIA)
	}
	return l.ultimo.dato
}

func (l *listaEnlazada[T]) Largo() int {
	return l.largo
}

func (l *listaEnlazada[T]) Iterar(visitar func(T) bool) {
	actual := l.primero
	for actual != nil {
		if !visitar(actual.dato) {
			break
		}
		actual = actual.siguiente
	}
}

func (l *listaEnlazada[T]) Iterador() IteradorLista[T] {
	return &iterListaEnlazada[T]{
		lista:    l,
		actual:   l.primero,
		anterior: nil}
}

func (it *iterListaEnlazada[T]) HaySiguiente() bool {
	return it.actual != nil
}

func (it *iterListaEnlazada[T]) VerActual() T {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITERADOR_TERMINADO)
	}
	return it.actual.dato
}

func (it *iterListaEnlazada[T]) Siguiente() {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITERADOR_TERMINADO)
	}
	it.anterior = it.actual
	it.actual = it.actual.siguiente
}

func (it *iterListaEnlazada[T]) Insertar(dato T) {
	nuevo := crearNodo(dato)

	if it.anterior == nil {
		nuevo.siguiente = it.lista.primero
		it.lista.primero = nuevo
		if it.lista.ultimo == nil {
			it.lista.ultimo = nuevo
		}
		it.actual = nuevo
	} else {
		nuevo.siguiente = it.actual
		it.anterior.siguiente = nuevo
		it.actual = nuevo
		if nuevo.siguiente == nil {
			it.lista.ultimo = nuevo
		}
	}
	it.lista.largo++
}

func (it *iterListaEnlazada[T]) Borrar() T {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITERADOR_TERMINADO)
	}
	nodoActual := it.actual
	dato := nodoActual.dato

	if it.anterior == nil {
		it.lista.primero = nodoActual.siguiente
		it.actual = nodoActual.siguiente
		if it.lista.primero == nil {
			it.lista.ultimo = nil
		}
	} else {
		it.anterior.siguiente = nodoActual.siguiente
		it.actual = nodoActual.siguiente
		if nodoActual == it.lista.ultimo {
			it.lista.ultimo = it.anterior
		}
	}
	it.lista.largo--
	return dato
}
