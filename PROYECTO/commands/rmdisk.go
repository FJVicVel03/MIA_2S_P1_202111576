package commands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParserRmdisk analiza el comando rmdisk y extrae el parámetro path.

type RMDISK struct {
	path string // Ruta del archivo del disco
}

/*
	rmdisk -path=/home/user/Disco1.mia
	rmdisk -path="/home/mis discos/Disco4.mia"
*/

func ParserRmdisk(tokens []string) (*RMDISK, error) {
	cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar el parámetro del comando rmdisk
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)
	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar el parámetro -path
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	// Lógica para eliminar el archivo del disco
	err := os.Remove(cmd.path)
	if err != nil {
		return nil, fmt.Errorf("error al eliminar el disco: %v", err)
	}

	fmt.Println("Disco eliminado exitosamente")
	return cmd, nil // Devuelve el comando RMDISK creado
}
