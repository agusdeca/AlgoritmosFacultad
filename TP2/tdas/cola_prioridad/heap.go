package cola_prioridad

const (
	TAM_INICIAL       = 10
	VALOR_REDIMENSION = 2
	VALOR_DISMINUCION = 4
	MSJ_PANICO        = "La cola esta vacia"
)

type cola_prioridad[T any] struct {
	datos []T
	cant  int
	cmp   func(T, T) int
}

func CrearHeap[T any](funcion_cmp func(T, T) int) ColaPrioridad[T] {
	return crearHeapConCap[T](TAM_INICIAL, funcion_cmp)
}

// CrearHeapArr crea un heap a partir de un arreglo en O(n) (heapify)
func CrearHeapArr[T any](arreglo []T, funcion_cmp func(T, T) int) ColaPrioridad[T] {
	// Copiamos los eleme
	n := len(arreglo)
	capacidad := n
	if capacidad < TAM_INICIAL {
		capacidad = TAM_INICIAL
	}
	datos := make([]T, capacidad)
	copy(datos, arreglo)
	h := &cola_prioridad[T]{datos: datos, cant: n, cmp: funcion_cmp}
	heapify(h.datos[:h.cant], h.cmp)
	return h
}

// func aux para crear el heap con capacidad específica
func crearHeapConCap[T any](capacidad int, cmp func(T, T) int) *cola_prioridad[T] {
	if capacidad < 1 {
		capacidad = TAM_INICIAL
	}
	datos := make([]T, capacidad)
	return &cola_prioridad[T]{datos: datos, cant: 0, cmp: cmp}
}

// EstaVacia devuelve true si la la cola se encuentra vacía, false en caso contrario.
func (h *cola_prioridad[T]) EstaVacia() bool {
	return h.cant == 0
}

// Encolar Agrega un elemento al heap.
func (h *cola_prioridad[T]) Encolar(elem T) {
	h.redimensionarSubir()
	h.datos[h.cant] = elem
	h.cant++
	h.upheap(h.cant - 1)
}

// VerMax devuelve el elemento con máxima prioridad. Si está vacía, entra en pánico con un mensaje
// "La cola esta vacia".
func (h *cola_prioridad[T]) VerMax() T {
	if h.EstaVacia() {
		panic(MSJ_PANICO)
	}
	return h.datos[0]
}

// Cantidad devuelve la cantidad de elementos que hay en la cola de prioridad.
func (h *cola_prioridad[T]) Cantidad() int {
	return h.cant
}

// auxiliares
func (h *cola_prioridad[T]) padre(i int) int {
	return (i - 1) / 2
}

func hijoIzq(i int) int { return 2*i + 1 }
func hijoDer(i int) int { return 2*i + 2 }

// upheap: mueve el elemento en i hacia arriba hasta cumplir heap
func (h *cola_prioridad[T]) upheap(i int) {
	for i > 0 {
		padre := h.padre(i)
		if h.cmp(h.datos[i], h.datos[padre]) > 0 {
			h.swap(i, padre)
			i = padre
		} else {
			break
		}
	}
}

// downheap: mueve el elemento en i hacia abajo, con límite length = h.cant
func downheap[T any](arr []T, i, cant int, cmp func(T, T) int) {
	for {
		izq := hijoIzq(i)
		der := hijoDer(i)
		mayor := i

		if izq < cant && cmp(arr[izq], arr[mayor]) > 0 {
			mayor = izq
		}
		if der < cant && cmp(arr[der], arr[mayor]) > 0 {
			mayor = der
		}
		if mayor == i {
			break
		}
		arr[i], arr[mayor] = arr[mayor], arr[i]
		i = mayor
	}
}

func heapify[T any](arr []T, cmp func(T, T) int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	for i := (n - 2) / 2; i >= 0; i-- {
		downheap(arr, i, n, cmp)
		if i == 0 {
			break
		}
	}
}

func (h *cola_prioridad[T]) swap(i, j int) {
	h.datos[i], h.datos[j] = h.datos[j], h.datos[i]
}

func (h *cola_prioridad[T]) redimensionar(nuevaCap int) {
	nuevosDatos := make([]T, nuevaCap)
	copy(nuevosDatos, h.datos[:h.cant])
	h.datos = nuevosDatos
}

// redimensionar si la capacidad es poca se duplica
func (h *cola_prioridad[T]) redimensionarSubir() {
	if h.cant < len(h.datos) {
		return
	}
	nuevaCap := len(h.datos) * VALOR_REDIMENSION
	if nuevaCap == 0 {
		nuevaCap = TAM_INICIAL
	}
	h.redimensionar(nuevaCap)
}

func (h *cola_prioridad[T]) Desencolar() T {
	if h.EstaVacia() {
		panic(MSJ_PANICO)
	}
	max := h.datos[0]
	h.swap(0, h.cant-1)
	h.cant--
	downheap(h.datos, 0, h.cant, h.cmp)

	h.redimensionarBajar()
	return max
}

func (h *cola_prioridad[T]) redimensionarBajar() {
	capacidad := len(h.datos)
	if h.cant*VALOR_DISMINUCION <= capacidad && capacidad > TAM_INICIAL {
		nuevaCap := capacidad / VALOR_REDIMENSION
		if nuevaCap < TAM_INICIAL {
			nuevaCap = TAM_INICIAL
		}
		h.redimensionar(nuevaCap)
	}
}

func HeapSort[T any](elementos []T, funcion_cmp func(T, T) int) {
	n := len(elementos)
	if n <= 1 {
		return
	}

	heapify(elementos, funcion_cmp)

	for i := n - 1; i > 0; i-- {
		elementos[0], elementos[i] = elementos[i], elementos[0]
		downheap(elementos, 0, i, funcion_cmp)
	}
}
