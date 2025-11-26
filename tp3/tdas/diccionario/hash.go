package diccionario

import (
	"fmt"
)

const (
	CAPACIDAD_INICIAL      = 127
	FACTOR_CARGA_MAX       = 0.7
	FACTOR_CARGA_MIN       = 0.2
	FACTOR_REDIMENSION     = 2
	CAPACIDAD_MINIMA       = 127
	MENSAJE_CLAVE_INEXIST  = "La clave no pertenece al diccionario"
	MENSAJE_ITER_TERMINADO = "El iterador termino de iterar"
)

type estadoCelda int

const (
	VACIO   estadoCelda = iota // VACIO = 0
	OCUPADO                    // OCUPADO = 1
	BORRADO                    // BORRADO = 2
)

type celda[K any, V any] struct {
	clave  K
	valor  V
	estado estadoCelda
}

type hashCerrado[K any, V any] struct {
	tabla     []celda[K, V]
	capacidad int
	cantidad  int
	borrados  int
	igualdad  func(K, K) bool
}

type iterHash[K any, V any] struct {
	hash     *hashCerrado[K, V]
	posicion int
}

// FNV-1a Hash (Optimizado)
func funcionHash[K any](clave K, capacidad int) int {
	const offset64 = 14695981039346656037
	const prime64 = 1099511628211

	var hash uint64 = offset64

	if s, ok := any(clave).(string); ok {
		for i := 0; i < len(s); i++ {
			hash ^= uint64(s[i])
			hash *= prime64
		}
	} else {
		s := fmt.Sprintf("%v", clave)
		for i := 0; i < len(s); i++ {
			hash ^= uint64(s[i])
			hash *= prime64
		}
	}

	idx := int(hash % uint64(capacidad))
	if idx < 0 {
		idx += capacidad
	}
	return idx
}

func (h *hashCerrado[K, V]) crearTabla(capacidad int) {
	h.tabla = make([]celda[K, V], capacidad)
	h.capacidad = capacidad
	h.cantidad = 0
	h.borrados = 0
}

func CrearHash[K any, V any](igualdad func(K, K) bool) Diccionario[K, V] {
	h := &hashCerrado[K, V]{igualdad: igualdad}
	h.crearTabla(CAPACIDAD_INICIAL)
	return h
}

func (h *hashCerrado[K, V]) factorCarga() float64 {
	return float64(h.cantidad+h.borrados) / float64(h.capacidad)
}

func (h *hashCerrado[K, V]) debeAgrandar() bool {
	return h.factorCarga() > FACTOR_CARGA_MAX
}

func (h *hashCerrado[K, V]) debeAchicar() bool {
	return h.capacidad > CAPACIDAD_MINIMA && h.factorCarga() < FACTOR_CARGA_MIN
}

func (h *hashCerrado[K, V]) redimensionar(nuevaCapacidad int) {
	tablaVieja := h.tabla
	h.crearTabla(nuevaCapacidad)

	for _, c := range tablaVieja {
		if c.estado == OCUPADO {
			h.Guardar(c.clave, c.valor)
		}
	}
}

func (h *hashCerrado[K, V]) buscar(clave K) (int, bool) {
	pos := funcionHash(clave, h.capacidad)
	inicio := pos

	for h.tabla[pos].estado != VACIO {
		if h.tabla[pos].estado == OCUPADO && h.igualdad(h.tabla[pos].clave, clave) {
			return pos, true
		}
		pos = (pos + 1) % h.capacidad
		if pos == inicio {
			break
		}
	}
	return pos, false
}

func (h *hashCerrado[K, V]) buscarParaInsertar(clave K) (int, bool) {
	pos := funcionHash(clave, h.capacidad)
	inicio := pos
	primerBorrado := -1

	for h.tabla[pos].estado != VACIO {
		if h.tabla[pos].estado == OCUPADO && h.igualdad(h.tabla[pos].clave, clave) {
			return pos, true
		}
		if h.tabla[pos].estado == BORRADO && primerBorrado == -1 {
			primerBorrado = pos
		}
		pos = (pos + 1) % h.capacidad
		
		if pos == inicio {
			if primerBorrado != -1 {
				return primerBorrado, false
			}
			return -1, false 
		}
	}

	if primerBorrado != -1 {
		return primerBorrado, false
	}
	return pos, false
}

func (h *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if h.debeAgrandar() {
		h.redimensionar(h.capacidad * FACTOR_REDIMENSION)
	}

	pos, existe := h.buscarParaInsertar(clave)

	// SAFETY CHECK: Recuperación si la tabla estaba llena
	if pos == -1 {
		h.redimensionar(h.capacidad * FACTOR_REDIMENSION)
		pos, existe = h.buscarParaInsertar(clave)
	}

	if existe {
		h.tabla[pos].valor = valor
	} else {
		if h.tabla[pos].estado == BORRADO {
			h.borrados--
		}
		h.tabla[pos] = celda[K, V]{
			clave:  clave,
			valor:  valor,
			estado: OCUPADO,
		}
		h.cantidad++
	}
}

func (h *hashCerrado[K, V]) Pertenece(clave K) bool {
	_, existe := h.buscar(clave)
	return existe
}

func (h *hashCerrado[K, V]) Obtener(clave K) V {
	pos, existe := h.buscar(clave)
	if !existe {
		panic(MENSAJE_CLAVE_INEXIST)
	}
	return h.tabla[pos].valor
}

func (h *hashCerrado[K, V]) Borrar(clave K) V {
	pos, existe := h.buscar(clave)
	if !existe {
		panic(MENSAJE_CLAVE_INEXIST)
	}

	valor := h.tabla[pos].valor
	h.tabla[pos].estado = BORRADO
	h.cantidad--
	h.borrados++

	if h.debeAchicar() {
		nuevaCap := h.capacidad / FACTOR_REDIMENSION
		if nuevaCap < CAPACIDAD_MINIMA {
			nuevaCap = CAPACIDAD_MINIMA
		}
		h.redimensionar(nuevaCap)
	}

	return valor
}

func (h *hashCerrado[K, V]) Cantidad() int {
	return h.cantidad
}

func (h *hashCerrado[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for _, c := range h.tabla {
		if c.estado == OCUPADO {
			if !visitar(c.clave, c.valor) {
				return
			}
		}
	}
}

func (h *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	iter := &iterHash[K, V]{hash: h, posicion: -1}
	iter.avanzar()
	return iter
}

func (it *iterHash[K, V]) HaySiguiente() bool {
	return it.posicion < it.hash.capacidad && it.posicion != -1
}

func (it *iterHash[K, V]) VerActual() (K, V) {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITER_TERMINADO)
	}
	celda := it.hash.tabla[it.posicion]
	return celda.clave, celda.valor
}

func (it *iterHash[K, V]) Siguiente() {
	if !it.HaySiguiente() {
		panic(MENSAJE_ITER_TERMINADO)
	}
	it.avanzar()
}

func (it *iterHash[K, V]) avanzar() {
	for i := it.posicion + 1; i < it.hash.capacidad; i++ {
		if it.hash.tabla[i].estado == OCUPADO {
			it.posicion = i
			return
		}
	}
	it.posicion = -1
}