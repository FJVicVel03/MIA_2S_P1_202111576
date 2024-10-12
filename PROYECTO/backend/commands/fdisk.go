package commands

import (
	"backend/structures"
	"backend/utils"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
Este comando se encarga de la administración de particiones en el archivo
que representa al disco duro virtual. Permite a los usuarios crear, eliminar o
modificar particiones dentro del archivo de disco duro.
En caso de que no se pueda realizar la operación solicitada sobre una
partición, la aplicación deberá mostrar un mensaje de error detallado. El
mensaje de error especificará claramente la razón por la cual la operación
falló, como podría ser la falta de espacio disponible, restricciones en el
número máximo de particiones permitidas, errores en los parámetros
ingresados, o cualquier otra causa relevante. Esta retroalimentación detallada
ayudará a los usuarios a entender y corregir los problemas para completar la
operación de manera exitosa.
*/

type FDISK struct {
	size int    // Tamaño de la partición
	unit string // Unidad de medida del tamaño
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del disco
	typ  string // Tipo de partición (P, E, L)
	name string // Nombre de la partición
}

// fdisk -size=300 -path=/home/Disco1.mia -name=Particion1

func ParserFdisk(tokens []string) (string, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return "", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" {
				return "", errors.New("la unidad debe ser K o M")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return "", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}

	// Crear la partición con los parámetros proporcionados
	err := commandFdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return "FDISK: Partición creada exitosamente", nil // Devuelve el comando FDISK creado
}
func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	fmt.Printf("Tamaño convertido a bytes: %d\n", sizeBytes) // Agregar esta línea para depuración

	if fdisk.typ == "P" {
		// Crear partición primaria
		err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return err
		}
	} else if fdisk.typ == "E" {
		// Crear partición extendida
		err = createExtendedPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición extendida:", err)
		}
	} else if fdisk.typ == "L" {
		// Crear partición lógica
		err = createLogicalPartition(fdisk, sizeBytes/1024)
		if err != nil {
			fmt.Println("Error creando partición lógica:", err)
		}
	}

	return nil
}

func createPrimaryPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición esté disponible
	fmt.Println("\nPartición disponible:")
	availablePartition.Print()

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Print para verificar que la partición se haya creado correctamente
	fmt.Println("\nPartición creada (modificada):")
	availablePartition.Print()

	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Imprimir las particiones del MBR
	fmt.Println("\nParticiones del MBR:")
	mbr.PrintPartitions()

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

// fdisk -type=E -path=Disco2.mia -unit=K -name=Particion2 -size=300
func createExtendedPartition(fdisk *FDISK, sizeBytes int) error {
	var mbr structures.MBR

	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'E' {
			return errors.New("ya existe una partición extendida")
		}
	}

	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		return errors.New("no hay particiones disponibles")
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición esté disponible
	fmt.Println("\nPartición disponible:")
	availablePartition.Print()

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, "E", fdisk.fit, fdisk.name)

	// Print para verificar que la partición se haya creado correctamente
	fmt.Println("\nPartición creada (modificada):")
	availablePartition.Print()

	mbr.Mbr_partitions[indexPartition] = *availablePartition

	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	fmt.Println("Partición extendida creada exitosamente")
	return nil
}

// fdisk -size=1 -type=L -unit=M -fit=BF -path=Disco2.mia -name=Particion3

func createLogicalPartition(fdisk *FDISK, sizeBytes int) error {
	var mbr structures.MBR

	// Deserializar el MBR desde el archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Encontrar la partición extendida
	var extendedPartition *structures.PARTITION
	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'E' {
			extendedPartition = &mbr.Mbr_partitions[i]
			break
		}
	}

	if extendedPartition == nil {
		return errors.New("no se encontró partición extendida")
	}

	// Deserializar el primer EBR en la partición extendida
	var ebr structures.EBR
	err = ebr.DeserializeEBR(fdisk.path, int64(extendedPartition.Part_start))
	if err != nil {
		fmt.Println("Error deserializando el EBR:", err)
		return err
	}

	// Si es la primera partición lógica, el EBR debe estar al inicio de la partición extendida
	if ebr.Part_size == 0 {
		ebr.Part_status = '1'
		ebr.Part_fit = fdisk.fit[0]
		ebr.Part_start = int64(extendedPartition.Part_start) + int64(binary.Size(ebr)) // El primer EBR va justo al inicio de la partición extendida
		ebr.Part_size = int64(sizeBytes)
		ebr.Part_next = -1 // Sin siguiente partición lógica por ahora
		copy(ebr.Part_name[:], fdisk.name)

		// Serializar el primer EBR
		err = ebr.SerializeEBR(fdisk.path, int64(extendedPartition.Part_start))
		if err != nil {
			return err
		}
		fmt.Println(" ")
		ebr.Print()
		return nil
	}

	// Buscar el último EBR en la cadena de particiones lógicas
	var prevEBR *structures.EBR
	for ebr.Part_next != -1 {
		prevEBR = &ebr
		err = ebr.DeserializeEBR(fdisk.path, ebr.Part_next)
		if err != nil {
			return err
		}
	}

	// Calcular el inicio de la nueva partición lógica
	newStart := ebr.Part_start + ebr.Part_size + int64(binary.Size(ebr))

	// Verificar si hay suficiente espacio en la partición extendida
	if newStart+int64(sizeBytes) > int64(extendedPartition.Part_start)+int64(extendedPartition.Part_size) {
		return errors.New("no hay suficiente espacio en la partición extendida")
	}

	// Crear el nuevo EBR para la partición lógica
	newEBR := structures.EBR{
		Part_status: '1',
		Part_fit:    fdisk.fit[0],
		Part_start:  newStart,
		Part_size:   int64(sizeBytes),
		Part_next:   -1, // El nuevo EBR no apunta a ningún otro por ahora
	}
	copy(newEBR.Part_name[:], fdisk.name)

	// Si existe un EBR previo, actualizamos su Part_next para apuntar a la nueva partición lógica
	if prevEBR != nil {
		prevEBR.Part_next = newEBR.Part_start
		err = prevEBR.SerializeEBR(fdisk.path, prevEBR.Part_start)
		if err != nil {
			return err
		}
	}

	// Serializar el nuevo EBR
	err = newEBR.SerializeEBR(fdisk.path, newEBR.Part_start)
	if err != nil {
		return err
	}
	fmt.Println(" ")
	newEBR.Print()
	return nil
}
