// login.go
package commands

import (
	"backend/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParserLogin analiza el comando login y extrae los parámetros user, pass e id.
// Luego llama a la función commandLogin para iniciar sesión.
func ParserLogin(tokens []string) (string, error) {
	if utils.IsLoggedIn() {
		return "", errors.New("ya hay un usuario logueado, por favor cierre sesión antes de iniciar una nueva")
	}

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-user="[^"]+"|-user=[^\s]+|-pass="[^"]+"|-pass=[^\s]+|-id="[^"]+"|-id=[^\s]+`)
	matches := re.FindAllString(args, -1)

	var user, pass, id string
	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}
		switch key {
		case "-user":
			user = value
		case "-pass":
			pass = value
		case "-id":
			id = value
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if user == "" || pass == "" || id == "" {
		return "", errors.New("faltan parámetros requeridos: -user, -pass, -id")
	}

	return "LOGIN: realizado correctamente", commandLogin(user, pass, id)
}

// commandLogin verifica las credenciales del usuario y establece la sesión.
func commandLogin(user, pass, id string) error {
	if !verifyUser(user, pass) {
		return errors.New("autenticación fallida: usuario o contraseña incorrectos")
	}

	// Establecer la sesión
	utils.GlobalSession = &utils.Session{
		User: user,
		ID:   id,
	}

	return nil
}

// verifyUser es una función ficticia que verifica las credenciales del usuario.
func verifyUser(user, pass string) bool {
	return user == "root" && pass == "123" // Ejemplo estático
}
