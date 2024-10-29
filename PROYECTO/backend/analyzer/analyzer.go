package analyzer

import (
	"backend/commands"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Analyzer analiza el comando de entrada y llama a la función de análisis correspondiente.
// Soporta los comandos mkdisk, rep y execute.
func Analyzer(input string) (string, error) {
	// Dividir la entrada en líneas y procesar cada línea individualmente
	lines := strings.Split(input, "\n")
	var result strings.Builder

	for _, line := range lines {
		// Ignorar líneas que comienzan con "#"
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		var output string
		var err error

		switch tokens[0] {
		case "mkdisk":
			output, err = commands.ParserMkdisk(tokens[1:])
		case "rmdisk":
			output, err = commands.ParserRmdisk(tokens[1:])
		case "fdisk":
			output, err = commands.ParserFdisk(tokens[1:])
		case "mkfs":
			output, err = commands.ParserMkfs(tokens[1:])
		case "mkdir":
			output, err = commands.ParserMkdir(tokens[1:])
		case "mkfile":
			output, err = commands.ParserMkfile(tokens[1:])
		case "rep":
			output, err = commands.ParserRep(tokens[1:])
		case "execute":
			output, err = commands.ParserExecute(tokens[1:])
		case "mount":
			output, err = commands.ParserMount(tokens[1:])
		case "unmount":
			output, err = commands.ParserUnmount(tokens[1:])
		case "cat":
			output, err = commands.ParserCat(tokens[1:])
		case "login":
			output, err = commands.ParserLogin(tokens[1:])
		case "logout":
			output, err = commands.ParserLogout()
		case "clear":
			// Crea un comando para limpiar la terminal
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			if err != nil {
				output, err = "", errors.New("error al limpiar la terminal")
			}
		default:
			output, err = "", fmt.Errorf("comando desconocido: %s", tokens[0])
		}

		if err != nil {
			return "", err
		}

		result.WriteString(output)
		result.WriteString("\n")
	}

	return result.String(), nil
}
