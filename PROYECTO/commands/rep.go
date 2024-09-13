package commands

import (
	"PROYECTO/structures"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParserRep analiza el comando rep y extrae el parámetro path.
// Luego llama a la función commandRep para deserializar el MBR y mostrar su información.
func ParserRep(tokens []string) (interface{}, error) {
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

	return nil, commandRep(path)
}

// commandRep deserializa el MBR desde el archivo en la ruta especificada y muestra su información.
// Llama a la función DeserializeMBR para obtener la estructura MBR y luego a la función Print para mostrarla.
func commandRep(path string) error {
	mbr, err := structures.DeserializeMBR(path)
	if err != nil {
		return err
	}

	mbr.Print()

	return nil
}
