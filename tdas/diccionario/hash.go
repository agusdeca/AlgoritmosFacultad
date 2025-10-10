package diccionario

import "fmt"

type estadoCelda int

const (
	CAPACIDAD_INICIAL                  = 17
	FACTOR_CARGA_MAX                   = 0.7
	FACTOR_CARGA_MIN                   = 0.2
	FACTOR_REDIMENSION                 = 2
	CAPACIDAD_MINIMA                   = 17
	MENSAJE_CLAVE_INEXIST              = "La clave no pertenece al diccionario"
	MENSAJE_ITER_TERMINADO             = "El iterador termino de iterar"
	VACIO                  estadoCelda = iota
	OCUPADO
	BORRADO
)

type celda[K any, V any] struct {
	clave  K
	valor  V
	estado estadoCelda
}

type hashCerrado[K any, V any] struct {
	tabla     []*celda[K, V]
	capacidad int
	cantidad  int
	borrados  int
	igualdad  func(K, K) bool
}

type iterHash[K any, V any] struct {
	hash     *hashCerrado[K, V]
	posicion int
}

func convertirABytes[K any](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

// MurmurHash3 (32-bit)
// Fuente original: https://github.com/aappleby/smhasher/blob/master/src/MurmurHash3.cpp
func funcionHash[K any](clave K, capacidad int) int {
	bytes := convertirABytes(clave)

	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
		r1 uint32 = 15
		r2 uint32 = 13
		m  uint32 = 5
		n  uint32 = 0xe6546b64
	)

	seed := uint32(0)
	hash := seed
	length := len(bytes)

	nblocks := length / 4
	for i := 0; i < nblocks; i++ {
		k := uint32(bytes[i*4]) | uint32(bytes[i*4+1])<<8 |
			uint32(bytes[i*4+2])<<16 | uint32(bytes[i*4+3])<<24

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	tail := bytes[nblocks*4:]
	var k1 uint32
	switch len(tail) {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1
		k1 = (k1 << r1) | (k1 >> (32 - r1))
		k1 *= c2
		hash ^= k1
	}

	hash ^= uint32(length)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return int(hash) % capacidad
}

func (h *hashCerrado[K, V]) crearTabla(capacidad int) {
	h.tabla = make([]*celda[K, V], capacidad)
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
		if c != nil && c.estado == OCUPADO {
			h.Guardar(c.clave, c.valor)
		}
	}
}

func (h *hashCerrado[K, V]) buscar(clave K) (int, bool) {
	pos := funcionHash(clave, h.capacidad)
	inicio := pos

	for {
		if h.tabla[pos] == nil {
			return pos, false
		}

		if h.tabla[pos].estado == OCUPADO && h.igualdad(h.tabla[pos].clave, clave) {
			return pos, true
		}

		pos = (pos + 1) % h.capacidad
		if pos == inicio {
			return pos, false
		}
	}
}

func (h *hashCerrado[K, V]) buscarParaInsertar(clave K) (int, bool) {
	pos := funcionHash(clave, h.capacidad)
	inicio := pos
	primerBorrado := -1
	claveExiste := false

	for {
		if h.tabla[pos] == nil {
			break
		}

		if h.tabla[pos].estado == OCUPADO && h.igualdad(h.tabla[pos].clave, clave) {
			return pos, true
		}

		if h.tabla[pos].estado == BORRADO && primerBorrado == -1 {
			primerBorrado = pos
		}

		pos = (pos + 1) % h.capacidad
		if pos == inicio {
			break
		}
	}

	if primerBorrado != -1 {
		return primerBorrado, claveExiste
	}
	return pos, claveExiste
}

func (h *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if h.debeAgrandar() {
		h.redimensionar(h.capacidad * FACTOR_REDIMENSION)
	}

	pos, existe := h.buscarParaInsertar(clave)

	if existe {
		h.tabla[pos].valor = valor
	} else {
		if h.tabla[pos] != nil && h.tabla[pos].estado == BORRADO {
			h.borrados--
		}
		h.tabla[pos] = &celda[K, V]{
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

// Iterador interno
func (h *hashCerrado[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for _, c := range h.tabla {
		if c != nil && c.estado == OCUPADO {
			if !visitar(c.clave, c.valor) {
				return
			}
		}
	}
}

// Iterador externo
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
		if it.hash.tabla[i] != nil && it.hash.tabla[i].estado == OCUPADO {
			it.posicion = i
			return
		}
	}
	it.posicion = -1
}