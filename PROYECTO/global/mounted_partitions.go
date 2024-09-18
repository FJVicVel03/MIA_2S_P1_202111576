package global

import (
	"PROYECTO/structures"
	"errors"
	"path/filepath"
	"strings"
)

const Carnet string = "76" // 202111576
// Declaración de las particiones montadas

var (
	MountedPartitions map[string]string = make(map[string]string)
)
var PartitionCounter int = 0

// GetMountedPartition obtiene la partición montada con el id especificado
func GetMountedPartition(id string) (*structures.PARTITION, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, path, nil
}

// GetDiskNameByID obtiene el nombre del disco a partir del ID de la partición montada
func GetDiskNameByID(id string) (string, error) {
	// Buscar el path del disco usando el ID de la partición montada
	path, exists := MountedPartitions[id]
	if !exists {
		return "", errors.New("no se encontró una partición montada con el ID especificado")
	}

	// Extraer el nombre del disco desde el path
	// filepath.Base devuelve el último componente del path (el nombre del archivo)
	diskName := filepath.Base(path)

	// Si el disco tiene una extensión, la eliminamos
	// Esto es útil si tienes discos con nombres como "disco1.dsk" y quieres solo "disco1"
	diskName = strings.TrimSuffix(diskName, filepath.Ext(diskName))

	return diskName, nil
}

// GetMountedMBR obtiene el MBR de la partición montada con el id especificado
func GetMountedPartitionRep(id string) (*structures.MBR, *structures.SuperBlock, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	// Crear una instancia de SuperBlock
	var sb structures.SuperBlock

	// Deserializar la estructura SuperBlock desde un archivo binario
	err = sb.Deserialize(path, int64(partition.Part_start))
	if err != nil {
		return nil, nil, "", err
	}

	return &mbr, &sb, path, nil
}
