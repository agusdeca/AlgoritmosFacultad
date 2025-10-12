package diccionario_test

import (
	"fmt"
	"strings"
	"testing"

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

var TAMS_VOLUMEN_ABB = []int{1000, 5000, 10000, 20000}

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
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	// Construir árbol específico
	dic.Guardar("D", 4)
	dic.Guardar("B", 2)
	dic.Guardar("F", 6)
	dic.Guardar("A", 1)
	dic.Guardar("C", 3)
	dic.Guardar("E", 5)
	dic.Guardar("G", 7)

	// Borrar hoja
	dic.Borrar("A")
	require.False(t, dic.Pertenece("A"))
	require.EqualValues(t, 6, dic.Cantidad())

	// Borrar nodo con un hijo
	dic.Borrar("B")
	require.False(t, dic.Pertenece("B"))
	require.True(t, dic.Pertenece("C"))

	// Borrar nodo con dos hijos
	dic.Borrar("D")
	require.False(t, dic.Pertenece("D"))
	require.True(t, dic.Pertenece("E"))
	require.True(t, dic.Pertenece("F"))
}

func TestABBGuardarYBorrarRepetidasVeces(t *testing.T) {
	t.Log("Guardar y borrar repetidas veces verifica estabilidad del árbol")
	dic := TDADiccionario.CrearABB[int, int](cmpInts)

	for i := 0; i < 1000; i++ {
		dic.Guardar(i, i)
		require.True(t, dic.Pertenece(i))
		dic.Borrar(i)
		require.False(t, dic.Pertenece(i))
	}
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

	// Avanzar el primer iterador
	iter1.Siguiente()
	iter1.Siguiente()

	// El segundo no debería haberse movido
	clave, _ := iter2.VerActual()
	require.EqualValues(t, "A", clave)
	require.True(t, iter2.HaySiguiente())
}

func ejecutarPruebaVolumenABB(b *testing.B, n int) {
	dic := TDADiccionario.CrearABB[string, int](cmpStrings)

	claves := make([]string, n)
	valores := make([]int, n)

	// Insertar n elementos
	for i := 0; i < n; i++ {
		valores[i] = i
		claves[i] = fmt.Sprintf("%08d", i)
		dic.Guardar(claves[i], valores[i])
	}

	require.EqualValues(b, n, dic.Cantidad(), "La cantidad de elementos es incorrecta")

	// Verificar Pertenece y Obtener
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
	require.True(b, ok, "Pertenece y Obtener con muchos elementos no funciona correctamente")

	// Verificar Borrar
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
	require.True(b, ok, "Borrar muchos elementos no funciona correctamente")
	require.EqualValues(b, 0, dic.Cantidad())
}

func BenchmarkABB(b *testing.B) {
	b.Log("Prueba de stress del ABB. Prueba guardando distinta cantidad de elementos (muy grandes), " +
		"ejecutando muchas veces las pruebas para generar un benchmark")
	for _, n := range TAMS_VOLUMEN_ABB {
		b.Run(fmt.Sprintf("Prueba %d elementos", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ejecutarPruebaVolumenABB(b, n)
			}
		})
	}
}

func ejecutarPruebasVolumenIteradorABB(b *testing.B, n int) {
	dic := TDADiccionario.CrearABB[string, *int](cmpStrings)

	claves := make([]string, n)
	valores := make([]int, n)

	for i := 0; i < n; i++ {
		claves[i] = fmt.Sprintf("%08d", i)
		valores[i] = i
		dic.Guardar(claves[i], &valores[i])
	}

	// Prueba de iteración
	iter := dic.Iterador()
	require.True(b, iter.HaySiguiente())

	ok := true
	var i int
	var clave string
	var valor *int

	for i = 0; i < n; i++ {
		if !iter.HaySiguiente() {
			ok = false
			break
		}
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

	require.True(b, ok, "Iteracion en volumen no funciona correctamente")
	require.EqualValues(b, n, i, "No se recorrió todo el largo")
	require.False(b, iter.HaySiguiente(), "El iterador debe estar al final")

	ok = true
	for i = 0; i < n; i++ {
		if valores[i] != n {
			ok = false
			break
		}
	}
	require.True(b, ok, "No se cambiaron todos los elementos")
}

func BenchmarkIteradorABB(b *testing.B) {
	b.Log("Prueba de stress del Iterador del ABB")
	for _, n := range TAMS_VOLUMEN_ABB {
		b.Run(fmt.Sprintf("Prueba %d elementos", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ejecutarPruebasVolumenIteradorABB(b, n)
			}
		})
	}
}

func TestABBVolumenIteradorCorte(t *testing.T) {
	t.Log("Prueba de volumen de iterador interno, validando que siempre que se indique corte, se corte")

	dic := TDADiccionario.CrearABB[int, int](cmpInts)

	for i := 0; i < 10000; i++ {
		dic.Guardar(i, i)
	}

	seguirEjecutando := true
	siguioEjecutandoCuandoNoDebia := false

	dic.Iterar(func(c int, v int) bool {
		if !seguirEjecutando {
			siguioEjecutandoCuandoNoDebia = true
			return false
		}
		if c%100 == 0 && c > 0 {
			seguirEjecutando = false
			return false
		}
		return true
	})

	require.False(t, seguirEjecutando, "Se tendría que haber encontrado un elemento que genere el corte")
	require.False(t, siguioEjecutandoCuandoNoDebia, "No debería haber seguido ejecutando después del corte")
}
