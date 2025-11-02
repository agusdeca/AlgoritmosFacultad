package diccionario_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	TDADiccionario "tdas/diccionario"

	"github.com/stretchr/testify/require"
)

func cmpStrings(a, b string) int {
	return strings.Compare(a, b)
}

func cmpInts(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

var TAMS_VOLUMEN_PRUEBA = 50000
var TAMS_VOLUMEN_ABB = []int{1000, 10000, 50000}

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
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, 1, dic.Obtener("uno"))

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
	require.False(t, dic.Pertenece("D"))

	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("D") })
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

func TestABBConClavesNumericas(t *testing.T) {
	t.Log("Valida que funcione con claves numéricas")
	dic := TDADiccionario.CrearABB[int, string](cmpInts)

	dic.Guardar(10, "diez")
	dic.Guardar(5, "cinco")
	dic.Guardar(15, "quince")

	require.EqualValues(t, 3, dic.Cantidad())
	require.True(t, dic.Pertenece(10))
	require.EqualValues(t, "diez", dic.Obtener(10))
	require.EqualValues(t, "cinco", dic.Borrar(5))
	require.False(t, dic.Pertenece(5))
}

func TestABBValorNulo(t *testing.T) {
	t.Log("Probamos que el valor puede ser nil sin problemas")
	dic := TDADiccionario.CrearABB[string, *int](cmpStrings)
	clave := "Pez"
	dic.Guardar(clave, nil)
	require.True(t, dic.Pertenece(clave))
	require.EqualValues(t, 1, dic.Cantidad())
	require.EqualValues(t, (*int)(nil), dic.Obtener(clave))
	require.EqualValues(t, (*int)(nil), dic.Borrar(clave))
	require.False(t, dic.Pertenece(clave))
}

func TestBorrarCasosEspecificos(t *testing.T) {
	t.Log("Verifica los tres casos de borrado: hoja, un hijo, dos hijos")
	// Árbol:
	//      D
	//    /   \
	//   B     F
	//  / \   / \
	// A   C E   G
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)
	dic.Guardar("D", 4)
	dic.Guardar("B", 2)
	dic.Guardar("F", 6)
	dic.Guardar("A", 1)
	dic.Guardar("C", 3)
	dic.Guardar("E", 5)
	dic.Guardar("G", 7)
	require.EqualValues(t, 7, dic.Cantidad())

	// Caso 1: Borrar Hoja (A)
	dic.Borrar("A")
	require.False(t, dic.Pertenece("A"))
	require.EqualValues(t, 6, dic.Cantidad())

	// Caso 2: Borrar nodo con 1 hijo (B, tiene a C)
	// Árbol debe quedar:
	//      D
	//    /   \
	//   C     F
	//        / \
	//       E   G
	dic.Borrar("B")
	require.False(t, dic.Pertenece("B"))
	require.True(t, dic.Pertenece("C"))
	require.EqualValues(t, 5, dic.Cantidad())

	// Caso 3: Borrar nodo con 2 hijos (D, raíz)
	// Sucesores 'E'.
	// Árbol debe quedar:
	//      E
	//    /   \
	//   C     F
	//          \
	//           G
	dic.Borrar("D")
	require.False(t, dic.Pertenece("D"))
	require.EqualValues(t, 4, dic.Cantidad())
	require.True(t, dic.Pertenece("E"))
	require.True(t, dic.Pertenece("F"))
	require.True(t, dic.Pertenece("C"))
	require.True(t, dic.Pertenece("G"))
}

func TestABBGuardarYBorrarRepetidasVeces(t *testing.T) {
	t.Log("Guardar y borrar repetidas veces verifica estabilidad del árbol")
	dic := TDADiccionario.CrearABB[int, int](cmpInts)

	n := 1000
	for i := 0; i < n; i++ {
		dic.Guardar(i, i)
		require.True(t, dic.Pertenece(i))
		require.EqualValues(t, i+1, dic.Cantidad())
	}

	for i := 0; i < n; i++ {
		require.EqualValues(t, i, dic.Borrar(i))
		require.False(t, dic.Pertenece(i))
		require.EqualValues(t, n-1-i, dic.Cantidad())
	}

	require.EqualValues(t, 0, dic.Cantidad())
}

