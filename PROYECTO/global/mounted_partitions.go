package global

import (
	"PROYECTO/structures"
	"PROYECTO/utils"
	"errors"
	"fmt"
)

const Carnet string = "76" // 202111576
// Declaración de las particiones montadas
var (
	MountedPartitions map[string]string = make(map[string]string)
)

// GetMountedPartitionRep obtiene la partición montada por su ID
func GetMountedPartition(id string) (*structures.MBR, *structures.SuperBlock, string, error) {
	path, exists := utils.GlobalMounts[id]
	if !exists {
		return nil, nil, "", errors.New("la partición no está montada")
	}

	var mbr structures.MBR
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, nil, "", fmt.Errorf("error deserializando el MBR: %v", err)
	}

	// Aquí puedes agregar la lógica para obtener el SuperBlock si es necesario
	var sb structures.SuperBlock

	return &mbr, &sb, path, nil
}

// GetMountedPartitionSuperblock obtiene el SuperBlock de la partición montada con el id especificado
func GetMountedPartitionSuperblock(id string) (*structures.SuperBlock, *structures.PARTITION, string, error) {
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

	return &sb, partition, path, nil
}

// GetMountedPartitionRep obtiene el MBR, SuperBlock y path de la partición montada con el id especificado
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
