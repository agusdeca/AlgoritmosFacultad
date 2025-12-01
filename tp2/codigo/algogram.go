package codigo

import (
	"tdas/diccionario"
)

const (
	MSJ_HOLA              = "Hola %s"
	MSJ_ADIOS             = "Adios"
	ERR_NO_LOGGEADO       = "Error: no habia usuario loggeado"
	ERR_USUARIO_NO_EXISTE = "Error: usuario no existente"
	ERR_USUARIO_EXISTE    = "Error: Ya habia un usuario loggeado"
	ERR_POST_INEXISTENTE  = "Error: Usuario no loggeado o Post inexistente"
	ERR_LOGUEADO_POST     = "Usuario no loggeado o no hay mas posts para ver"
	ERR_SIN_LIKES         = "Error: Post inexistente o sin likes"
	MSJ_POST_PUBLICADO    = "Post publicado"
	MSJ_POST_LIKEADO      = "Post likeado"
)

type AlgoGram interface {
	// Login valida e inicia sesión de un usuario.
	Login(nombre string) string

	// Logout cierra la sesión actual.
	Logout() string

	// Publicar crea un nuevo post del usuario logeado y lo distribuye.
	Publicar(texto string) string

	// VerSiguienteFeed obtiene el post prioritario del feed del usuario.
	VerSiguienteFeed() string

	// LikearPost agrega un like a un post.
	LikearPost(id int) string

	// MostrarLikes devuelve la lista de usuarios que likearon un post.
	MostrarLikes(id int) string
}

// constructor sin usuarios (vacío)
func NewAlgoGram() AlgoGram {
	return &algogramImpl{
		usuarios_hash:    diccionario.CrearHash[string, Usuario](func(a, b string) bool { return a == b }),
		posts:            diccionario.CrearHash[int, Post](func(a, b int) bool { return a == b }),
		usuario_loggeado: nil,
		proximo_id_post:  0,
	}
}

// constructor que recibe una lista de nombres de usuario (la lectura del archivo queda en main)
func NewAlgoGramFromUsers(nombres []string) AlgoGram {
	ag := NewAlgoGram().(*algogramImpl)
	for i, nombre := range nombres {
		u := nuevoUsuario(nombre, i)
		ag.usuarios_hash.Guardar(nombre, u)
	}
	return ag
}
