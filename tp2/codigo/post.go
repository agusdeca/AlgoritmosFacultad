package codigo

import (
	"fmt"
	"strings"
	tdadict "tdas/diccionario"
)

type post struct {
	id    int
	autor *usuario
	texto string
	likes tdadict.Diccionario[string, bool]
}

func nuevoPost(id int, autor *usuario, texto string) *post {
	return &post{
		id:    id,
		autor: autor,
		texto: texto,
		likes: tdadict.CrearHash[string, bool](func(a, b string) bool { return a == b }),
	}
}

// darLike agrega el like de un usuario, si no existe ya
func (p *post) darLike(nombre string) {
	if !p.likes.Pertenece(nombre) {
		p.likes.Guardar(nombre, true)
	}
}

// cantidadLikes devuelve la cant de usuarios que dieron like
func (p *post) cantidadLikes() int {
	return p.likes.Cantidad()
}

// obtenerUsuariosLikes devuelve una lista de nombres por orden alfabetico
func (p *post) obtenerUsuariosLikes() []string {
	nombres := make([]string, 0, p.likes.Cantidad())
	likes := tdadict.CrearABB[string, bool](strings.Compare)
	p.likes.Iterar(func(clave string, valor bool) bool {
		likes.Guardar(clave, valor)
		return true
	})
	iter := likes.Iterador()
	for iter.HaySiguiente() {
		nombre, _ := iter.VerActual()
		nombres = append(nombres, nombre)
		iter.Siguiente()
	}
	return nombres
}

// Imprimir
func (p *post) String() string {
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d",
		p.id, p.autor.nombre, p.texto, p.cantidadLikes())
}
