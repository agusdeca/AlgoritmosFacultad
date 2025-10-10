package diccionario_test

import (
	"strings"
	"testing"

	TDADiccionario "tdas/diccionario"

	"github.com/stretchr/testify/require"
)

func cmpStrings(a, b string) int {
	return strings.Compare(a, b)
}

func TestAbbVacio(t *testing.T) {
	t.Log("Un ABB vacío no tiene claves ni elementos")

	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("A") })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("A") })
}

func TestGuardarYObtener(t *testing.T) {
	t.Log("Guardar y obtener elementos correctamente")

	dic := TDADiccionario.CrearABB[string, string](cmpStrings)
	dic.Guardar("B", "b")
	dic.Guardar("A", "a")
	dic.Guardar("C", "c")

	require.EqualValues(t, 3, dic.Cantidad())
	require.True(t, dic.Pertenece("B"))
	require.EqualValues(t, "b", dic.Obtener("B"))
}

func TestReemplazarDato(t *testing.T) {
	t.Log("Reemplaza el dato si la clave ya existía")

	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	dic.Guardar("uno", 1)
	dic.Guardar("uno", 100)
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, 100, dic.Obtener("uno"))
}

func TestBorrarElementos(t *testing.T) {
	t.Log("Borra claves correctamente en distintos casos")

	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	claves := []string{"D", "B", "A", "C", "F", "E", "G"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}
	require.EqualValues(t, 7, dic.Cantidad())

	require.EqualValues(t, 0, dic.Borrar("D"))
	require.EqualValues(t, 6, dic.Cantidad())

	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("Z") })
}

func TestIteradorExternoVacio(t *testing.T) {
	t.Log("Iterador sobre ABB vacío")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	iter := dic.Iterador()
	require.False(t, iter.HaySiguiente())
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.VerActual() })
	require.PanicsWithValue(t, "El iterador termino de iterar", func() { iter.Siguiente() })
}

func TestIteradorExternoOrdenado(t *testing.T) {
	t.Log("Itera elementos en orden")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	claves := []string{"D", "B", "F", "A", "C", "E", "G"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}
	
	iter := dic.Iterador()
	esperado := []string{"A", "B", "C", "D", "E", "F", "G"}
	i := 0
	
	for iter.HaySiguiente() {
		clave, _ := iter.VerActual()
		require.EqualValues(t, esperado[i], clave)
		iter.Siguiente()
		i++
	}
	require.EqualValues(t, len(esperado), i)
}

func TestIteradorRango(t *testing.T) {
	t.Log("Itera solo en el rango especificado")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	claves := []string{"A", "B", "C", "D", "E", "F", "G"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}
	
	desde := "C"
	hasta := "F"
	iter := dic.IteradorRango(&desde, &hasta)
	
	esperado := []string{"C", "D", "E", "F"}
	i := 0
	
	for iter.HaySiguiente() {
		clave, _ := iter.VerActual()
		require.EqualValues(t, esperado[i], clave)
		iter.Siguiente()
		i++
	}
	require.EqualValues(t, len(esperado), i)
}
