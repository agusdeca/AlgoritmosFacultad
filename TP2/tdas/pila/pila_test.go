package pila_test

import (
	"github.com/stretchr/testify/require"
	TDAPila "tdas/pila"
	"testing"
)

const (
	CANTIDAD = 10000
)

func TestPilaVacia(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()
	require.True(t, pila.EstaVacia())
	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.VerTope() })
	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.Desapilar() })
}

func TestApilarUnElemento(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	pila.Apilar(5)
	require.False(t, pila.EstaVacia())
	require.Equal(t, 5, pila.VerTope())

	require.Equal(t, 5, pila.VerTope())
	require.False(t, pila.EstaVacia())
}

func TestApilarDesapilarUnElemento(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	pila.Apilar(10)
	require.Equal(t, 10, pila.Desapilar())
	require.True(t, pila.EstaVacia())
}

func TestIntercalarOperaciones(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	pila.Apilar(1)
	pila.Apilar(2)
	require.Equal(t, 2, pila.Desapilar())

	pila.Apilar(3)
	require.Equal(t, 3, pila.Desapilar())
	require.Equal(t, 1, pila.Desapilar())

	require.True(t, pila.EstaVacia())
}

func TestLIFOReducido(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	pila.Apilar(10)
	pila.Apilar(20)
	pila.Apilar(30)

	require.Equal(t, 30, pila.Desapilar())
	require.Equal(t, 20, pila.Desapilar())
	require.Equal(t, 10, pila.Desapilar())

	require.True(t, pila.EstaVacia())
}

func TestInvarianteLIFO(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()

	//Apilo
	elementos := []int{1, 5, 10, 15, 20}
	for _, elem := range elementos {
		pila.Apilar(elem)
		require.Equal(t, elem, pila.VerTope())
		require.False(t, pila.EstaVacia())
	}

	//Desapilo
	for i := len(elementos) - 1; i >= 0; i-- {
		require.Equal(t, elementos[i], pila.VerTope())
		require.Equal(t, elementos[i], pila.Desapilar())
	}

	require.True(t, pila.EstaVacia())
}

func TestVolumen(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[int]()
	for i := 0; i < CANTIDAD; i++ {
		pila.Apilar(i * 2)
		require.Equal(t, i*2, pila.VerTope())
		require.False(t, pila.EstaVacia())
	}

	for i := CANTIDAD - 1; i >= 0; i-- {
		valorEsperado := i * 2
		require.Equal(t, valorEsperado, pila.VerTope())
		require.Equal(t, valorEsperado, pila.Desapilar())

		if i > 0 {
			require.Equal(t, (i-1)*2, pila.VerTope())
		}
	}

	require.True(t, pila.EstaVacia())
}

func TestPilaVaciaEqRecienCreada(t *testing.T) {
	pila := TDAPila.CrearPilaDinamica[string]()

	pila.Apilar("hola")
	pila.Apilar("asd")
	pila.Desapilar()
	pila.Desapilar()

	require.True(t, pila.EstaVacia())
	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.VerTope() })
	require.PanicsWithValue(t, "La pila esta vacia", func() { pila.Desapilar() })

	pila.Apilar("ja")
	require.Equal(t, "ja", pila.VerTope())
	require.False(t, pila.EstaVacia())
}