func TestIteradorInternoOrdenado(t *testing.T) {
	t.Log("El iterador interno recorre elementos en orden")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	dic.Guardar("D", 4)
	dic.Guardar("B", 2)
	dic.Guardar("F", 6)
	dic.Guardar("A", 1)
	dic.Guardar("C", 3)

	claves := []string{}
	dic.Iterar(func(clave string, dato int) bool {
		claves = append(claves, clave)
		return true
	})

	esperado := []string{"A", "B", "C", "D", "F"}
	require.EqualValues(t, esperado, claves)
}

func TestIteradorInternoCorte(t *testing.T) {
	t.Log("El iterador interno debe detenerse cuando visitar devuelve false")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	for i := 0; i < 10; i++ {
		dic.Guardar(fmt.Sprintf("%d", i), i)
	}

	contador := 0
	dic.Iterar(func(_ string, _ int) bool {
		contador++
		return contador < 5
	})

	require.EqualValues(t, 5, contador)
}

func TestIterarRangoCompleto(t *testing.T) {
	t.Log("IterarRango con límites nil debe iterar todo")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := []string{"D", "B", "F", "A", "C"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}

	resultado := []string{}
	dic.IterarRango(nil, nil, func(clave string, _ int) bool {
		resultado = append(resultado, clave)
		return true
	})

	esperado := []string{"A", "B", "C", "D", "F"}
	require.EqualValues(t, esperado, resultado)
}

func TestIterarRangoSoloDesde(t *testing.T) {
	t.Log("IterarRango con solo desde definido")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := []string{"A", "B", "C", "D", "E", "F"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}

	desde := "C"
	resultado := []string{}
	dic.IterarRango(&desde, nil, func(clave string, _ int) bool {
		resultado = append(resultado, clave)
		return true
	})

	esperado := []string{"C", "D", "E", "F"}
	require.EqualValues(t, esperado, resultado)
}

func TestIterarRangoSoloHasta(t *testing.T) {
	t.Log("IterarRango con solo hasta definido")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := []string{"A", "B", "C", "D", "E", "F"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}

	hasta := "D"
	resultado := []string{}
	dic.IterarRango(nil, &hasta, func(clave string, _ int) bool {
		resultado = append(resultado, clave)
		return true
	})

	esperado := []string{"A", "B", "C", "D"}
	require.EqualValues(t, esperado, resultado)
}

func TestIteradorExternoTrasBorrados(t *testing.T) {
	t.Log("Iterador creado tras borrar elementos no los incluye")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := []string{"A", "B", "C", "D", "E"}
	for i, k := range claves {
		dic.Guardar(k, i)
	}

	dic.Borrar("B")
	dic.Borrar("D")

	iter := dic.Iterador()
	resultado := []string{}
	for iter.HaySiguiente() {
		clave, _ := iter.VerActual()
		resultado = append(resultado, clave)
		iter.Siguiente()
	}

	esperado := []string{"A", "C", "E"}
	require.EqualValues(t, esperado, resultado)
}

func TestMultiplesIteradores(t *testing.T) {
	t.Log("Múltiples iteradores sobre el mismo ABB funcionan independientemente")
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	dic.Guardar("A", 1)
	dic.Guardar("B", 2)
	dic.Guardar("C", 3)

	iter1 := dic.Iterador()
	iter2 := dic.Iterador()

	iter1.Siguiente()
	iter1.Siguiente()

	clave, _ := iter2.VerActual()
	require.EqualValues(t, "A", clave)
	require.True(t, iter2.HaySiguiente())
}

func ejecutarPruebaVolumenABB(tb testing.TB, n int, ordenado bool) {
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := make([]string, n)
	valores := make([]int, n)

	for i := 0; i < n; i++ {
		valores[i] = i
		claves[i] = fmt.Sprintf("%08d", i)
	}

	if !ordenado {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(claves), func(i, j int) {
			claves[i], claves[j] = claves[j], claves[i]
			valores[i], valores[j] = valores[j], valores[i]
		})
	}

	for i := 0; i < n; i++ {
		dic.Guardar(claves[i], valores[i])
	}

	require.EqualValues(tb, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	ok := true
	for i := 0; i < n; i++ {
		ok = dic.Pertenece(claves[i])
		if !ok {
			break
		}
		ok = dic.Obtener(claves[i]) == valores[i]
		if !ok {
			break
		}
	}
	require.True(tb, ok, "Pertenece y Obtener con muchos elementos no funciona correctamente")

	for i := 0; i < n; i++ {
		ok = dic.Borrar(claves[i]) == valores[i]
		if !ok {
			break
		}
		ok = !dic.Pertenece(claves[i])
		if !ok {
			break
		}
	}
	require.True(tb, ok, "Borrar muchos elementos no funciona correctamente")
	require.EqualValues(tb, 0, dic.Cantidad())
}

