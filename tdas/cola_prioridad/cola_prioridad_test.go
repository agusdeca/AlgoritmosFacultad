package cola_prioridad_test

import (
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
	heap := TDAHeap.CrearHeapArr[int](arr, cmpInt)
	require.Equal(t, len(arr), heap.Cantidad())
	require.Equal(t, 9, heap.VerMax())

	prev := 9999999
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
