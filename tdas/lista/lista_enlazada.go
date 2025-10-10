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
)

//Lista

func CrearListaEnlazada[T any]() Lista[T] {
	return &listaEnlazada[T]{
		primero: nil,
		ultimo:  nil,
		largo:   0,
	}
}

func crearNodo[T any](dato T) *nodoLista[T] {
	return &nodoLista[T]{
		dato:      dato,
		siguiente: nil,
	}
}

func (l *listaEnlazada[T]) verificarNoVacia() {
	if l.EstaVacia() {
		panic(MENSAJE_LISTA_VACIA)
	}
}

func (l *listaEnlazada[T]) actualizarSiQuedaVacia() {
	if l.primero == nil {
		l.ultimo = nil
	}
}

func (l *listaEnlazada[T]) EstaVacia() bool {
	return l.largo == 0
}

func (l *listaEnlazada[T]) InsertarPrimero(dato T) {
	nuevoNodo := crearNodo(dato)

	if l.EstaVacia() {
		l.ultimo = nuevoNodo
	} else {
		nuevoNodo.siguiente = l.primero
	}
	l.primero = nuevoNodo
	l.largo++
}

func (l *listaEnlazada[T]) InsertarUltimo(dato T) {
	nuevoNodo := crearNodo(dato)

	if l.EstaVacia() {
		l.primero = nuevoNodo
	} else {
		l.ultimo.siguiente = nuevoNodo
	}
	l.ultimo = nuevoNodo
	l.largo++
}

func (l *listaEnlazada[T]) BorrarPrimero() T {
	l.verificarNoVacia()

	dato := l.primero.dato
	l.primero = l.primero.siguiente

	l.actualizarSiQuedaVacia()

	l.largo--
	return dato
}

func (l *listaEnlazada[T]) VerPrimero() T {
	l.verificarNoVacia()
	return l.primero.dato
}

func (l *listaEnlazada[T]) VerUltimo() T {
	l.verificarNoVacia()
	return l.ultimo.dato
}

func (l *listaEnlazada[T]) Largo() int {
	return l.largo
}

func (l *listaEnlazada[T]) Iterar(visitar func(T) bool) {
	for actual := l.primero; actual != nil; actual = actual.siguiente {
		if !visitar(actual.dato) {
			break
		}
	}
}

func (l *listaEnlazada[T]) Iterador() IteradorLista[T] {
	return &iterListaEnlazada[T]{
		lista:    l,
		actual:   l.primero,
		anterior: nil}
}

//Iterador externo

func (it *iterListaEnlazada[T]) verificarNoTerminado() {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITERADOR_TERMINADO)
	}
}

func (it *iterListaEnlazada[T]) HaySiguiente() bool {
	return it.actual != nil
}

func (it *iterListaEnlazada[T]) VerActual() T {
	it.verificarNoTerminado()
	return it.actual.dato
}

func (it *iterListaEnlazada[T]) Siguiente() {
	it.verificarNoTerminado()
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
	} else {
		nuevo.siguiente = it.actual
		it.anterior.siguiente = nuevo
		if it.actual == nil {
			it.lista.ultimo = nuevo
		}
	}
	it.actual = nuevo
	it.lista.largo++
}

func (it *iterListaEnlazada[T]) Borrar() T {
	it.verificarNoTerminado()
	dato := it.actual.dato

	if it.anterior == nil {
		it.lista.primero = it.actual.siguiente
		if it.lista.primero == nil {
			it.lista.ultimo = nil
		}
	} else {
		it.anterior.siguiente = it.actual.siguiente
		if it.actual == it.lista.ultimo {
			it.lista.ultimo = it.anterior
		}
	}
	it.actual = it.actual.siguiente
	it.lista.largo--
	return dato
}
