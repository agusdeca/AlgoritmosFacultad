package codigo

import (
	"fmt"
	"tdas/diccionario"
)

type algogramImpl struct {
	usuarios_hash    diccionario.Diccionario[string, Usuario]
	posts            diccionario.Diccionario[int, Post]
	usuario_loggeado Usuario
	proximo_id_post  int
}

func (ag *algogramImpl) Login(nombre string) string {
	if ag.usuario_loggeado != nil {
		return ERR_USUARIO_EXISTE
	}
	if !ag.usuarios_hash.Pertenece(nombre) {
		return ERR_USUARIO_NO_EXISTE
	}
	ag.usuario_loggeado = ag.usuarios_hash.Obtener(nombre)
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
	nuevo := nuevoPost(idActual, ag.usuario_loggeado.(*usuario), texto)
	ag.posts.Guardar(idActual, nuevo)
	ag.proximo_id_post++

	ag.usuarios_hash.Iterar(func(_ string, u Usuario) bool {
		if u != ag.usuario_loggeado {
			u.RecibirPost(nuevo)
		}
		return true
	})

	return MSJ_POST_PUBLICADO
}

func (ag *algogramImpl) VerSiguienteFeed() string {
	if ag.usuario_loggeado == nil {
		return ERR_LOGUEADO_POST
	}
	post := ag.usuario_loggeado.ProximoPost()
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
		// No las uno en un or porque Obtener() me da panic si la clave no existe
		return ERR_POST_INEXISTENTE
	}

	post := ag.posts.Obtener(id)
	post.DarLike(ag.usuario_loggeado.Nombre())
	return MSJ_POST_LIKEADO
}

func (ag *algogramImpl) MostrarLikes(id int) string {
	if !ag.posts.Pertenece(id) {
		return ERR_SIN_LIKES
	}
	post := ag.posts.Obtener(id)
	if post.CantidadLikes() == 0 {
		return ERR_SIN_LIKES
	}
	return post.MostrarLikes()
}
