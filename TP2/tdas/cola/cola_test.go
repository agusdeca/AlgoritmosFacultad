package cola_test

import (
	"github.com/stretchr/testify/require"
	TDACola "tdas/cola"
	"testing"
)

const (
	CANTIDAD = 10000
)

func TestColaVacia(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	require.True(t, cola.EstaVacia())
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.VerPrimero() })
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.Desencolar() })
}

func TestEncolarUnElemento(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	cola.Encolar(5)
	require.False(t, cola.EstaVacia())
	require.Equal(t, 5, cola.VerPrimero())
	require.Equal(t, 5, cola.VerPrimero())
	require.False(t, cola.EstaVacia())
}

func TestEncolarDesencolarUnElemento(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	cola.Encolar(10)
	require.Equal(t, 10, cola.Desencolar())
	require.True(t, cola.EstaVacia())
}

func TestIntercalarOperaciones(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	cola.Encolar(1)
	cola.Encolar(2)
	require.Equal(t, 1, cola.Desencolar())
	cola.Encolar(3)
	require.Equal(t, 2, cola.Desencolar())
	require.Equal(t, 3, cola.Desencolar())
	require.True(t, cola.EstaVacia())
}

func TestInvarianteFIFO(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	elementos := []int{1, 5, 10, 15, 20}
	for _, elem := range elementos {
		cola.Encolar(elem)
		require.Equal(t, elementos[0], cola.VerPrimero())
		require.False(t, cola.EstaVacia())
	}

	for i := 0; i < len(elementos); i++ {
		require.Equal(t, elementos[i], cola.VerPrimero())
		require.Equal(t, elementos[i], cola.Desencolar())
	}
	require.True(t, cola.EstaVacia())
}

func TestVolumen(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[int]()
	for i := 0; i < CANTIDAD; i++ {
		cola.Encolar(i * 2)
		require.Equal(t, 0, cola.VerPrimero())
		require.False(t, cola.EstaVacia())
	}

	for i := 0; i < CANTIDAD; i++ {
		valorEsperado := i * 2
		require.Equal(t, valorEsperado, cola.VerPrimero())
		require.Equal(t, valorEsperado, cola.Desencolar())
		if i < CANTIDAD-1 {
			require.Equal(t, (i+1)*2, cola.VerPrimero())
		}
	}
	require.True(t, cola.EstaVacia())
}

func TestColaVaciaEqRecienCreada(t *testing.T) {
	cola := TDACola.CrearColaEnlazada[string]()
	cola.Encolar("hola")
	cola.Encolar("mundo")
	cola.Desencolar()
	cola.Desencolar()
	require.True(t, cola.EstaVacia())
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.VerPrimero() })
	require.PanicsWithValue(t, "La cola esta vacia", func() { cola.Desencolar() })

	cola.Encolar("nuevo")
	require.Equal(t, "nuevo", cola.VerPrimero())
	require.False(t, cola.EstaVacia())
}