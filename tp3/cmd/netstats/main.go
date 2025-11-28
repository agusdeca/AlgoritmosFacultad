package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"tp3"
	"tp3/tdas/grafo"
)

const (
	MAX_CAPACITY = 10 * 1024 * 1024
	ERR_FALTAN_PARAMETROS      = "Error: faltan parámetros"
	ERR_FORMATO_INCORRECTO     = "Error: formato incorrecto"
	ERR_FALTA_N                = "Error: falta el parámetro n"
	ERR_PARAMETRO_NO_NUMERO    = "Error: parámetro debe ser un número"
	ERR_FALTA_PAGINA           = "Error: falta la página"
	ERR_N_NO_NUMERO            = "Error: n debe ser un número"
	ERR_FALTAN_PAGINAS         = "Error: faltan páginas"
	ERR_NO_RECORRIDO           = "No se encontro recorrido"
	ERR_NO_FORMA_ORDEN         = "No existe forma de leer las paginas en orden"
	ERR_NO_HAY_CAMINO          = "No hay camino"
)

func cargarGrafo(archivo string) (grafo.GrafoNoPesado[string], error) {
	file, err := os.Open(archivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	debug.SetGCPercent(-1)
	defer debug.SetGCPercent(100)

	g := grafo.CrearGrafoNoPesado[string](true)
	scanner := bufio.NewScanner(file)

	buf := make([]byte, MAX_CAPACITY)
	scanner.Buffer(buf, MAX_CAPACITY)

	fmt.Fprintf(os.Stderr, "Cargando grafo desde %s...\n", archivo)
	lineas := 0

	for scanner.Scan() {
		linea := scanner.Text()
		if len(linea) == 0 {
			continue
		}

		primerTab := strings.IndexByte(linea, '\t')

		var origen string
		if primerTab == -1 {
			origen = linea
			g.AgregarVertice(origen)
		} else {
			origen = linea[:primerTab]
			g.AgregarVertice(origen)

			resto := linea[primerTab+1:]
			start := 0
			for i := 0; i < len(resto); i++ {
				if resto[i] == '\t' {
					if i > start {
						destino := resto[start:i]
						g.AgregarArista(origen, destino)
					}
					start = i + 1
				}
			}
			if start < len(resto) {
				destino := resto[start:]
				if destino != "" {
					g.AgregarArista(origen, destino)
				}
			}
		}

		lineas++
		if lineas%20000 == 0 {
			fmt.Fprintf(os.Stderr, "Procesadas %d líneas...\r", lineas)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	debug.FreeOSMemory()
	
	fmt.Fprintf(os.Stderr, "\nGrafo cargado: %d vértices procesados.\n", g.Cantidad())
	return g, nil
}

func procesarComando(g grafo.GrafoNoPesado[string], comando string) {
	partes := strings.SplitN(comando, " ", 2)
	cmd := partes[0]

	switch cmd {
	case "listar_operaciones":
		fmt.Println("camino")
		fmt.Println("mas_importantes")
		fmt.Println("conectados")
		fmt.Println("ciclo")
		fmt.Println("lectura")
		fmt.Println("diametro")
		fmt.Println("rango")
		fmt.Println("navegacion")
		fmt.Println("clustering")
		fmt.Println("comunidad")

	case "camino":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTAN_PARAMETROS)
			return
		}
		params := strings.Split(partes[1], ",")
		if len(params) != 2 {
			fmt.Println(ERR_FORMATO_INCORRECTO)
			return
		}
		origen := strings.TrimSpace(params[0])
		destino := strings.TrimSpace(params[1])

		camino, costo := tp3.CaminoMinimo(g, origen, destino)
		if camino == nil {
			fmt.Println(ERR_NO_RECORRIDO)
		} else {
			fmt.Println(strings.Join(camino, " -> "))
			fmt.Printf("Costo: %d\n", costo)
		}

	case "mas_importantes":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTA_N)
			return
		}
		n, err := strconv.Atoi(strings.TrimSpace(partes[1]))
		if err != nil {
			fmt.Println(ERR_PARAMETRO_NO_NUMERO)
			return
		}

		importantes := tp3.MasImportantes(g, n)
		fmt.Println(strings.Join(importantes, ", "))

	case "conectados":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTA_PAGINA)
			return
		}
		pagina := strings.TrimSpace(partes[1])

		conectados := tp3.Conectados(g, pagina)
		fmt.Println(strings.Join(conectados, ", "))

	case "ciclo":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTAN_PARAMETROS)
			return
		}
		params := strings.Split(partes[1], ",")
		if len(params) != 2 {
			fmt.Println(ERR_FORMATO_INCORRECTO)
			return
		}
		pagina := strings.TrimSpace(params[0])
		n, err := strconv.Atoi(strings.TrimSpace(params[1]))
		if err != nil {
			fmt.Println(ERR_N_NO_NUMERO)
			return
		}

		ciclo := tp3.Ciclo(g, pagina, n)
		if ciclo == nil {
			fmt.Println(ERR_NO_RECORRIDO)
		} else {
			fmt.Println(strings.Join(ciclo, " -> "))
		}

	case "lectura":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTAN_PAGINAS)
			return
		}
		paginasStr := strings.Split(partes[1], ",")
		paginas := make([]string, len(paginasStr))
		for i, p := range paginasStr {
			paginas[i] = strings.TrimSpace(p)
		}

		orden := tp3.Lectura(g, paginas)
		if orden == nil {
			fmt.Println(ERR_NO_FORMA_ORDEN)
		} else {
			fmt.Println(strings.Join(orden, ", "))
		}

	case "diametro":
		camino, costo := tp3.Diametro(g)
		if camino == nil {
			fmt.Println(ERR_NO_HAY_CAMINO)
		} else {
			fmt.Println(strings.Join(camino, " -> "))
			fmt.Printf("Costo: %d\n", costo)
		}

	case "rango":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTAN_PARAMETROS)
			return
		}
		params := strings.Split(partes[1], ",")
		if len(params) != 2 {
			fmt.Println(ERR_FORMATO_INCORRECTO)
			return
		}
		pagina := strings.TrimSpace(params[0])
		n, err := strconv.Atoi(strings.TrimSpace(params[1]))
		if err != nil {
			fmt.Println(ERR_N_NO_NUMERO)
			return
		}

		cantidad := tp3.EnRango(g, pagina, n)
		fmt.Println(cantidad)

	case "navegacion":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTA_PAGINA)
			return
		}
		pagina := strings.TrimSpace(partes[1])

		camino := tp3.Navegacion(g, pagina)
		fmt.Println(strings.Join(camino, " -> "))

	case "comunidad":
		if len(partes) < 2 {
			fmt.Println(ERR_FALTA_PAGINA)
			return
		}
		pagina := strings.TrimSpace(partes[1])

		comunidad := tp3.Comunidad(g, pagina)
		fmt.Println(strings.Join(comunidad, ", "))

	case "clustering":
		if len(partes) < 2 {
			// Clustering promedio
			vertices := g.ObtenerVertices()
			suma := 0.0
			for _, v := range vertices {
				suma += tp3.Clustering(g, v)
			}
			promedio := suma / float64(len(vertices))
			fmt.Printf("%.3f\n", promedio)
		} else {
			// Clustering de una página
			pagina := strings.TrimSpace(partes[1])
			coef := tp3.Clustering(g, pagina)
			fmt.Printf("%.3f\n", coef)
		}

	default:
		fmt.Printf("Comando desconocido: %s\n", cmd)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: ./netstats <archivo>")
		os.Exit(1)
	}

	archivo := os.Args[1]

	g, err := cargarGrafo(archivo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al cargar el archivo: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comando := strings.TrimSpace(scanner.Text())
		if comando == "" || strings.HasPrefix(comando, "#") {
			continue
		}
		procesarComando(g, comando)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al leer entrada: %v\n", err)
		os.Exit(1)
	}
}
