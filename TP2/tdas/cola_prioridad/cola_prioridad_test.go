package cola_prioridad_test

import (
	"math/rand"
	"testing"

	TDAHeap "tdas/cola_prioridad"

	"github.com/stretchr/testify/require"
)

const (
	CANTIDAD = 10000
)

func cmpInt(a, b int) int {
	return a - b
}

func TestHeapVacio(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)

	require.True(t, heap.EstaVacia())
	require.Equal(t, 0, heap.Cantidad())
	require.PanicsWithValue(t, "La cola esta vacia", func() { heap.VerMax() })
	require.PanicsWithValue(t, "La cola esta vacia", func() { heap.Desencolar() })
}

func TestEncolarUnElemento(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	heap.Encolar(10)
	require.False(t, heap.EstaVacia())
	require.Equal(t, 10, heap.VerMax())
	require.Equal(t, 10, heap.VerMax())
	require.Equal(t, 1, heap.Cantidad())
}

func TestEncolarYDesencolar(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	heap.Encolar(5)
	heap.Encolar(10)
	require.Equal(t, 10, heap.VerMax())
	require.Equal(t, 10, heap.Desencolar())
	require.Equal(t, 5, heap.Desencolar())
	require.True(t, heap.EstaVacia())
	require.Equal(t, 0, heap.Cantidad())
}

func TestHeapMantieneMaximo(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	valores := []int{3, 7, 1, 9, 5}
	for _, v := range valores {
		heap.Encolar(v)
	}
	require.Equal(t, 9, heap.VerMax())
	require.Equal(t, 9, heap.Desencolar())
	require.Equal(t, 7, heap.VerMax())
}

func TestCrearHeapArr(t *testing.T) {
	arr := []int{3, 1, 4, 1, 5, 9, 2}

	arrOriginal := make([]int, len(arr))
	copy(arrOriginal, arr)

	heap := TDAHeap.CrearHeapArr[int](arr, cmpInt)
	require.Equal(t, len(arr), heap.Cantidad())
	require.Equal(t, 9, heap.VerMax())

	require.Equal(t, arrOriginal, arr, "CrearHeapArr NO debe modificar el arreglo original")

	prev := 9999999
	for !heap.EstaVacia() {
		val := heap.Desencolar()
		require.LessOrEqual(t, val, prev)
		prev = val
	}
}

func TestCrearHeapArrVacio(t *testing.T) {
	arr := []int{}
	heap := TDAHeap.CrearHeapArr[int](arr, cmpInt)
	require.True(t, heap.EstaVacia())
	require.Equal(t, 0, heap.Cantidad())
	require.PanicsWithValue(t, "La cola esta vacia", func() { heap.VerMax() })
}

func TestCrearHeapArrUnElemento(t *testing.T) {
	arr := []int{5}
	heap := TDAHeap.CrearHeapArr[int](arr, cmpInt)
	require.False(t, heap.EstaVacia())
	require.Equal(t, 1, heap.Cantidad())
	require.Equal(t, 5, heap.VerMax())
	require.Equal(t, 5, heap.Desencolar())
	require.True(t, heap.EstaVacia())
}

func TestCrearHeapArrVolumen(t *testing.T) {
	arr := make([]int, CANTIDAD)
	for i := 0; i < CANTIDAD; i++ {
		arr[i] = rand.Intn(CANTIDAD * 10)
	}

	arrOriginal := make([]int, len(arr))
	copy(arrOriginal, arr)

	heap := TDAHeap.CrearHeapArr[int](arr, cmpInt)
	require.Equal(t, CANTIDAD, heap.Cantidad())
	require.Equal(t, arrOriginal, arr, "CrearHeapArr NO debe modificar el arreglo original")

	prev := heap.Desencolar()
	for !heap.EstaVacia() {
		val := heap.Desencolar()
		require.LessOrEqual(t, val, prev)
		prev = val
	}
}

func TestPanicHeapVacio(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	require.PanicsWithValue(t, "La cola esta vacia", func() { heap.VerMax() })
	require.PanicsWithValue(t, "La cola esta vacia", func() { heap.Desencolar() })
}

func TestVolumen(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	for i := 0; i < CANTIDAD; i++ {
		heap.Encolar(i)
		require.Equal(t, i, heap.VerMax())
		require.False(t, heap.EstaVacia())
	}
	require.Equal(t, CANTIDAD, heap.Cantidad())

	for i := CANTIDAD - 1; i >= 0; i-- {
		require.Equal(t, i, heap.VerMax())
		require.Equal(t, i, heap.Desencolar())
	}
	require.True(t, heap.EstaVacia())
}

