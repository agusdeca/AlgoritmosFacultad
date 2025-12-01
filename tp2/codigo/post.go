package codigo

import (
	"fmt"
	"strings"
	tdadict "tdas/diccionario"
)

type Post interface {
	ID() int
	Autor() string
	Texto() string
	DarLike(nombre string)
	CantidadLikes() int
	MostrarLikes() string
	String() string
}

type post struct {
	id    int
	autor *usuario
	texto string
	likes tdadict.Diccionario[string, bool]
}

func nuevoPost(id int, autor *usuario, texto string) Post {
	return &post{
		id:    id,
		autor: autor,
		texto: texto,
		likes: tdadict.CrearABB[string, bool](strings.Compare),
	}
}

func (p *post) ID() int {
	return p.id
}

func (p *post) Autor() string {
	return p.autor.nombre
}
func (p *post) Texto() string {
	return p.texto
}

func (p *post) DarLike(nombre string) {
	if !p.likes.Pertenece(nombre) {
		p.likes.Guardar(nombre, true)
	}
}

func (p *post) CantidadLikes() int { return p.likes.Cantidad() }

func (p *post) MostrarLikes() string {
	if p.likes.Cantidad() == 0 {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("El post tiene %d likes:\n", p.likes.Cantidad()))

	iter := p.likes.Iterador()
	for iter.HaySiguiente() {
		nombre, _ := iter.VerActual()
		builder.WriteString(fmt.Sprintf("\t%s\n", nombre))
		iter.Siguiente()
	}

	s := builder.String()
	return strings.TrimSuffix(s, "\n")
}

func (p *post) String() string {
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d",
		p.id, p.autor.nombre, p.texto, p.CantidadLikes())
}
