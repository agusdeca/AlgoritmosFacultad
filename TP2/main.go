package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tp2/algogram"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error: Debe dar la ruta al archivo de usuarios")
		os.Exit(1)
	}
	rutaArchivo := os.Args[1]

	sistema := algogram.NewAlgoGram()

	//Cargar los usuarios iniciales
	if err := sistema.CargarUsuarios(rutaArchivo); err != nil {
		fmt.Fprintf(os.Stderr, "Error al cargar usuarios: %s\n", err.Error())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comando, parametro := parsearLinea(scanner.Text())

		switch comando {
		case "login":
			fmt.Println(sistema.Login(parametro))

		case "logout":
			fmt.Println(sistema.Logout())

		case "publicar":
			fmt.Println(sistema.Publicar(parametro))

		case "ver_siguiente_feed":
			fmt.Println(sistema.VerSiguienteFeed())

		case "likear_post":
			id, err := strconv.Atoi(parametro)
			if err != nil {
				fmt.Println("Error: Usuario no loggeado o Post inexistente")
				continue
			}
			fmt.Println(sistema.LikearPost(id))

		case "mostrar_likes":
			id, err := strconv.Atoi(parametro)
			if err != nil {
				fmt.Println("Error: Post inexistente o sin likes")
				continue
			}
			fmt.Println(sistema.MostrarLikes(id))
		
		}
	}
}

// Si no hay parámetro, devuelve parametro="".
func parsearLinea(linea string) (string, string) {
	partes := strings.SplitN(linea, " ", 2)
	
	comando := partes[0]
	
	if len(partes) == 1 {
		return comando, "" // No hay parámetro
	}
	
	return comando, partes[1] // Hay parámetro
}