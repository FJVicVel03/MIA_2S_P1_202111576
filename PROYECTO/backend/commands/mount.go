package commands

import (
	"backend/global"
	"backend/structures"
	"backend/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
	path string // Ruta del archivo del disco
	name string // Nombre de la partición
}

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func ParserMount(tokens []string) (string, error) {
	cmd := &MOUNT{} // Crea una nueva instancia de MOUNT

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
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
		case "-path":
			cmd.path = value
		case "-name":
			cmd.name = value
		default:
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := commandMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return "MOUNT: Partición montada exitosamente", nil // Devuelve el comando MOUNT creado
}

// mount -path=Disco2.mia -name=Particion3
func commandMount(mount *MOUNT) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(mount.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Buscar la partición con el nombre especificado
	partition, indexPartition := mbr.GetPartitionByName(mount.name)
	if partition == nil {
		// Si no se encuentra una partición primaria o extendida, buscar en las particiones lógicas
		for _, part := range mbr.Mbr_partitions {
			if part.Part_type[0] == 'E' {
				var ebr structures.EBR
				err := ebr.DeserializeEBR(mount.path, part.Part_start)
				if err != nil {
					fmt.Println("Error deserializando el EBR:", err)
					return err
				}
				for {
					ebrName := strings.Trim(string(ebr.Part_name[:]), "\x00 ")
					if strings.EqualFold(ebrName, mount.name) {
						var partID [4]byte
						copy(partID[:], ebr.Part_id[:])
						partition = &structures.PARTITION{
							Part_status:      [1]byte{ebr.Part_status},
							Part_type:        [1]byte{"L"[0]},
							Part_fit:         [1]byte{ebr.Part_fit},
							Part_start:       ebr.Part_start,
							Part_size:        ebr.Part_size,
							Part_name:        ebr.Part_name,
							Part_correlative: 0,
							Part_id:          partID,
						}
						indexPartition = -1
						break
					}
					if ebr.Part_next == -1 {
						break
					}
					err = ebr.DeserializeEBR(mount.path, ebr.Part_next)
					if err != nil {
						fmt.Println("Error deserializando el EBR:", err)
						return err
					}
				}
			}
		}
		if partition == nil {
			fmt.Println("Error: la partición no existe")
			return errors.New("la partición no existe")
		}
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición fue encontrada
	fmt.Println("\nPartición encontrada:")
	partition.Print()

	// Generar un id único para la partición
	logicalIndex := -1
	if partition.Part_type[0] == 'L' {
		logicalIndex = 0 // or any other logic to determine the logical partition index
	}
	idPartition, err := GenerateIdPartition(mount, logicalIndex)
	if err != nil {
		fmt.Println("Error generando el id de partición:", err)
		return err
	}

	// Guardar la partición montada en la lista de montajes globales
	global.MountedPartitions[idPartition] = mount.path

	// Modificamos la partición para indicar que está montada
	if indexPartition != -1 {
		partition.MountPartition(indexPartition, idPartition)
		mbr.Mbr_partitions[indexPartition] = *partition
	} else {
		// Handle logical partition mounting
		var ebr structures.EBR
		err := ebr.DeserializeEBR(mount.path, partition.Part_start)
		if err != nil {
			fmt.Println("Error deserializando el EBR:", err)
			return err
		}
		ebr.Part_status = '2' // Indicate that the partition is mounted
		copy(ebr.Part_id[:], idPartition)
		err = ebr.SerializeEBR(mount.path, partition.Part_start)
		if err != nil {
			fmt.Println("Error serializando el EBR:", err)
			return err
		}
	}

	// Serializar la estructura MBR en el archivo binario
	err = mbr.Serialize(mount.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición fue montada
	fmt.Println("\nPartición montada exitosamente:")
	partition.Print()

	return nil
}

func GenerateIdPartition(mount *MOUNT, indexPartition int) (string, error) {
	// Asignar una letra a la partición
	letter, err := utils.GetLetter(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", err
	}

	// Incrementar el contador global
	global.PartitionCounter++

	// Crear id de partición
	idPartition := fmt.Sprintf("%s%d%s", utils.Carnet, global.PartitionCounter, letter)

	return idPartition, nil
}
