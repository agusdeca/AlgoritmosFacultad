package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tp2/codigo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error: Debe dar la ruta al archivo de usuarios")
		os.Exit(1)
	}
	rutaArchivo := os.Args[1]

	// leemos el archivo para obtener la lista de nombres
	nombres, err := leerArchivoUsuarios(rutaArchivo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al cargar usuarios: %s\n", err.Error())
		os.Exit(1)
	}

	// Crear el sistema con la lista de usuarios
	sistema := codigo.NewAlgoGramFromUsers(nombres)

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
				fmt.Println(codigo.ERR_POST_INEXISTENTE)
				continue
			}
			fmt.Println(sistema.LikearPost(id))

		case "mostrar_likes":
			id, err := strconv.Atoi(parametro)
			if err != nil {
				fmt.Println(codigo.ERR_SIN_LIKES)
				continue
			}
			fmt.Println(sistema.MostrarLikes(id))
		}
	}
}

// leerArchivoUsuarios lee el archivo de usuarios y devuelve la lista de nombres
func leerArchivoUsuarios(ruta string) ([]string, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	var nombres []string
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		nombre := scanner.Text()
		nombres = append(nombres, nombre)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nombres, nil
}

// Si no hay parámetro, devuelve parametro=""
func parsearLinea(linea string) (string, string) {
	partes := strings.SplitN(linea, " ", 2)
	comando := partes[0]
	if len(partes) == 1 {
		return comando, ""
	}
	return comando, partes[1]
}