func TestVolumenDescendente(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)
	max := CANTIDAD - 1
	for i := max; i >= 0; i-- {
		heap.Encolar(i)
		require.Equal(t, max, heap.VerMax())
	}
	require.Equal(t, CANTIDAD, heap.Cantidad())

	for i := max; i >= 0; i-- {
		require.Equal(t, i, heap.Desencolar())
	}
	require.True(t, heap.EstaVacia())
}

func TestHeapGenericoString(t *testing.T) {
	cmpStr := func(a, b string) int {
		if a > b {
			return 1
		} else if a < b {
			return -1
		}
		return 0
	}
	heap := TDAHeap.CrearHeap[string](cmpStr)
	heap.Encolar("perro")
	heap.Encolar("gato")
	heap.Encolar("zorro")

	require.Equal(t, "zorro", heap.VerMax())
	require.Equal(t, "zorro", heap.Desencolar())
	require.Equal(t, "perro", heap.VerMax())
}

func TestRedimension(t *testing.T) {
	heap := TDAHeap.CrearHeap[int](cmpInt)

	for i := 0; i < 25; i++ {
		heap.Encolar(i)
	}
	require.Equal(t, 25, heap.Cantidad())
	require.Equal(t, 24, heap.VerMax())

	for i := 0; i < 15; i++ {
		heap.Desencolar()
	}
	require.Equal(t, 10, heap.Cantidad())
	require.Equal(t, 9, heap.VerMax())

	for i := 0; i < 5; i++ {
		heap.Desencolar()
	}
	require.Equal(t, 5, heap.Cantidad())

	heap.Desencolar()
	heap.Desencolar()
	heap.Desencolar()
	require.Equal(t, 2, heap.Cantidad())

	heap.Desencolar()
	heap.Desencolar()
	require.True(t, heap.EstaVacia())
}

func TestHeapSortVacio(t *testing.T) {
	var arr []int
	TDAHeap.HeapSort(arr, cmpInt)
	require.Nil(t, arr)

	arr = make([]int, 0)
	TDAHeap.HeapSort(arr, cmpInt)
	require.NotNil(t, arr)
	require.Len(t, arr, 0)
}

func TestHeapSortUnElemento(t *testing.T) {
	t.Log("Prueba: HeapSort con un elemento")
	arr := []int{5}
	TDAHeap.HeapSort(arr, cmpInt)
	require.Equal(t, []int{5}, arr)
}

func TestHeapSortBasico(t *testing.T) {
	t.Log("Prueba: HeapSort con varios elementos")
	arr := []int{5, 3, 8, 1, 10, 2, 7, 4, 6, 9}
	esperado := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	TDAHeap.HeapSort(arr, cmpInt)
	require.Equal(t, esperado, arr)
}

func TestHeapSortDuplicados(t *testing.T) {
	t.Log("Prueba: HeapSort con elementos duplicados")
	arr := []int{5, 1, 5, 8, 1, 10, 8, 5, 5, 1}
	esperado := []int{1, 1, 1, 5, 5, 5, 5, 8, 8, 10}
	TDAHeap.HeapSort(arr, cmpInt)
	require.Equal(t, esperado, arr)
}

func TestHeapSortYaOrdenadoDesc(t *testing.T) {
	arr := make([]int, 100)
	esperado := make([]int, 100)
	for i := 0; i < 100; i++ {
		arr[i] = 99 - i
		esperado[i] = i
	}

	TDAHeap.HeapSort(arr, cmpInt)
	require.Equal(t, esperado, arr)
}

func TestHeapSortYaOrdenadoAsc(t *testing.T) {
	arr := make([]int, 100)
	for i := 0; i < 100; i++ {
		arr[i] = i
	}
	esperado := make([]int, 100)
	copy(esperado, arr)

	TDAHeap.HeapSort(arr, cmpInt)
	require.Equal(t, esperado, arr)
}

func TestHeapSortVolumen(t *testing.T) {
	arr := make([]int, CANTIDAD)
	for i := 0; i < CANTIDAD; i++ {
		arr[i] = CANTIDAD - i
	}

	TDAHeap.HeapSort(arr, cmpInt)

	for i := 1; i < CANTIDAD; i++ {
		require.LessOrEqual(t, arr[i-1], arr[i])
	}
}
