package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func Int32ToBytes(n int32) [4]byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(n))
	return buf
}

func Float64ToBytes(f float64) [4]byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(f))
	return buf
}

func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
	case "K":
		return size * 1024, nil // Convierte kilobytes a bytes
	case "M":
		return size * 1024 * 1024, nil // Convierte megabytes a bytes
	default:
		return 0, errors.New("invalid unit") // Devuelve un error si la unidad es inválida
	}
}

// Carnet de estudiante
const Carnet string = "76" // 202111576
// Lista con todo el abecedario
var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// Mapa para almacenar la asignación de letras a los diferentes paths
var pathToLetter = make(map[string]string)

// Índice para la siguiente letra disponible en el abecedario
var nextLetterIndex = 0

// GetLetter obtiene la letra asignada a un path
func GetLetter(path string) (string, error) {
	// Asignar una letra al path si no tiene una asignada
	if _, exists := pathToLetter[path]; !exists {
		if nextLetterIndex < len(alphabet) {
			pathToLetter[path] = alphabet[nextLetterIndex]
			nextLetterIndex++
		} else {
			fmt.Println("Error: no hay más letras disponibles para asignar")
			return "", errors.New("no hay más letras disponibles para asignar")
		}
	}

	return pathToLetter[path], nil
}