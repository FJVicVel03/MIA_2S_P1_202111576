// logout.go
package commands

import (
	"backend/utils"
	"errors"
)

// ParserLogout analiza el comando logout y cierra la sesión si hay una activa.
func ParserLogout() (string, error) {
	if !utils.IsLoggedIn() {
		return "", errors.New("no hay ninguna sesión activa")
	}

	// Cerrar la sesión
	utils.GlobalSession = nil
	return "LOGOUT: sesión cerrada correctamente", nil
}
