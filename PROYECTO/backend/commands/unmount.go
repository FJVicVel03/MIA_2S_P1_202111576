package commands

import (
	"backend/global"
	"errors"
	"fmt"
	"regexp"
)

// ParserUnmount analiza el comando unmount y llama a la función correspondiente.
func ParserUnmount(tokens []string) (string, error) {
	if len(tokens) == 0 {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	// Expresión regular para encontrar el parámetro -id
	re := regexp.MustCompile(`-id=\S+`)
	match := re.FindString(tokens[0])
	if match == "" {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	// Extraer el valor del parámetro -id
	id := match[4:]

	return commandUnmount(id)
}

// commandUnmount desmonta la partición con el ID proporcionado.
func commandUnmount(id string) (string, error) {
	// Verificar si la partición está montada
	_, exists := global.MountedPartitions[id]
	if !exists {
		return "", errors.New("no se encontró la partición con id: " + id)
	}

	// Desmontar la partición
	delete(global.MountedPartitions, id)
	return fmt.Sprintf("UNMOUNT: Partición con id %s desmontada exitosamente", id), nil
}
