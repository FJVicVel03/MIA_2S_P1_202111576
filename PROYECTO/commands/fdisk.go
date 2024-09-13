package commands

import (
	"PROYECTO/structures"
	"PROYECTO/utils"
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

func ParserFdisk(tokens []string) (*FDISK, error) {
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
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
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
				return nil, errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" {
				return nil, errors.New("la unidad debe ser K o M")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return nil, errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
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

	return cmd, nil // Devuelve el comando FDISK creado
}
func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

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
		err = createLogicalPartition(fdisk, sizeBytes)
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

	// Deserialize the MBR from the binary file
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR", err)
		return err
	}

	// Find the extended partition
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

	// Deserialize the first EBR from the extended partition start
	var ebr structures.EBR
	err = ebr.DeserializeEBR(fdisk.path, extendedPartition.Part_start)
	if err != nil {
		fmt.Println("Error deserializando el EBR:", err)
		return err
	}

	// Find the first available EBR slot
	for ebr.Part_next != -1 {
		fmt.Println("EBR encontrado, siguiente EBR en:", ebr.Part_next)
		if ebr.Part_next <= 0 {
			break
		}
		err = ebr.DeserializeEBR(fdisk.path, ebr.Part_next)
		if err != nil {
			fmt.Println("Error deserializando el EBR:", err)
			return err
		}
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición esté disponible
	fmt.Println("\nPartición disponible:")
	ebr.Print()

	// Create the new logical partition
	ebr.Part_status = '1'
	ebr.Part_start = ebr.Part_start + ebr.Part_size
	ebr.Part_size = int64(sizeBytes)
	ebr.Part_fit = fdisk.fit[0]
	copy(ebr.Part_name[:], fdisk.name)

	// Print para verificar que la partición se haya creado correctamente
	fmt.Println("\nPartición creada (modificada):")
	ebr.Print()

	// Serialize the new EBR
	err = ebr.SerializeEBR(fdisk.path, ebr.Part_start)
	if err != nil {
		fmt.Println("Error serializando el EBR:", err)
		return err
	}

	fmt.Println("Partición lógica creada exitosamente")
	return nil
}
