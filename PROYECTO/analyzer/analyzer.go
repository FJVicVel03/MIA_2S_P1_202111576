package analyzer

import (
	"PROYECTO/commands"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Analyzer analiza el comando de entrada y llama a la función de análisis correspondiente.
// Soporta los comandos mkdisk, rep y execute.
func Analyzer(input string) (interface{}, error) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}

	switch tokens[0] {
	case "mkdisk":
		return commands.ParserMkdisk(tokens[1:])
	case "rmdisk":
		return commands.ParserRmdisk(tokens[1:])
	case "fdisk":
		return commands.ParserFdisk(tokens[1:])
	case "mkfs":
		return commands.ParserMkfs(tokens[1:])
	case "rep":
		return commands.ParserRep(tokens[1:])
	case "execute":
		return commands.ParserExecute(tokens[1:])
	case "mount":
		return commands.ParserMount(tokens[1:])
	case "clear":
		//crea un comando para limpiar la terminal
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return nil, errors.New("error al limpiar la terminal")
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
