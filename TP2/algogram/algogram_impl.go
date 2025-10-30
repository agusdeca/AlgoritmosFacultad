package algogram

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"tp2/tda"
)

type algogramImpl struct {
	usuarios_hash  tda.CrearHash[string, *usuario] 
	usuarios_lista []*usuario
	usuarios_cantidad int 
	posts tda.CrearHash[int, *post]
	usuario_loggeado *usuario
	proximo_id_post int
}

func (ag *algogramImpl) CargarUsuarios(ruta string) error {
	archivo, err := os.Open(ruta)
	if err != nil {
		return err
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	nro_linea := 0
	for scanner.Scan() {
		nombre := scanner.Text()
		nuevoUsuario := newUsuario(nombre, nro_linea) 
		
		if ag.usuarios_cantidad == cap(ag.usuarios_lista) {
			
			nuevaCapacidad := cap(ag.usuarios_lista) * 2
			nuevoSlice := make([]*usuario, nuevaCapacidad)
			copy(nuevoSlice, ag.usuarios_lista)
			ag.usuarios_lista = nuevoSlice
		}
		
		ag.usuarios_lista[ag.usuarios_cantidad] = nuevoUsuario
		ag.usuarios_cantidad++
		
		ag.usuarios_hash.Guardar(nombre, nuevoUsuario)
		nro_linea++
	}
	return scanner.Err()
}

func (ag *algogramImpl) Login(nombre string) string {
	if ag.usuario_loggeado != nil {
		return "Error: Ya habia un usuario logeado"
	}
	usuario, existe := ag.usuarios_hash.Obtener(nombre)
	if !existe {
		return "Error: usuario no existente"
	}
	ag.usuario_loggeado = usuario
	return fmt.Sprintf("Hola %s", nombre)
}

func (ag *algogramImpl) Logout() string {
	if ag.usuario_loggeado == nil {
		return "Error: no habia usuario loggeado"
	}
	ag.usuario_loggeado = nil
	return "Adios"
}

func (ag *algogramImpl) Publicar(texto string) string {
	if ag.usuario_loggeado == nil {
		return "Error: no habia usuario loggeado"
	}

	idActual := ag.proximo_id_post
	autor := ag.usuario_loggeado
	nuevoPost := newPost(idActual, autor, texto)

	ag.posts.Guardar(idActual, nuevoPost)
	ag.proximo_id_post++

	for i := 0; i < ag.usuarios_cantidad; i++ {
		u := ag.usuarios_lista[i]
		
		if u == autor {
			continue 
		}
		
		afinidad := int(math.Abs(float64(autor.nro_linea - u.nro_linea)))
		
		entradaFeed := postFeed{
			post:      nuevoPost,
			prioridad: afinidad,
			id_post:   idActual,
		}
		u.feed.Encolar(entradaFeed)
	}
	return "Post publicado"
}

func (ag *algogramImpl) VerSiguienteFeed() string {
	if ag.usuario_loggeado == nil {
		return "Usuario no loggeado o no hay mas posts para ver"
	}
	postFeed, err := ag.usuario_loggeado.obtenerProximoPost()
	if err != nil { 
		return "Usuario no loggeado o no hay mas posts para ver"
	}
	p := postFeed.post
	return fmt.Sprintf("Post ID %d\n%s dijo: %s\nLikes: %d",
		p.id,
		p.autor.nombre,
		p.texto,
		p.cant_likes)
}

func (ag *algogramImpl) LikearPost(id int) string {
	if ag.usuario_loggeado == nil {
		return "Error: Usuario no loggeado o Post inexistente"
	}
	post, existe := ag.posts.Obtener(id)
	if !existe {
		return "Error: Usuario no loggeado o Post inexistente"
	}
	post.darLike(ag.usuario_loggeado.nombre)
	return "Post likeado"
}

func (ag *algogramImpl) MostrarLikes(id int) string {
	post, existe := ag.posts.Obtener(id)
	if !existe || post.cant_likes == 0 {
		return "Error: Post inexistente o sin likes"
	}
	nombres := post.obtenerLikes()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("El post tiene %d likes:\n", post.cant_likes))
	for _, nombre := range nombres { 
		builder.WriteString(fmt.Sprintf("\t%s\n", nombre))
	}
	return strings.TrimSuffix(builder.String(), "\n")
}