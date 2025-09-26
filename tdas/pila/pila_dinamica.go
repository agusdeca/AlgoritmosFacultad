package pila

/* Definición del struct pila proporcionado por la cátedra. */
const (
	VALOR_REDIMENSION = 2
	LARGO_MINIMO      = 2
	VALOR_CUARTO      = 4
	CAPACIDAD_INICIAL = 10
	CAPACIDAD_MINIMA  = 5
)

type pilaDinamica[T any] struct {
	datos    []T
	cantidad int
}

func CrearPilaDinamica[T any]() Pila[T] {
	return &pilaDinamica[T]{datos: make([]T, CAPACIDAD_INICIAL), cantidad: 0}
}

func (p *pilaDinamica[T]) EstaVacia() bool {
	return p.cantidad == 0
}

func (p *pilaDinamica[T]) validarPilaVacia() {
	if p.EstaVacia() {
		panic("La pila esta vacia")
	}
}

func (p *pilaDinamica[T]) VerTope() T {
	p.validarPilaVacia()
	return p.datos[p.cantidad-1]
}

func (p *pilaDinamica[T]) redimensionar(capacidadNueva int) {
	nuevoSlice := make([]T, capacidadNueva)
	copy(nuevoSlice, p.datos[:p.cantidad])
	p.datos = nuevoSlice
}

func (p *pilaDinamica[T]) Apilar(elemento T) {
	if p.cantidad == len(p.datos) {
		p.redimensionar(len(p.datos) * VALOR_REDIMENSION)
	}

	p.datos[p.cantidad] = elemento

	p.cantidad++
}

func (p *pilaDinamica[T]) Desapilar() T {
	p.validarPilaVacia()
	elemento := p.datos[p.cantidad-1]
	p.cantidad--

	if p.cantidad > 0 && p.cantidad <= len(p.datos)/VALOR_CUARTO {
		nuevaCapacidad := len(p.datos) / VALOR_REDIMENSION

		if nuevaCapacidad < CAPACIDAD_MINIMA {
			nuevaCapacidad = CAPACIDAD_MINIMA
		}

		p.redimensionar(nuevaCapacidad)
	}

	return elemento
}
