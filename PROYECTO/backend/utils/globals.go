package utils

// Declaracion de variables globales
var (
	GlobalMounts  map[string]string = make(map[string]string)
	GlobalSession *Session
)

// Estructura para almacenar el estado de la sesión
type Session struct {
	User string
	ID   string
}
