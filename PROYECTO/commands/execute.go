package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParserExecute analiza el comando execute y extrae el parámetro path.
// Luego llama a la función commandExecute para crear un disco y mostrar su información.
func ParserExecute(tokens []string) (interface{}, error) {
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	matches := re.FindAllString(args, -1)

	var path string
	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}
		switch key {
		case "-path":
			path = value
		default:
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	return nil, commandExecute(path)
}

// commandExecute crea un disco de 5MB y muestra su información.
// Llama directamente a las funciones ParserMkdisk y ParserRep para evitar dependencias cíclicas.
func commandExecute(path string) error {
	// Crear el disco
	mkdiskTokens := []string{"-path=" + path, "-size=5"}
	_, err := ParserMkdisk(mkdiskTokens)
	if err != nil {
		return fmt.Errorf("error al crear el disco: %v", err)
	}

	// Mostrar la información del disco
	repTokens := []string{"-path=" + path}
	_, err = ParserRep(repTokens)
	if err != nil {
		return fmt.Errorf("error al mostrar la información del disco: %v", err)
	}

	return nil
}