func TestVolumenABBOrdenado(t *testing.T) {
	t.Log("Prueba de volumen (peor caso, inserción ordenada)")
	n := TAMS_VOLUMEN_PRUEBA
	ejecutarPruebaVolumenABB(t, n, true)
}

func TestVolumenABBAleatorio(t *testing.T) {
	t.Log("Prueba de volumen (caso promedio, inserción aleatoria)")
	n := TAMS_VOLUMEN_PRUEBA
	ejecutarPruebaVolumenABB(t, n, false)
}

func ejecutarPruebasVolumenIteradorABB(tb testing.TB, n int, ordenado bool) {
	dic := TDADiccionario.CrearABB[string, *int](cmpStrings)

	claves := make([]string, n)
	valores := make([]int, n)

	for i := 0; i < n; i++ {
		claves[i] = fmt.Sprintf("%08d", i)
		valores[i] = i
	}

	if !ordenado {
		// Desordenamos las claves
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(claves), func(i, j int) {
			claves[i], claves[j] = claves[j], claves[i]
			valores[i], valores[j] = valores[j], valores[i]
		})
	}

	for i := 0; i < n; i++ {
		dic.Guardar(claves[i], &valores[i])
	}

	// Prueba de iteración
	iter := dic.Iterador()
	require.True(tb, iter.HaySiguiente())

	ok := true
	var i int
	var clave string
	var valor *int

	for i = 0; iter.HaySiguiente(); i++ {
		c1, v1 := iter.VerActual()
		clave = c1
		if clave == "" {
			ok = false
			break
		}
		valor = v1
		if valor == nil {
			ok = false
			break
		}
		*valor = n
		iter.Siguiente()
	}

	require.True(tb, ok, "Iteracion en volumen no funciona correctamente")
	require.EqualValues(tb, n, i, "No se recorrió todo el largo")
	require.False(tb, iter.HaySiguiente(), "El iterador debe estar al final")

	ok = true
	for i = 0; i < n; i++ {
		if valores[i] != n {
			ok = false
			break
		}
	}
	require.True(tb, ok, "No se cambiaron todos los elementos")
}

func TestVolumenIteradorOrdenado(t *testing.T) {
	t.Log("Prueba de volumen del iterador (peor caso, inserción ordenada)")
	n := TAMS_VOLUMEN_PRUEBA
	ejecutarPruebasVolumenIteradorABB(t, n, true)
}

func TestVolumenIteradorAleatorio(t *testing.T) {
	t.Log("Prueba de volumen del iterador (caso promedio, inserción aleatoria)")
	n := TAMS_VOLUMEN_PRUEBA
	ejecutarPruebasVolumenIteradorABB(t, n, false)
}

func TestABBVolumenIteradorCorte(t *testing.T) {
	t.Log("Prueba de volumen de iterador interno, validando que siempre que se indique corte, se corte")

	dic := TDADiccionario.CrearABB[int, int](cmpInts)
	n := 50000
	for i := 0; i < n; i++ {
		dic.Guardar(i, i)
	}

	seguirEjecutando := true
	siguioEjecutandoCuandoNoDebia := false
	corteEn := n / 10
	contador := 0

	dic.Iterar(func(c int, v int) bool {
		if !seguirEjecutando {
			siguioEjecutandoCuandoNoDebia = true
			return false
		}
		if c == corteEn {
			seguirEjecutando = false
			return false
		}
		contador++
		return true
	})

	require.False(t, seguirEjecutando, "Se tendría que haber encontrado un elemento que genere el corte")
	require.False(t, siguioEjecutandoCuandoNoDebia, "No debería haber seguido ejecutando después del corte")
	require.EqualValues(t, corteEn, contador, "El contador no coincide con el punto de corte")
}
