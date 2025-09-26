package lista

type nodoLista[T any] struct {
	dato T
	prox *nodoLista[T]
}

type listaEnlazada[T any] struct {
	primero *nodoLista[T]
	ultimo  *nodoLista[T]
	largo   int
}

// TODO: struct de iterador externo


const(
	MENSAJE_LISTA_VACIA = "La lista esta vacia"
	LARGO_INICIAL= 0
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
		dato: dato,
		prox: nil,
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
		nuevoNodo.prox = l.primero
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
		l.ultimo.prox = nuevoNodo
		l.ultimo = nuevoNodo
	}
	l.largo++
}

func (l *listaEnlazada[T]) BorrarPrimero() T {
	if l.EstaVacia() {
		panic(MENSAJE_LISTA_VACIA)
	}
	
	dato := l.primero.dato
	l.primero = l.primero.prox
	
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

func (l *listaEnlazada[T]) Largo() T {
	return l.largo
}

func (l *listaEnlazada[T]) Iterar(visitar func(T) bool) {
	actual := l.primero
	for actual != nil {
		if !visitar(actual.dato) {
			break
		}
		actual = actual.prox
	}
}

// TODO: Iterador de la fimra de lista que implica al iterador externo

// TODO: firmas del iterador externo 