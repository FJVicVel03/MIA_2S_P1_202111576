package structures

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type MBR struct {
	Mbr_size           int32        // Tamaño del MBR en bytes
	Mbr_creation_date  float32      // Fecha y hora de creación del MBR
	Mbr_disk_signature int32        // Firma del disco
	Mbr_disk_fit       [1]byte      // Tipo de ajuste
	Mbr_partitions     [4]PARTITION // Particiones del MBR
}

// SerializeMBR serializa la estructura MBR y la escribe en un archivo en la ruta especificada.
func (mbr *MBR) SerializeMBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeMBR deserializa la estructura MBR desde un archivo en la ruta especificada.
func (mbr *MBR) DeserializeMBR(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("tamaño de MBR inválido: %d", mbrSize)
	}

	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Print imprime la información del MBR en un formato legible.
func (mbr *MBR) Print() {
	creationTime := time.Unix(int64(mbr.Mbr_creation_date), 0)

	fmt.Printf("Tamaño del MBR: %d\n", mbr.Mbr_size)
	fmt.Printf("Fecha de Creación: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Firma del Disco: %d\n", mbr.Mbr_disk_signature)
}

// DeserializeMBR deserializa la estructura MBR desde un archivo en la ruta especificada y la retorna.
func DeserializeMBR(path string) (*MBR, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mbr := &MBR{}
	err = binary.Read(file, binary.LittleEndian, mbr)
	if err != nil {
		return nil, err
	}

	return mbr, nil
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (mbr *MBR) Serialize(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Método para obtener la primera partición disponible
func (mbr *MBR) GetFirstAvailablePartition() (*PARTITION, int, int) {
	// Calcular el offset para el start de la partición
	offset := binary.Size(mbr) // Tamaño del MBR en bytes

	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Si el start de la partición es -1, entonces está disponible
		if mbr.Mbr_partitions[i].Part_start == -1 {
			// Devolver la partición, el offset y el índice
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return nil, -1, -1
}

// Método para obtener una partición por nombre
func (mbr *MBR) GetPartitionByName(name string) (*PARTITION, int) {
	// Recorrer las particiones del MBR
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.Part_name[:]), "\x00 ")
		// Convertir el nombre de la partición a string y eliminar los caracteres nulos
		inputName := strings.Trim(name, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionName, inputName) {
			return &partition, i
		}
	}
	return nil, -1
}

// Función para obtener una partición por ID
func (mbr *MBR) GetPartitionByID(id string) (*PARTITION, error) {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionID := strings.Trim(string(mbr.Mbr_partitions[i].Part_id[:]), "\x00 ")
		// Convertir el id a string y eliminar los caracteres nulos
		inputID := strings.Trim(id, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición
		if strings.EqualFold(partitionID, inputID) {
			return &mbr.Mbr_partitions[i], nil
		}
	}
	return nil, errors.New("partición no encontrada")
}

func CreateInitialMBR(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set file size to 5MB
	if err := file.Truncate(5 * 1024 * 1024); err != nil {
		return err
	}

	mbr := MBR{
		Mbr_size:           5 * 1024 * 1024,
		Mbr_creation_date:  float32(time.Now().Unix()),
		Mbr_disk_signature: 123456, // Example signature
	}

	if err := binary.Write(file, binary.LittleEndian, &mbr); err != nil {
		return err
	}

	return nil
}

// Método para imprimir las particiones del MBR
func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.Part_status[0])
		partType := rune(partition.Part_type[0])
		partFit := rune(partition.Part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.Part_name[:])
		// Convertir Part_id a string
		partID := string(partition.Part_id[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.Part_start)
		fmt.Printf("  Size: %d\n", partition.Part_size)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlative)
		fmt.Printf("  ID: %s\n", partID)
	}
}
