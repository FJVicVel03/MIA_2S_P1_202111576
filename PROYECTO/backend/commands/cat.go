package commands

import (
	"backend/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// CAT estructura que representa el comando cat con sus parámetros
type CAT struct {
	files []string // Lista de archivos a concatenar
}

// ParserCat parsea el comando cat y devuelve una instancia de CAT
func ParserCat(tokens []string) (string, error) {
	// Verificar si hay una sesión iniciada
	if !utils.IsLoggedIn() {
		return "", errors.New("no hay ninguna sesión activa")
	}

	cmd := &CAT{} // Crea una nueva instancia de CAT

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando cat
	re := regexp.MustCompile(`-file\d+="[^"]+"|-file\d+=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		var value string
		if len(kv) == 2 {
			value = kv[1]
		}

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Agregar el archivo a la lista de archivos
		cmd.files = append(cmd.files, value)
	}

	// Verificar que al menos un archivo haya sido proporcionado
	if len(cmd.files) == 0 {
		return "", errors.New("faltan parámetros requeridos: -fileN")
	}

	// Concatenar el contenido de los archivos
	content, err := commandCat(cmd)
	if err != nil {
		return "", err
	}

	return content, nil // Devuelve el contenido concatenado de los archivos
}

// Función para concatenar el contenido de los archivos
func commandCat(cat *CAT) (string, error) {
	var contentBuilder strings.Builder

	for _, file := range cat.files {
		// Verificar permisos de lectura
		if !hasReadPermission(file) {
			return "", fmt.Errorf("no tiene permiso de lectura para el archivo: %s", file)
		}

		// Leer el contenido del archivo
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("error al leer el archivo: %w", err)
		}

		// Agregar el contenido al builder
		contentBuilder.WriteString(string(content))
		contentBuilder.WriteString("\n")
	}

	return contentBuilder.String(), nil
}

// hasReadPermission verifica si el usuario tiene permiso de lectura para el archivo
func hasReadPermission(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}

	// Verificar permisos de lectura (esto es un ejemplo simplificado)
	return info.Mode().Perm()&(1<<(uint(7))) != 0
}
