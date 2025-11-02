package codigo

import (
	"sort"
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
	p.likes.Iterar(func(clave string, _ bool) bool {
		nombres = append(nombres, clave)
		return true
	})
	sort.Strings(nombres)
	return nombres
}
