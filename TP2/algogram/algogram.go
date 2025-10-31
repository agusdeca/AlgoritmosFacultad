package algogram

import (
    "tp2/tdas/diccionario"
)

type AlgoGram interface {
	// CargarUsuarios lee el archivo
	CargarUsuarios(ruta string) error

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

const CAPACIDAD_INICIAL_USUARIOS = 2

// Crea la instancia de la implementación
func NewAlgoGram() AlgoGram {
	return &algogramImpl{
		usuarios_hash:   diccionario.CrearHash[string, *usuario](func(a, b string) bool { return a == b }),
		usuarios_lista:  make([]*usuario, CAPACIDAD_INICIAL_USUARIOS),
		usuarios_cantidad: 0,
		posts:           diccionario.CrearHash[int, *post](func(a, b int) bool { return a == b }),
		usuario_loggeado: nil,
		proximo_id_post: 0,
	}
}