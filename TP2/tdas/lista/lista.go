package lista

type Lista[T any] interface {

	// EstaVacia devuelve verdadero si la lista no tiene elementos, false en caso contrario.
	EstaVacia() bool

	// InsertarPrimero agrega un nuevo elemento a la lista, al principio de la misma.
	InsertarPrimero(T)

	// InsertarUltimo agrega un nuevo elemento a la lista, al final de la misma.
	InsertarUltimo(T)

	// BorrarPrimero saca el primer elemento de la lista. Si la lista tiene elementos, se quita el primero de la misma,
	// y se devuelve ese valor. Si está vacía, entra en pánico con un mensaje "La lista esta vacia".
	BorrarPrimero() T

	// VerPrimero obtiene el valor del primero de la lista. Si está vacía, entra en pánico con un mensaje
	// "La lista esta vacia".
	VerPrimero() T

	// VerUltimo obtiene el valor del ultimo de la lista. Si está vacía, entra en pánico con un mensaje
	// "La lista esta vacia".
	VerUltimo() T

	// Largo devuelve la cantidad de elementos de la lista.
	Largo() int

	// Iterar aplica la función visitar a cada elemento de la lista en orden,
	// desde el primero hasta el último, hasta que se termine la lista o
	// la función visitar devuelva false.
	Iterar(visitar func(T) bool)

	// Iterador devuelve un iterador externo posicionado al inicio de la lista.
	Iterador() IteradorLista[T]
}

type IteradorLista[T any] interface {
	// VerActual devuelve el elemento actual del iterador.
	VerActual() T

	// HaySiguiente indica si hay un elemento siguiente en la iteración.
	HaySiguiente() bool

	// Siguiente avanza el iterador al siguiente elemento.
	Siguiente()

	// Insertar agrega un elemento en la posición actual del iterador.
	Insertar(T)

	// Borrar elimina y devuelve el elemento actual del iterador.
	Borrar() T
}