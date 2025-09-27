package lista_test

import (
	TDALista "tdas/lista"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListaVacia(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	require.True(t, l.EstaVacia())
	require.Equal(t, 0, l.Largo())
	require.PanicsWithValue(t, "La lista esta vacia", func() { l.VerPrimero() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { l.VerUltimo() })
	require.PanicsWithValue(t, "La lista esta vacia", func() { l.BorrarPrimero() })
}

func TestInsertar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	l.InsertarPrimero(10)
	require.False(t, l.EstaVacia())
	require.Equal(t, 1, l.Largo())
	require.Equal(t, 10, l.VerPrimero())
	require.Equal(t, 10, l.VerUltimo())

	l.InsertarPrimero(5)
	require.Equal(t, 2, l.Largo())
	require.Equal(t, 5, l.VerPrimero())
	require.Equal(t, 10, l.VerUltimo())

	l.InsertarUltimo(20)
	require.Equal(t, 3, l.Largo())
	require.Equal(t, 5, l.VerPrimero())
	require.Equal(t, 20, l.VerUltimo())
}

func TestBorrar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	l.InsertarUltimo(1)
	l.InsertarUltimo(2)

	require.Equal(t, 1, l.BorrarPrimero())
	require.Equal(t, 1, l.Largo())
	require.Equal(t, 2, l.VerPrimero())
	require.Equal(t, 2, l.VerUltimo())

	require.Equal(t, 2, l.BorrarPrimero())
	require.True(t, l.EstaVacia())
	require.PanicsWithValue(t, "La lista esta vacia", func() { l.BorrarPrimero() })
}

// Pruebas iterador interno
func TestIteradorInterno(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	for i := 1; i <= 10; i++ {
		l.InsertarUltimo(i)
	}
	suma := 0
	l.Iterar(func(v int) bool {
		suma += v
		return true
	})
	require.Equal(t, 55, suma)

	sum2 := 0
	l.Iterar(func(v int) bool {
		if v == 7 {
			return false
		}
		if v%2 == 0 {
			sum2 += v
		}
		return true
	})
	require.Equal(t, 12, sum2)
}

// Pruebas iterador externo
func TestIteradorExtRecorrer(t *testing.T) {
	l := TDALista.CrearListaEnlazada[string]()
	l.InsertarUltimo("a")
	l.InsertarUltimo("b")
	l.InsertarUltimo("c")

	iter := l.Iterador()
	resultado := []string{}
	for iter.HaySiguiente() {
		resultado = append(resultado, iter.VerActual())
		iter.Siguiente()
	}
	require.Equal(t, []string{"a", "b", "c"}, resultado)
}

func TestIteradorExtInsertar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	l.InsertarUltimo(2)
	l.InsertarUltimo(3)

	// Insertar al principio
	iter := l.Iterador()
	iter.Insertar(1)
	require.Equal(t, 3, l.Largo())
	require.Equal(t, 1, l.VerPrimero())

	// Insertar en el medio
	iter.Siguiente()
	iter.Insertar(5)
	arr := []int{}
	l.Iterar(func(v int) bool {
		arr = append(arr, v)
		return true
	})
	require.Equal(t, []int{1, 5, 2, 3}, arr)

	// Insertar al final
	for iter.HaySiguiente() {
		iter.Siguiente()
	}
	iter.Insertar(10)
	require.Equal(t, 5, l.Largo())
	require.Equal(t, 10, l.VerUltimo())
}

func TestIteradorExtBorrar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	l.InsertarUltimo(1)
	l.InsertarUltimo(2)
	l.InsertarUltimo(3)

	iter := l.Iterador()
	require.Equal(t, 1, iter.Borrar()) // borrar primero
	require.Equal(t, 2, l.VerPrimero())

	iter = l.Iterador()
	require.Equal(t, 2, iter.Borrar()) // borrar del medio
	require.Equal(t, 1, l.Largo())

	iter = l.Iterador()
	require.Equal(t, 3, iter.Borrar()) // borrar ultimo
	require.True(t, l.EstaVacia())
}

func TestIteradorExtPanics(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	l.InsertarUltimo(7)
	iter := l.Iterador()
	iter.Borrar()

	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.VerActual() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Siguiente() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Borrar() })
}

func TestIteradorExtVolumen(t *testing.T) {
	const n = 1000
	l := TDALista.CrearListaEnlazada[int]()
	for i := 0; i < n; i++ {
		l.InsertarUltimo(i)
	}
	iter := l.Iterador()
	cont := 0
	for iter.HaySiguiente() {
		require.Equal(t, cont, iter.VerActual())
		iter.Siguiente()
		cont++
	}
	require.Equal(t, n, cont)
}
