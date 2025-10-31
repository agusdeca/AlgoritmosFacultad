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
	valores := []int{10, 5, 20}
	esperadosPrimero := []int{10, 5, 5}
	esperadosUltimo := []int{10, 10, 20}
	largos := []int{1, 2, 3}

	l.InsertarPrimero(valores[0])
	require.False(t, l.EstaVacia())

	for i := 0; i < len(valores); i++ {
		if i == 1 {
			l.InsertarPrimero(valores[i])
		} else if i == 2 {
			l.InsertarUltimo(valores[i])
		}
		require.Equal(t, largos[i], l.Largo())
		require.Equal(t, esperadosPrimero[i], l.VerPrimero())
		require.Equal(t, esperadosUltimo[i], l.VerUltimo())
	}
}

func TestBorrar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	valores := []int{1, 2}

	for _, v := range valores {
		l.InsertarUltimo(v)
	}

	require.Equal(t, 1, l.BorrarPrimero())
	require.Equal(t, 1, l.Largo())
	require.Equal(t, 2, l.VerPrimero())
	require.Equal(t, 2, l.VerUltimo())

	require.Equal(t, 2, l.BorrarPrimero())
	require.True(t, l.EstaVacia())
	require.PanicsWithValue(t, "La lista esta vacia", func() { l.BorrarPrimero() })
}

func TestVolumen(t *testing.T) {
	const n = 10000
	l := TDALista.CrearListaEnlazada[int]()

	// Inserto al final
	for i := 0; i < n; i++ {
		l.InsertarUltimo(i)
	}
	require.Equal(t, n, l.Largo())
	require.Equal(t, 0, l.VerPrimero())
	require.Equal(t, n-1, l.VerUltimo())

	// Borrar todos
	for i := 0; i < n; i++ {
		require.Equal(t, i, l.BorrarPrimero())
	}
	require.True(t, l.EstaVacia())

	// Insertar al principio
	for i := 0; i < n; i++ {
		l.InsertarPrimero(i)
	}
	require.Equal(t, n, l.Largo())
	require.Equal(t, n-1, l.VerPrimero())
	require.Equal(t, 0, l.VerUltimo())
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
	valores := []string{"a", "b", "c"}

	for _, v := range valores {
		l.InsertarUltimo(v)
	}

	iter := l.Iterador()
	resultado := []string{}
	for iter.HaySiguiente() {
		resultado = append(resultado, iter.VerActual())
		iter.Siguiente()
	}
	require.Equal(t, valores, resultado)
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

func TestIteradorExtInsertarListaVacia(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	iter := l.Iterador()

	iter.Insertar(1)
	require.Equal(t, 1, l.Largo())
	require.Equal(t, 1, l.VerPrimero())
	require.Equal(t, 1, l.VerUltimo())

	iter.Insertar(2)
	require.Equal(t, 2, l.Largo())
	require.Equal(t, 2, l.VerPrimero())
	require.Equal(t, 1, l.VerUltimo())
}

func TestIteradorExtBorrar(t *testing.T) {
	l := TDALista.CrearListaEnlazada[int]()
	valores := []int{1, 2, 3}

	for _, v := range valores {
		l.InsertarUltimo(v)
	}

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
	const n = 10000
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

	iter = l.Iterador()
	for i := 0; i < n/2 && iter.HaySiguiente(); i++ {
		iter.Borrar()
	}
	require.Equal(t, n/2, l.Largo())

	iter = l.Iterador()
	for i := 0; i < 100 && iter.HaySiguiente(); i++ {
		iter.Insertar(-1)
		iter.Siguiente()
	}
	require.Equal(t, n/2+100, l.Largo())
}