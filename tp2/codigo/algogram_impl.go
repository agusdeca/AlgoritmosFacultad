package codigo

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"tdas/diccionario"
)

type algogramImpl struct {
	usuarios_hash    diccionario.Diccionario[string, *usuario]
	usuarios_lista   []*usuario
	posts            diccionario.Diccionario[int, *post]
	usuario_loggeado *usuario
	proximo_id_post  int
}

func CrearAlgogram() *algogramImpl {
	return &algogramImpl{
		usuarios_hash:  diccionario.CrearHash[string, *usuario](func(a, b string) bool { return a == b }),
		usuarios_lista: []*usuario{},
		posts:          diccionario.CrearHash[int, *post](func(a, b int) bool { return a == b }),
	}
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
		nuevoUsuario := nuevoUsuario(nombre, nro_linea)
		ag.usuarios_lista = append(ag.usuarios_lista, nuevoUsuario)
		ag.usuarios_hash.Guardar(nombre, nuevoUsuario)
		nro_linea++
	}
	return scanner.Err()
}

func (ag *algogramImpl) Login(nombre string) string {
	if ag.usuario_loggeado != nil {
		return ERR_USUARIO_EXISTE
	}
	if !ag.usuarios_hash.Pertenece(nombre) {
		return ERR_USUARIO_NO_EXISTE
	}
	usuario := ag.usuarios_hash.Obtener(nombre)
	ag.usuario_loggeado = usuario
	return fmt.Sprintf(MSJ_HOLA, nombre)
}

func (ag *algogramImpl) Logout() string {
	if ag.usuario_loggeado == nil {
		return ERR_NO_LOGGEADO
	}
	ag.usuario_loggeado = nil
	return MSJ_ADIOS
}

func (ag *algogramImpl) Publicar(texto string) string {
	if ag.usuario_loggeado == nil {
		return ERR_NO_LOGGEADO
	}

	idActual := ag.proximo_id_post
	autor := ag.usuario_loggeado
	nuevoPost := nuevoPost(idActual, autor, texto)
	ag.posts.Guardar(idActual, nuevoPost)
	ag.proximo_id_post++

	for _, u := range ag.usuarios_lista {
		if u == autor {
			continue
		}

		afinidad := int(math.Abs(float64(autor.posicion - u.posicion)))
		entradaFeed := entradaFeed{
			post:      nuevoPost,
			prioridad: afinidad,
			id_post:   idActual,
		}
		u.feed.Encolar(entradaFeed)
	}
	return MSJ_POST_PUBLICADO
}

func (ag *algogramImpl) VerSiguienteFeed() string {
	if ag.usuario_loggeado == nil {
		return ERR_LOGUEADO_POST
	}

	post := ag.usuario_loggeado.proximoPost()
	if post == nil {
		return ERR_LOGUEADO_POST
	}

	return post.String()
}

func (ag *algogramImpl) LikearPost(id int) string {
	if ag.usuario_loggeado == nil {
		return ERR_POST_INEXISTENTE
	}
	if !ag.posts.Pertenece(id) {
		return ERR_POST_INEXISTENTE
		// No las uno en un or porque Obtener() me da panic si la clave no existe
	}

	post := ag.posts.Obtener(id)
	post.darLike(ag.usuario_loggeado.nombre)
	return "Post likeado"
}

func (ag *algogramImpl) MostrarLikes(id int) string {
	if !ag.posts.Pertenece(id) {
		return ERR_SIN_LIKES
	}
	post := ag.posts.Obtener(id)
	if post.cantidadLikes() == 0 {
		return ERR_SIN_LIKES
	}
	nombres := post.obtenerUsuariosLikes()
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("El post tiene %d likes:\n", post.cantidadLikes()))
	for _, nombre := range nombres {
		builder.WriteString(fmt.Sprintf("\t%s\n", nombre))
	}
	return strings.TrimSuffix(builder.String(), "\n")
}
