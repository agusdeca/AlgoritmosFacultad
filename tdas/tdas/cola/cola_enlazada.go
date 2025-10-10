package cola

type nodo[T any] struct {
	dato      T
	siguiente *nodo[T]
}

type colaEnlazada[T any] struct {
	primero *nodo[T]
	ultimo  *nodo[T]
}

func CrearColaEnlazada[T any]() Cola[T] {
	return &colaEnlazada[T]{primero: nil, ultimo: nil}
}

func (c *colaEnlazada[T]) EstaVacia() bool {
	return c.primero == nil
}

func (c *colaEnlazada[T]) validarColaVacia() {
	if c.EstaVacia() {
		panic("La cola esta vacia")
	}
}

func (c *colaEnlazada[T]) VerPrimero() T {
	c.validarColaVacia()
	return c.primero.dato
}

func (c *colaEnlazada[T]) Encolar(elemento T) {
	nuevoNodo := &nodo[T]{dato: elemento, siguiente: nil}

	if c.EstaVacia() {
		c.primero = nuevoNodo
	} else {
		c.ultimo.siguiente = nuevoNodo
	}
	c.ultimo = nuevoNodo

}

func (c *colaEnlazada[T]) Desencolar() T {
	c.validarColaVacia()

	elemento := c.primero.dato
	c.primero = c.primero.siguiente

	if c.EstaVacia() {
		c.ultimo = nil
	}

	return elemento
}
